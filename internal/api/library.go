package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	v1 "github.com/itscontained/automarkwatched/api/v1"
	"github.com/itscontained/automarkwatched/internal/config"
)

func libraryRoutes(r chi.Router) {
	r.Get("/libraries", getLibraries)
	r.Post("/libraries", setLibraries)
}

func getLibraries(w http.ResponseWriter, r *http.Request) {
	user := GetContextUser(r)
	if err := PullLibraries(user); err != nil {
		ServerError(w)
		return
	}
	render.JSON(w, r, &user.Libraries)
}

func PullLibraries(user *v1.User) error {
	libraries, err := db.GetLibraries(user)
	if err != nil {
		return err
	}
	if user.Libraries == nil {
		user.Libraries = make(map[string]*v1.Library)
	}
	for i := range libraries {
		user.AttachLibrary(libraries[i])
	}
	return nil
}

func setLibraries(w http.ResponseWriter, r *http.Request) {
	libraries := make(map[string]*v1.Library)
	if err := render.DecodeJSON(r.Body, &libraries); err != nil {
		SendError(w, err)
		return
	}
	user := GetContextUser(r)
	if err := PullLibraries(user); err != nil {
		ServerError(w)
		return
	}
	missing := make(map[string]*v1.Library)
	ignored := 0
	for i := range libraries {
		if libraries[i].UserID != config.App.OwnerID {
			ignored++
			continue
		}
		if _, ok := user.Libraries[i]; !ok {
			missing[i] = libraries[i]
		}
	}
	if len(missing) > 0 {
		err := db.AddLibraries(missing)
		if err != nil {
			SendError(w, err)
			return
		}
	}
	untouched := len(libraries) - len(missing) - ignored
	msg := map[string]int{
		"Untouched": untouched,
		"Added":     len(missing),
		"Ignored":   ignored,
	}
	for i := range missing {
		user.AttachLibrary(missing[i])
	}
	render.JSON(w, r, &msg)
}
