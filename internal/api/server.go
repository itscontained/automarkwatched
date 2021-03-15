package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"

	v1 "github.com/itscontained/automarkwatched/api/v1"
)

func serverRoutes(r chi.Router) {
	r.Get("/servers", getServers)
	r.Post("/servers", setServers)
}

func getServers(w http.ResponseWriter, r *http.Request) {
	user := GetContextUser(r)
	if err := user.SyncServers(); err != nil {
		ErrorResponse(w, r, ErrPlexAPI)
		return
	}
	for i, s := range user.Servers {
		if s.OwnerId == 0 || s.OwnerId == user.ID {
			s.OwnerId = user.ID
			continue
		}
		delete(user.Servers, i)
	}
	render.JSON(w, r, &user.Servers)
}

func setServers(w http.ResponseWriter, r *http.Request) {
	servers := make(map[string]*v1.Server)
	if err := render.DecodeJSON(r.Body, &servers); err != nil {
		log.WithError(err).Error("problem parsing json")
		return
	}
	savedServersMap, err := db.GetServers()
	if err != nil {
		resp := Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		render.JSON(w, r, resp)
		return
	}

	missing := make(map[string]*v1.Server)
	modified := 0
	for i := range servers {
		if v, ok := savedServersMap[i]; ok {
			if compareServer(servers[i], v) {
				logr := log.WithFields(log.Fields{
					"server":             servers[i].Name,
					"machine_identifier": i,
				})
				db.UpdateServer(servers[i])
				modified++
				logr.Debug("updated series")
				continue
			}
			continue
		}
		missing[i] = servers[i]
	}
	if len(missing) > 0 {
		db.AddServers(missing)
	}

	untouched := len(servers) - len(missing) - modified
	msg := fmt.Sprintf("Untouched: %d, Added: %d, Modified: %d", untouched, len(missing), modified)
	resp := Response{
		Code:    http.StatusOK,
		Message: msg,
	}
	go setUserServers(GetContextUser(r), missing)
	render.JSON(w, r, &resp)
}

func setUserServers(user *v1.User, servers map[string]*v1.Server) {
	userServers := db.GetUserServers(user.ID)

	savedUserServersMap := make(map[string]*v1.UserServer)
	for i := range userServers {
		savedUserServersMap[i] = userServers[i]
	}

	missing := make(map[string]*v1.UserServer)
	for i := range servers {
		if _, ok := savedUserServersMap[i]; !ok {
			missing[i] = &v1.UserServer{
				UserID:                  user.ID,
				ServerMachineIdentifier: i,
				AccessToken:             servers[i].AccessToken,
				Enabled:                 true,
				User:                    user,
				Server:                  servers[i],
			}
		}
	}

	db.AddUserServers(missing)
}

func getUserServers(user *v1.User) map[string]*v1.UserServer {
	userServers := db.GetUserServers(user.ID)
	if userServers == nil {
		return nil
	}
	mappedUserServers := make(map[string]*v1.UserServer, 0)
	for i := range userServers {
		s := db.GetServer(i)
		s.Libraries = make(map[string]*v1.Library)
		s.AccessToken = userServers[i].AccessToken
		user.Servers[i] = s
		userServers[i].Server = s
		userServers[i].User = user
		mappedUserServers[i] = userServers[i]
	}
	return mappedUserServers
}

func compareServer(s1, s2 *v1.Server) bool {
	mismatched := make([]string, 0)
	if s1.Name != s2.Name {
		mismatched = append(mismatched, "name")
	}
	if s1.Scheme != s2.Scheme {
		mismatched = append(mismatched, "scheme")
	}
	if s1.Host != s2.Host {
		mismatched = append(mismatched, "host")
	}
	if s1.Port != s2.Port {
		mismatched = append(mismatched, "port")
	}
	if s1.LocalAddresses != s2.LocalAddresses {
		mismatched = append(mismatched, "local addresses")
	}
	if s1.Address != s2.Address {
		mismatched = append(mismatched, "address")
	}
	if s1.Version != s2.Version {
		mismatched = append(mismatched, "version")
	}
	if s1.Enabled != s2.Enabled {
		mismatched = append(mismatched, "enabled")
	}
	if len(mismatched) > 0 {
		return false
	}
	return true
}
