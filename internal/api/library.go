package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	v1 "github.com/itscontained/automarkwatched/api/v1"
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
	msg, err := SaveLibraries(user)
	if err != nil {
		ServerError(w)
		return
	}
	render.JSON(w, r, &msg)
}

func SaveLibraries(user *v1.User) (map[string]int, error) {
	missing := make(map[string]*v1.Library)
	ownerLibraries, err := db.GetOwnerLibraries()
	if err != nil {
		return nil, err
	}
	libraries, err := db.GetLibraries(user)
	if err != nil {
		return nil, err
	}
	for i := range ownerLibraries {
		if _, ok := libraries[i]; !ok {
			ownerLibraries[i].Enabled = true
			ownerLibraries[i].UserID = user.ID
			missing[i] = ownerLibraries[i]
		}
	}
	if len(missing) > 0 {
		err = db.AddLibraries(missing)
		if err != nil {
			return nil, err
		}
	}
	untouched := len(user.Libraries) - len(missing)
	msg := map[string]int{
		"Untouched": untouched,
		"Added":     len(missing),
	}
	for i := range missing {
		user.AttachLibrary(missing[i])
	}
	return msg, nil
}
