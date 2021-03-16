package v1

import (
	"fmt"

	"github.com/itscontained/automarkwatched/internal/config"
	"github.com/itscontained/automarkwatched/pkg/provider/plex"
)

// Server is a distinct plex server
type Server struct {
	*plex.Server
	Enabled   bool                `json:"enabled" db:"enabled"`
	Libraries map[string]*Library `json:"-" db:"-"`
}

func (s *Server) URL() string {
	return fmt.Sprintf("%s://%s:%d", s.Scheme, s.Host, s.Port)
}

func (s *Server) AttachLibrary(library *Library) error {
	err := s.AttachPlexLibrary(library.Library)
	if err != nil && err != ErrNoAppOwner {
		return err
	}
	s.Libraries[library.UUID].Enabled = library.Enabled
	s.Libraries[library.UUID].AttachServer(s)
	return nil
}

func (s *Server) AttachPlexLibrary(library *plex.Library) error {
	if config.App.OwnerID == 0 {
		return ErrNoAppOwner
	}
	if s.Libraries == nil {
		s.Libraries = make(map[string]*Library)
	}
	if _, ok := s.Libraries[library.UUID]; !ok {
		s.Libraries[library.UUID] = &Library{
			Library:                 library,
			ServerMachineIdentifier: s.MachineIdentifier,
			Enabled:                 true,
			Server:                  s,
		}
		return nil
	}
	s.Libraries[library.UUID].comparePlexLibrary(library)
	return nil
}

func (s *Server) SyncLibraries() error {
	if s.AccessToken != "" {
		return ErrNoUserAccessToken
	}
	libraries, err := plex.GetTVLibraries(s.URL(), s.AccessToken)
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
