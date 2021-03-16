package v1

// UserLibrary holds per library config for a User
type UserLibrary struct {
	UserID                  int      `db:"user_id"`
	ServerMachineIdentifier string   `db:"server_machine_identifier"`
	LibraryUUID             string   `db:"library_uuid"`
	Enabled                 bool     `db:"enabled"`
	User                    *User    `json:"-" db:"-"`
	Library                 *Library `json:"-" db:"-"`
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
