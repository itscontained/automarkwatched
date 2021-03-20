package v1

import (
	"github.com/DirtyCajunRice/go-plex"
)

// Series holds per User config of a plex.Series
type Series struct {
	RatingKey     int    `json:"rating_key" db:"rating_key"`
	Title         string `json:"title" db:"title"`
	ContentRating string `json:"contentRating,omitempty" db:"content_rating"`
	Year          int    `json:"year" db:"year"`
	Studio        string `json:"studio" db:"studio"`
	Scrobble      bool   `db:"scrobble"`
	Enabled       bool   `json:"enabled" db:"enabled"`
	// User key map
	UserID int `json:"user_id" db:"user_id"`
	// Server key map
	ServerID string `json:"server_machine_identifier" db:"server_machine_identifier"`
	// Library key map
	LibraryID string `json:"library_uuid" db:"library_uuid"`
	User      *User  `json:"-" db:"-"`
	p         *plex.Series
}

func newSeries(u *User, serverID, libraryUUID string, s *plex.Series) *Series {
	return &Series{
		UserID:        u.ID,
		RatingKey:     s.RatingKey,
		Title:         s.Title,
		ContentRating: s.ContentRating,
		Year:          s.Year,
		Scrobble:      false,
		Enabled:       true,
		Studio:        s.Studio,
		ServerID:      serverID,
		LibraryID:     libraryUUID,
		User:          u,
		p:             s,
	}
}

func (l *Library) SyncSeries() error {
	series, err := l.p.Series(false)
	if err != nil {
		return err
	}
	for _, s := range series {
		err = l.User.AttachPlexSeries(l.ServerID, l.UUID, &s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) AttachSeries(series *Series) {
	if _, ok := u.Series[series.RatingKey]; !ok {
		u.Series[series.RatingKey] = series
	}
	if u.Series[series.RatingKey].User == nil {
		u.Series[series.RatingKey].User = u
	}
}

func (u *User) AttachPlexSeries(serverID, libraryID string, series *plex.Series) error {
	if _, ok := u.Series[series.RatingKey]; !ok {
		u.Series[series.RatingKey] = newSeries(u, serverID, libraryID, series)
	}
	u.Series[series.RatingKey].update(series)
	if u.Series[series.RatingKey].User == nil {
		u.Series[series.RatingKey].User = u
	}
	return nil
}

func (s *Series) update(s2 *plex.Series) bool {
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
	if s.Year != s2.Year {
		s.Year = s2.Year
		updated = true
	}
	if s.ContentRating != s2.ContentRating {
		s.ContentRating = s2.ContentRating
		updated = true
	}
	if s.ContentRating != s2.ContentRating {
		s.ContentRating = s2.ContentRating
		updated = true
	}
	if s.ContentRating != s2.ContentRating {
		s.ContentRating = s2.ContentRating
		updated = true
	}
	if s.ContentRating != s2.ContentRating {
		s.ContentRating = s2.ContentRating
		updated = true
	}
	return updated
}
