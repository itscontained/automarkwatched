package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"

	v1 "github.com/itscontained/automarkwatched/api/v1"
)

func libraryRoutes(r chi.Router) {
	r.Get("/libraries", getLibraries)
	r.Post("/libraries", setLibraries)
}

func getLibraries(w http.ResponseWriter, r *http.Request) {
	user := GetContextUser(r)
	userServers := getUserServers(user)
	if userServers == nil {
		render.NoContent(w, r)
		return
	}
	al := make(map[string]*v1.Library)
	for i := range userServers {
		err := userServers[i].SyncLibraries()
		if err != nil {
			log.WithError(err).Error("problem getting users plex libraries")
			return
		}
		for v := range userServers[i].Server.Libraries {
			al[v] = userServers[i].Server.Libraries[v]
		}
	}
	render.JSON(w, r, &al)
}

func setLibraries(w http.ResponseWriter, r *http.Request) {
	libraries := make(map[string]*v1.Library)
	if err := render.DecodeJSON(r.Body, &libraries); err != nil {
		log.WithError(err).Error("problem parsing json")
		return
	}
	savedLibraries := db.GetLibraries()
	if savedLibraries == nil {
		render.NoContent(w, r)
	}
	missing := make(map[string]*v1.Library)
	modified := 0
	for i := range libraries {
		if _, ok := savedLibraries[libraries[i].UUID]; ok {
			if !compareLibrary(libraries[i], savedLibraries[i]) {
				l := log.WithFields(log.Fields{
					"name": libraries[i].Title,
					"uuid": libraries[i].UUID,
				})
				db.UpdateLibrary(libraries[i])
				modified++
				l.Debug("updated library")
				continue
			}
			continue
		}
		missing[i] = libraries[i]
	}

	if len(missing) > 0 {
		db.AddLibraries(missing)
	}
	untouched := len(libraries) - len(missing) - modified
	msg := fmt.Sprintf("Untouched: %d, Added: %d, Modified: %d", untouched, len(missing), modified)
	resp := Response{
		Code:    http.StatusOK,
		Message: msg,
	}
	go setUserLibraries(GetContextUser(r), missing)
	render.JSON(w, r, &resp)
}

func setUserLibraries(user *v1.User, libraries map[string]*v1.Library) {
	userLibraries := db.GetUserLibraries(user)
	if userLibraries == nil {
		return
	}
	missing := make(map[string]*v1.UserLibrary)
	for i := range libraries {
		if _, ok := userLibraries[i]; !ok {
			s := db.GetServer(libraries[i].ServerMachineIdentifier)
			missing[i] = &v1.UserLibrary{
				UserID:                  user.ID,
				ServerMachineIdentifier: s.MachineIdentifier,
				LibraryUUID:             libraries[i].UUID,
				Enabled:                 true,
				User:                    user,
				Library:                 libraries[i],
			}
			log.Printf("%#v", user.ID)
			log.Printf("%#v", user)
		}
	}
	if len(missing) > 0 {
		db.AddUserLibraries(missing)
	}
}

func getUserLibraries(user *v1.User) map[string]*v1.UserLibrary {
	userLibraries := db.GetUserLibraries(user)
	if userLibraries == nil {
		return nil
	}

	attached := make(map[string]*v1.UserLibrary, 0)
	for i := range userLibraries {
		userLibraries[i].User = user
		userLibraries[i].Library = db.GetLibrary(userLibraries[i].LibraryUUID)
		userLibraries[i].Library.Series = make(map[int]*v1.Series)
		attached[i] = userLibraries[i]
		if _, ok := user.Servers[userLibraries[i].ServerMachineIdentifier]; ok {
			userLibraries[i].Library.Server = user.Servers[userLibraries[i].ServerMachineIdentifier]
			continue
		}
		userServer := db.GetUserServer(user, userLibraries[i].ServerMachineIdentifier)
		userLibraries[i].Library.Server = db.GetServer(userLibraries[i].ServerMachineIdentifier)
		userLibraries[i].Library.Server.AccessToken = userServer.AccessToken
		user.Servers[userLibraries[i].ServerMachineIdentifier] = userLibraries[i].Library.Server
	}
	return attached
}

func compareLibrary(s1, s2 *v1.Library) bool {
	mismatched := make([]string, 0)
	if s1.Title != s2.Title {
		mismatched = append(mismatched, "title")
	}
	if s1.Key != s2.Key {
		mismatched = append(mismatched, "key")
	}
	if s1.Agent != s2.Agent {
		mismatched = append(mismatched, "agent")
	}
	if s1.Scanner != s2.Scanner {
		mismatched = append(mismatched, "scanner")
	}
	if s1.Enabled != s2.Enabled {
		mismatched = append(mismatched, "address")
	}
	if len(mismatched) > 0 {
		return false
	}
	return true
}
