package v1

import (
	"github.com/itscontained/automarkwatched/pkg/provider/plex"
)

// User is a wrapper struct for a plex.User
type User struct {
	*plex.User
	Owner   bool               `json:"owner" db:"owner" goqu:"skipupdate"`
	Exists  bool               `json:"exists" db:"-"`
	Enabled bool               `json:"enabled" db:"enabled"`
	Servers map[string]*Server `json:"-" db:"-"`
}

// Server is a distinct plex server
type Server struct {
	*plex.Server
	Enabled   bool                `db:"enabled"`
	Libraries map[string]*Library `json:"-" db:"-"`
}

// UserServer holds per server config for a User
type UserServer struct {
	UserID                  int     `db:"user_id" goqu:"skipupdate"`
	ServerMachineIdentifier string  `db:"server_machine_identifier" goqu:"skipupdate"`
	AccessToken             string  `db:"access_token"`
	Enabled                 bool    `db:"enabled"`
	User                    *User   `json:"-" db:"-"`
	Server                  *Server `json:"-" db:"-"`
}

// Library is a wrapper struct for a plex.Library
type Library struct {
	*plex.Library
	ServerMachineIdentifier string          `db:"server_machine_identifier" goqu:"skipupdate"`
	Enabled                 bool            `json:"enabled" db:"enabled"`
	Series                  map[int]*Series `json:"-" db:"-"`
	Server                  *Server         `json:"-" db:"-"`
}

// UserLibrary holds per library config for a User
type UserLibrary struct {
	UserID                  int      `db:"user_id"`
	ServerMachineIdentifier string   `db:"server_machine_identifier"`
	LibraryUUID             string   `db:"library_uuid"`
	Enabled                 bool     `db:"enabled"`
	User                    *User    `json:"-" db:"-"`
	Library                 *Library `json:"-" db:"-"`
}

// Series is a wrapper struct for a plex.Series
type Series struct {
	*plex.Series
	Enabled                 bool     `db:"enabled"`
	ServerMachineIdentifier string   `db:"server_machine_identifier"`
	LibraryUUID             string   `db:"library_uuid"`
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
