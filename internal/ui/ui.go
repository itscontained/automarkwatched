package ui

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/itscontained/automarkwatched/internal/api"
	. "github.com/itscontained/automarkwatched/internal/config"
	"github.com/itscontained/automarkwatched/internal/database"
)

var (
	r         *chi.Mux
	templates *template.Template
	db        = database.DB
	server    *http.Server
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

	l := log.WithFields(log.Fields{
		"process": "webserver",
	})
	r.Group(func(r chi.Router) {
		r.Use(api.TokenContext)
		r.Use(api.UserContext)
		r.Get("/", index)
	})
	r.Route("/", func(r chi.Router) {
		r.Get("/setup", setup)
		r.Get("/login", login)
		r.With(middleware.NoCache).Get("/static/*", static)
	})
	templates, err = template.New("base").ParseGlob("web/dist/templates/*.tmpl")
	if err != nil {
		l.WithError(err).Error("could not parse templates")
	}
	server = &http.Server{
		Addr:    ":" + strconv.Itoa(App.WebPort),
		Handler: r,
	}
	go func() {
		l.Infof("ui listening on :%d", App.WebPort)
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
		log.WithError(err).Fatalf("could not gracefully shutdown webserver")
	}
}

func static(w http.ResponseWriter, r *http.Request) {
	fs := http.StripPrefix("/static", http.FileServer(http.Dir("./web/dist/static")))
	fs.ServeHTTP(w, r)
}

func setup(w http.ResponseWriter, _ *http.Request) {
	renderTemplate(w, "setup.tmpl", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	_, _, err := api.GetRequestCreds(r)
	if err != nil {
		renderTemplate(w, "login.tmpl", nil)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func index(w http.ResponseWriter, r *http.Request) {
	user := api.GetContextUser(r)

	if err := api.PullAll(user); err != nil {
		api.SendError(w, err)
		return
	}
	data := map[string]interface{}{
		"user": user,
	}
	renderTemplate(w, "index.tmpl", data)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
	t, err := template.New("base").ParseGlob("web/dist/templates/*.tmpl")
	if err != nil {
		log.WithError(err).Error("could not parse templates")
	}
	err = t.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("error 500: %s", err.Error()), http.StatusInternalServerError)
	}
}
