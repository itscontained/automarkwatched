package v1

import (
	"fmt"
	"reflect"

	"github.com/DirtyCajunRice/go-plex"

	"github.com/itscontained/automarkwatched/internal/config"
)

// Server holds per User app config for a plex.Server
type Server struct {
	UserID            int    `db:"user_id" goqu:"skipupdate"`
	Address           string `json:"address" db:"address"`
	Host              string `json:"host" db:"host"`
	LocalAddresses    string `json:"local_addresses" db:"local_addresses"`
	Name              string `json:"name" db:"name"`
	Port              int    `json:"port" db:"port"`
	Scheme            string `json:"scheme" db:"scheme"`
	Version           string `json:"version" db:"version"`
	OwnerId           int    `json:"owner_id" db:"owner_id"`
	MachineIdentifier string `json:"machine_identifier" db:"machine_identifier" goqu:"skipupdate"`
	AccessToken       string `json:"access_token" db:"access_token"`
	Enabled           bool   `json:"enabled" db:"enabled"`
	User              *User  `json:"-" db:"-"`
	p                 *plex.Server
}

func newUserServer(u *User, s *plex.Server) *Server {
	return &Server{
		UserID:            u.ID,
		Address:           s.Address,
		Host:              s.Host,
		LocalAddresses:    s.LocalAddresses,
		Name:              s.Name,
		Port:              s.Port,
		Scheme:            s.Scheme,
		Version:           s.Version,
		OwnerId:           s.OwnerId,
		MachineIdentifier: s.MachineIdentifier,
		AccessToken:       s.AccessToken,
		Enabled:           true,
		User:              u,
		p:                 s,
	}
}

func (s *Server) URL() string {
	return fmt.Sprintf("%s://%s:%d", s.Scheme, s.Host, s.Port)
}

func (u *User) SyncServers() error {
	if u.AuthToken == "" {
		return ErrNoUserAuthToken
	}
	servers, err := u.p.Servers()
	if err != nil {
		return err
	}
	for _, server := range servers {
		err = u.AttachPlexServer(&server)
		if err != nil || err == ErrNoAppOwner {
			return err
		}
	}
	return nil
}

func (u *User) AttachServer(server *Server) {
	if server.p == nil {
		server.p = &plex.Server{
			Host:        server.Host,
			Scheme:      server.Scheme,
			Port:        server.Port,
			AccessToken: server.AccessToken,
		}
		App.AttachServer(server.p)
	}
	if _, ok := u.Servers[server.MachineIdentifier]; !ok {
		u.Servers[server.MachineIdentifier] = server
	}
	if u.Servers[server.MachineIdentifier].User == nil {
		u.Servers[server.MachineIdentifier].User = u
	}
}

func (u *User) AttachPlexServer(server *plex.Server) error {
	if config.App.OwnerID == 0 {
		return ErrNoAppOwner
	}
	if u.ID == config.App.OwnerID && server.OwnerId == 0 {
		server.OwnerId = u.ID
	}
	if server.OwnerId != config.App.OwnerID {
		return ErrNotApplicationOwned
	}
	if _, ok := u.Servers[server.MachineIdentifier]; !ok {
		u.Servers[server.MachineIdentifier] = newUserServer(u, server)
		return nil
	}
	u.Servers[server.MachineIdentifier].update(server)
	if u.Servers[server.MachineIdentifier].User == nil {
		u.Servers[server.MachineIdentifier].User = u
	}
	return nil
}

func (s *Server) update(s2 *plex.Server) bool {
	if reflect.DeepEqual(s.p, s2) {
		return false
	}
	updated := false
	if s.Name != s2.Name {
		s.Name = s2.Name
		updated = true
	}
	if s.Scheme != s2.Scheme {
		s.Scheme = s2.Scheme
		updated = true
	}
	if s.Host != s2.Host {
		s.Host = s2.Host
		updated = true
	}
	if s.Port != s2.Port {
		s.Port = s2.Port
		updated = true
	}
	if s.LocalAddresses != s2.LocalAddresses {
		s.LocalAddresses = s2.LocalAddresses
		updated = true
	}
	if s.Address != s2.Address {
		s.Address = s2.Address
		updated = true
	}
	if s.Version != s2.Version {
		s.Version = s2.Version
		updated = true
	}
	s.p = s2
	return updated
}

func (s *Server) Scrobble(ratingKey int) error {
	return s.p.Scrobble(ratingKey)
}
