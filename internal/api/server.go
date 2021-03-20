package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	v1 "github.com/itscontained/automarkwatched/api/v1"
	"github.com/itscontained/automarkwatched/internal/config"
)

func serverRoutes(r chi.Router) {
	r.Get("/servers", getServers)
	r.Post("/servers", setServers)
}

func getServers(w http.ResponseWriter, r *http.Request) {
	user := GetContextUser(r)
	if err := PullServers(user); err != nil {
		ServerError(w)
		return
	}
	render.JSON(w, r, &user.Servers)
}

func PullServers(user *v1.User) error {
	servers, err := db.GetServers(user)
	if err != nil {
		return err
	}
	if user.Servers == nil {
		user.Servers = make(map[string]*v1.Server)
	}
	for i := range servers {
		user.AttachServer(servers[i])
	}
	return nil
}

func setServers(w http.ResponseWriter, r *http.Request) {
	servers := make(map[string]*v1.Server)
	if err := render.DecodeJSON(r.Body, &servers); err != nil {
		SendError(w, err)
		return
	}
	user := GetContextUser(r)
	if err := PullServers(user); err != nil {
		ServerError(w)
		return
	}
	missing := make(map[string]*v1.Server)
	ignored := 0
	for i := range servers {
		if servers[i].OwnerId != config.App.OwnerID {
			ignored++
			continue
		}
		if _, ok := user.Servers[i]; !ok {
			missing[i] = servers[i]
		}
	}
	if len(missing) > 0 {
		err := db.AddServers(missing)
		if err != nil {
			SendError(w, err)
			return
		}
	}
	untouched := len(servers) - len(missing) - ignored
	msg := map[string]int{
		"Untouched": untouched,
		"Added":     len(missing),
		"Ignored":   ignored,
	}
	for i := range missing {
		user.AttachServer(missing[i])
		err := user.SyncServers()
		if err != nil {
			SendError(w, err)
			return
		}
	}
	render.JSON(w, r, &msg)
}
