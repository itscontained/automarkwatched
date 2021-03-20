package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/DirtyCajunRice/go-plex"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"

	v1 "github.com/itscontained/automarkwatched/api/v1"
	. "github.com/itscontained/automarkwatched/internal/config"
	"github.com/itscontained/automarkwatched/internal/database"
)

var (
	decoder = schema.NewDecoder()
	db      = database.DB
	server  *http.Server
	r       *chi.Mux
	A       *plex.App
)

func Start() {
	var err error
	r = chi.NewRouter()
	// injects a request ID into the context of each request. A request ID is a
	// string of the form "host.example.com/random-0001", where "random" is a base62
	// random string that uniquely identifies this go process, and where the last
	// number is an atomically incremented request counter
	r.Use(middleware.RequestID)
	// sets a http.Request's RemoteAddr to the results of parsing either the
	// X-Forwarded-For header or the X-Real-IP header (in that order)
	r.Use(middleware.RealIP)
	// logs the start and end of each request, along with some useful data about what
	// was requested, what the response status was, and how long it took to return
	r.Use(middleware.Logger)
	// recovers from, logs (and a backtrace), and returns a HTTP 500 status if possible
	r.Use(middleware.Recoverer)
	// CleanPath middleware will clean out double slash mistakes from a user's request path
	r.Use(middleware.CleanPath)
	// cancels ctx after a given timeout and return a 504 Gateway Timeout error to the client
	r.Use(middleware.Timeout(60 * time.Second))
	// enforces a whitelist of request Content-Types otherwise responds with a 415 Unsupported Media Type status
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "X-Plex-Token", "X-Plex-ID", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// routes with token required
			r.Group(func(r chi.Router) {
				r.Use(TokenContext)
				r.Get("/user", getUser)
				r.Post("/user", setUser)
				r.Group(func(r chi.Router) {
					r.Use(UserContext)
					serverRoutes(r)
					libraryRoutes(r)
					seriesRoutes(r)
				})
			})
			// routes without token required
			r.Group(func(r chi.Router) {
				r.Get("/login", login)
			})
		})
	})

	server = &http.Server{
		Addr:    ":" + strconv.Itoa(App.APIPort),
		Handler: r,
	}
	l := log.WithFields(log.Fields{
		"process": "api-server",
	})

	go func() {
		l.Infof("api listening on :%d", App.APIPort)
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.WithError(err).Fatal("")
		}
	}()
}

func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Fatalf("could not gracefully shutdown api-server")
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	pin, err := A.GeneratePin()
	if err != nil {
		log.Error(err)
		return
	}
	data := map[string]interface{}{
		"url": pin.AuthUrl(),
		"pin": pin,
	}
	render.JSON(w, r, data)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	token, id, err := getContextUserData(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	newUser := v1.NewPartialUser(id, token)
	savedUser := db.GetUser(newUser.ID)
	if savedUser == nil {
		render.JSON(w, r, newUser)
		return
	}
	newUser.Update(savedUser)
	render.JSON(w, r, savedUser)
}

func setUser(w http.ResponseWriter, r *http.Request) {
	var user *v1.User
	if err := render.DecodeJSON(r.Body, &user); err != nil {
		log.WithError(err).Error("problem parsing json")
		return
	}
	if err := db.UpdateUser(user); err != nil {
		log.WithError(err).Error("problem updating  parsing json")
		return
	}
	render.NoContent(w, r)
}

func TokenContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, id, err := getContextUserData(r)
		if err != nil {
			token = getRequestToken(r)
			id, _ = getRequestUserID(r)
			if token != "" && id != 0 {
				tokenIdMap := map[string]interface{}{"token": token, "id": id}
				ctx := context.WithValue(r.Context(), "plexAuthToken", tokenIdMap)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetContextUser(r)
		if user == nil {
			token, id, err := getContextUserData(r)
			if err != nil {
				log.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			u := v1.NewPartialUser(id, token)
			savedUser := db.GetUser(u.ID)
			if savedUser == nil {
				log.Error("fuck")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			u.Update(savedUser)
			if !u.Enabled {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			for j := range u.Libraries {
				log.Warnf("indexset: %+v", u.Libraries[j].PrintKey())
			}
			ctx := context.WithValue(r.Context(), "user", u)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetContextUser(r *http.Request) *v1.User {
	user, ok := r.Context().Value("user").(*v1.User)
	if !ok {
		return nil
	}
	savedUser := db.GetUser(user.ID)
	user.Update(savedUser)
	return user
}

func getRequestToken(r *http.Request) string {
	if t := r.URL.Query().Get("X-Plex-Token"); t != "" {
		return t
	}
	if t := r.Header.Get("X-Plex-Token"); t != "" {
		return t
	}
	if t, err := r.Cookie("X-Plex-Token"); err == nil {
		if t.Value != "" {
			return t.Value
		}
	}
	return ""
}

func getRequestUserID(r *http.Request) (int, error) {
	if t := r.URL.Query().Get("X-Plex-ID"); t != "" {
		return strconv.Atoi(t)
	}
	if t := r.Header.Get("X-Plex-ID"); t != "" {
		return strconv.Atoi(t)
	}
	if t, err := r.Cookie("X-Plex-ID"); err == nil {
		if t.Value != "" {
			return strconv.Atoi(t.Value)
		}
	}
	return 0, errors.New("no plex ID found")
}

func getContextUserData(r *http.Request) (string, int, error) {
	tokenIdMap, ok := r.Context().Value("plexAuthToken").(map[string]interface{})
	if !ok {
		return "", 0, errors.New("no context user data")
	}
	token := tokenIdMap["token"].(string)
	id := tokenIdMap["id"].(int)
	return token, id, nil
}
