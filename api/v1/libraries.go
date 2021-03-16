package v1

import (
	"github.com/itscontained/automarkwatched/internal/config"
	"github.com/itscontained/automarkwatched/pkg/provider/plex"
)

// Library is a wrapper struct for a plex.Library
type Library struct {
	*plex.Library
	ServerMachineIdentifier string          `db:"server_machine_identifier" goqu:"skipupdate"`
	Enabled                 bool            `json:"enabled" db:"enabled"`
	Series                  map[int]*Series `json:"-" db:"-"`
	Server                  *Server         `json:"-" db:"-"`
}

func (l *Library) AttachServer(server *Server) {
	l.Server = server
	l.ServerMachineIdentifier = server.MachineIdentifier
}

func (l *Library) AttachSeries(series *Series) error {
	err := l.AttachPlexSeries(series.Series)
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
