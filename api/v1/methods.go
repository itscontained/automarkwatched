package v1

import (
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"

	"github.com/itscontained/automarkwatched/pkg/provider/plex"
)

func (u *User) Sync() error {
	user, err := plex.GetUser(u.AuthToken)
	if err != nil {
		return err
	}
	u.User = &user
	u.Servers = make(map[string]*Server)
	return nil
}

func (u *User) SyncServers() error {
	servers, err := plex.GetServers(u.AuthToken)
	if err != nil {
		return err
	}
	for _, server := range servers {
		if _, ok := u.Servers[server.MachineIdentifier]; !ok {
			u.Servers[server.MachineIdentifier] = &Server{
				Enabled:   true,
				Libraries: make(map[string]*Library),
			}
		}
		u.Servers[server.MachineIdentifier].Server = &server
	}
	return nil
}

func (s *Server) URL() string {
	return fmt.Sprintf("%s://%s:%d", s.Scheme, s.Host, s.Port)
}

func (s *Server) SyncLibraries(accessToken string) error {
	libraries, err := plex.GetTVLibraries(s.URL(), accessToken)
	if err != nil {
		return err
	}
	for i := range libraries {
		if _, ok := s.Libraries[libraries[i].UUID]; !ok {
			s.Libraries[libraries[i].UUID] = &Library{
				ServerMachineIdentifier: s.MachineIdentifier,
				Server:                  s,
				Enabled:                 true,
				Series:                  make(map[int]*Series),
			}
		}
		s.Libraries[libraries[i].UUID].Library = &libraries[i]
	}
	return nil
}

func (us *UserServer) SyncLibraries() error {
	return us.Server.SyncLibraries(us.AccessToken)
}

func (ul *UserLibrary) SyncSeries() error {
	return ul.Library.SyncSeries(ul.Library.Server.AccessToken)
}
func (l *Library) SyncSeries(accessToken string) error {
	log.Printf("%s, %s, %d", l.Server.URL(), accessToken, l.Library.Key)
	plexSeries, err := plex.GetTVSeries(l.Server.URL(), accessToken, l.Library.Key, false)
	if err != nil {
		return err
	}

	for i := range plexSeries {
		if _, ok := l.Series[plexSeries[i].RatingKey]; !ok {
			l.Series[plexSeries[i].RatingKey] = &Series{
				Library:     l,
				Enabled:     true,
				LibraryUUID: l.UUID,
			}
		}
		l.Series[plexSeries[i].RatingKey].Series = &plexSeries[i]
	}
	return nil
}

func NewSortedUserSeriesSlice(userSeries map[int]*UserSeries) []*UserSeries {
	keys := make([]string, 0)
	maps := make(map[string]int)
	for ratingKey := range userSeries {
		keys = append(keys, userSeries[ratingKey].Series.Title)
		maps[userSeries[ratingKey].Series.Title] = ratingKey
	}
	sort.Strings(keys)
	newSlice := make([]*UserSeries, len(userSeries)+1)
	for i, v := range keys {
		newSlice[i] = userSeries[maps[v]]
	}
	return newSlice
}
