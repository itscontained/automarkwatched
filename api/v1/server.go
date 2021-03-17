package v1

import (
	"github.com/DirtyCajunRice/go-plex"
)

// Server is an app config wrapper for a plex.Server
type Server struct {
	*plex.Server
	Enabled bool `json:"enabled" db:"enabled"`
}

func (s *Server) SyncLibraries() error {
	if s.AccessToken == "" {
		return ErrNoUserAccessToken
	}
	libraries, err := s.Server.Libraries()
	if err != nil {
		return err
	}
	for _, l := range libraries {
		err = s.AttachPlexLibrary(&l)
		if err != nil || err == ErrNoAppOwner {
			return err
		}
	}
	return nil
}

func (s *Server) compareServer(s2 *Server) bool {
	updated := s.comparePlexServer(s2.Server)
	if s.Enabled != s2.Enabled {
		s.Enabled = s2.Enabled
		updated = true
	}
	return updated
}

func (s *Server) comparePlexServer(s2 *plex.Server) bool {
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
	return updated
}

// UserServer holds per User app config for a Server
type UserServer struct {
	UserID                  int     `db:"user_id" goqu:"skipupdate"`
	ServerMachineIdentifier string  `db:"server_machine_identifier" goqu:"skipupdate"`
	AccessToken             string  `db:"access_token"`
	Enabled                 bool    `db:"enabled"`
	User                    *User   `json:"-" db:"-"`
	Server                  *Server `json:"-" db:"-"`
}

func (s *UserServer) URL() string {
	return s.Server.URL()
}
