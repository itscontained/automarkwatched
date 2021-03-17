package v1

import (
	"sort"

	"github.com/DirtyCajunRice/go-plex"
)

// Series is a wrapper struct for a plex.Series
type Series struct {
	*plex.Series
	Enabled                 bool        `db:"enabled"`
	ServerMachineIdentifier string      `db:"server_machine_identifier"`
	LibraryUUID             string      `db:"library_uuid"`
	Library                 *Library    `json:"-" db:"-"`
	UserSeries              *UserSeries `json:"-" db:"-"`
}

func (s *Series) AttachLibrary(library *Library) {
	s.Library = library
	s.ServerMachineIdentifier = library.ServerMachineIdentifier
}

func (s *Series) compareSeries(s2 *Series) bool {
	updated := s.comparePlexSeries(s2.Series)
	if s.Enabled != s2.Enabled {
		s.Title = s2.Title
		updated = true
	}
	return updated
}

func (s *Series) comparePlexSeries(s2 *plex.Series) bool {
	updated := false
	if s.Title != s2.Title {
		s.Title = s2.Title
		updated = true
	}
	if s.RatingKey != s2.RatingKey {
		s.RatingKey = s2.RatingKey
		updated = true
	}
	if s.Studio != s2.Studio {
		s.Studio = s2.Studio
		updated = true
	}
	if s.ContentRating != s2.ContentRating {
		s.ContentRating = s2.ContentRating
		updated = true
	}
	return updated
}

func (ul *UserLibrary) SyncSeries() error {
	return ul.Library.SyncSeries()
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
