package v1

import (
	"github.com/DirtyCajunRice/go-plex"
)

// Library holds per User config of a plex.Library
type Library struct {
	Key      int    `json:"key,string" db:"key"`
	Type     string `json:"type" db:"type"`
	Title    string `json:"title" db:"title"`
	Agent    string `json:"agent" db:"agent"`
	Scanner  string `json:"scanner" db:"scanner"`
	UUID     string `json:"uuid" db:"uuid" goqu:"skipupdate"`
	UserID   int    `json:"user_id" db:"user_id"`
	Enabled  bool   `json:"enabled" db:"enabled"`
	ServerID string `json:"server_machine_identifier" db:"server_machine_identifier"`
	User     *User  `json:"-" db:"-"`
	p        plex.Library
}

func (l *Library) PrintKey() int {
	return l.p.Key
}

func newLibrary(u *User, serverID string, l plex.Library) *Library {
	return &Library{
		Key:      l.Key,
		Type:     l.Type,
		Title:    l.Title,
		Agent:    l.Agent,
		Scanner:  l.Scanner,
		UUID:     l.UUID,
		UserID:   u.ID,
		Enabled:  true,
		ServerID: serverID,
		User:     u,
		p:        l,
	}
}

func (s *Server) SyncLibraries() error {
	libraries, err := s.p.Libraries()
	if err != nil {
		return err
	}
	for _, l := range libraries {
		if l.Type != "show" {
			continue
		}
		err = s.User.AttachPlexLibrary(s.MachineIdentifier, l)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) AttachLibrary(library *Library) {
	if _, ok := u.Libraries[library.UUID]; !ok {
		u.Libraries[library.UUID] = library
	}
	if u.Libraries[library.UUID].User == nil {
		u.Libraries[library.UUID].User = u
	}
}

func (u *User) AttachPlexLibrary(serverID string, library plex.Library) error {
	if u.Libraries == nil {
		u.Libraries = make(map[string]*Library)
	}
	if _, ok := u.Libraries[library.UUID]; !ok {
		u.Libraries[library.UUID] = newLibrary(u, serverID, library)
		return nil
	}

	u.Libraries[library.UUID].update(library)
	if u.Libraries[library.UUID].User == nil {
		u.Libraries[library.UUID].User = u
	}
	return nil
}

func (l *Library) update(l2 plex.Library) bool {

	updated := false
	if l.Title != l2.Title {
		l.Title = l2.Title
		updated = true
	}
	if l.Type != l2.Type {
		l.Type = l2.Type
		updated = true
	}
	if l.Key != l2.Key {
		l.Key = l2.Key
		updated = true
	}
	if l.Agent != l2.Agent {
		l.Agent = l2.Agent
		updated = true
	}
	if l.Scanner != l2.Scanner {
		l.Scanner = l2.Scanner
		updated = true
	}
	if l.UUID != l2.UUID {
		l.UUID = l2.UUID
		updated = true
	}
	l.p = l2
	return updated
}

func (l *Library) Unwatched() ([]plex.Series, error) {
	return l.p.Series(true)
}
