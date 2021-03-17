package v1

import (
	"github.com/DirtyCajunRice/go-plex"

	"github.com/itscontained/automarkwatched/internal/config"
)

func (s *UserServer) AttachUserLibrary(library *UserLibrary) error {
	err := s.AttachLibrary(library.Library)
	if err != nil && err != ErrNoAppOwner {
		return err
	}
	s.Libraries[library.LibraryUUID].UserID = library.UserID
	s.Libraries[library.LibraryUUID].ServerMachineIdentifier = library.ServerMachineIdentifier
	s.Libraries[library.LibraryUUID].LibraryUUID = library.LibraryUUID
	s.Libraries[library.LibraryUUID].Enabled = library.Enabled
	s.Libraries[library.LibraryUUID].User = s.User
	return nil
}

func (s *UserServer) AttachLibrary(library *Library) error {
	err := s.AttachPlexLibrary(library.Library)
	if err != nil && err != ErrNoAppOwner {
		return err
	}
	s.Libraries[library.UUID].Library.ServerMachineIdentifier = library.ServerMachineIdentifier
	s.Libraries[library.UUID].Library.Enabled = library.Enabled
	s.Libraries[library.UUID].Library.Server = s.Server
	return nil
}

func (s *UserServer) AttachPlexLibrary(library *plex.Library) error {
	if s.Libraries == nil {
		s.Libraries = make(map[string]*UserLibrary)
	}
	if _, ok := s.Libraries[library.UUID]; !ok {
		s.Libraries[library.UUID] = &UserLibrary{
			Library: &Library{
				Library:                 nil,
				ServerMachineIdentifier: "",
				Enabled:                 false,
				Series:                  nil,
				Server:                  nil,
			},
			ServerMachineIdentifier: s.MachineIdentifier,
			Enabled:                 true,
			Server:                  s,
		}
		return nil
	}
	s.Libraries[library.UUID].comparePlexLibrary(library)
	return nil
}

func (s *UserServer) SyncLibraries() error {
	return us.Server.SyncLibraries()
}

// UserSeries holds per series config for a User
type UserSeries struct {
	UserID                  int     `db:"user_id"`
	ServerMachineIdentifier string  `db:"server_machine_identifier"`
	LibraryUUID             string  `db:"library_uuid"`
	SeriesRatingKey         int     `db:"series_rating_key"`
	Scrobble                bool    `db:"scrobble"`
	Enabled                 bool    `db:"enabled"`
	User                    *User   `json:"-" db:"-"`
	Series                  *Series `json:"-" db:"-"`
}
