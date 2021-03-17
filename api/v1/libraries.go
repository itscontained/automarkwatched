package v1

import (
	"github.com/DirtyCajunRice/go-plex"

	"github.com/itscontained/automarkwatched/internal/config"
)

// Library is a wrapper struct for a plex.Library
type Library struct {
	*plex.Library
	ServerMachineIdentifier string `db:"server_machine_identifier" goqu:"skipupdate"`
	Enabled                 bool   `json:"enabled" db:"enabled"`
}

// UserLibrary holds per User config of a Library
type UserLibrary struct {
	UserID      int      `json:"user_id" db:"user_id"`
	LibraryUUID string   `json:"library_uuid" db:"library_uuid"`
	Enabled     bool     `json:"enabled" db:"enabled"`
	User        *User    `json:"-" db:"-"`
	Library     *Library `json:"-" db:"-"`
}

// AttachUserServer adds the UserServer pointer to the Library struct
func (l *Library) AttachUserServer(userServer *UserServer) {
	l.Server = userServer
	l.ServerMachineIdentifier = userServer.Server.MachineIdentifier
}

func (l *Library) AttachUserSeries(userSeries *UserSeries) error {
	err := l.AttachPlexSeries(s.Series)
	if err != nil && err != ErrNoAppOwner {
		return err
	}
	l.Series[series.RatingKey].Enabled = series.Enabled
	return nil
}

func (l *Library) AttachPlexSeries(series *plex.Series) error {
	if config.App.OwnerID == 0 {
		return ErrNoAppOwner
	}
	if l.Series == nil {
		l.Series = make(map[int]*Series)
	}
	if _, ok := l.Series[series.RatingKey]; !ok {
		l.Series[series.RatingKey] = &Series{
			Series:  series,
			Enabled: true,
		}
	}
	l.Series[series.RatingKey].comparePlexSeries(series)
	l.Series[series.RatingKey].AttachLibrary(l)
	return nil
}
func (l *Library) SyncSeries() error {
	plexSeries, err := plex.GetTVSeries(l.Server.URL(), l.Server.AccessToken, l.Library.Key, false)
	if err != nil {
		return err
	}
	for _, s := range plexSeries {
		err = l.AttachPlexSeries(&s)
		if err != nil || err == ErrNoAppOwner {
			return err
		}
	}
	return nil
}

func (l *Library) compareLibrary(s2 *Library) bool {
	updated := l.comparePlexLibrary(s2.Library)
	if l.Enabled != s2.Enabled {
		l.Enabled = s2.Enabled
		updated = true
	}
	return updated
}

func (l *Library) comparePlexLibrary(s2 *plex.Library) bool {
	updated := false
	if l.Title != s2.Title {
		l.Title = s2.Title
		updated = true
	}
	if l.Type != s2.Type {
		l.Type = s2.Type
		updated = true
	}
	if l.Key != s2.Key {
		l.Key = s2.Key
		updated = true
	}
	if l.Agent != s2.Agent {
		l.Agent = s2.Agent
		updated = true
	}
	if l.Scanner != s2.Scanner {
		l.Scanner = s2.Scanner
		updated = true
	}
	return updated
}
