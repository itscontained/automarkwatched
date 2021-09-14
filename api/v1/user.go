package v1

import (
	"errors"
	"reflect"

	"github.com/DirtyCajunRice/go-plex"

	"github.com/itscontained/automarkwatched/internal/config"
)

var (
	ErrNotApplicationOwned = errors.New("not application owner server")
	ErrNoUserAuthToken     = errors.New("auth token not set")
	ErrNoUserAccessToken   = errors.New("access token not set")
	ErrNoAppOwner          = errors.New("app owner not set")
	App                    *plex.App
)

// User is a wrapper struct for a plex.User
type User struct {
	ID        int                 `json:"id" db:"id" goqu:"skipupdate"`
	UUID      string              `json:"uuid" db:"uuid"`
	Username  string              `json:"username" db:"username"`
	Email     string              `json:"email" db:"email"`
	Thumb     string              `json:"thumb" db:"thumb"`
	Owner     bool                `json:"owner" db:"owner" goqu:"skipupdate"`
	Enabled   bool                `json:"enabled" db:"enabled"`
	AuthToken string              `json:"auth_token" db:"-"`
	Servers   map[string]*Server  `json:"-" db:"-"`
	Libraries map[string]*Library `json:"-" db:"-"`
	Series    map[int]*Series     `json:"-" db:"-"`
	p         *plex.User
}

func NewPartialUser(id int, authToken string) *User {
	return &User{
		ID:        id,
		Enabled:   true,
		AuthToken: authToken,
		Servers:   make(map[string]*Server),
		Libraries: make(map[string]*Library),
		Series:    make(map[int]*Series),
	}
}

func (u *User) Attach() *User {
	if u.p == nil {
		u.createPlexUser()
	}
	App.AttachUser(u.p)
	return u
}

func (u *User) SyncUser() error {
	if u.AuthToken == "" {
		return ErrNoUserAuthToken
	}
	user, err := App.User(u.AuthToken)
	if err != nil {
		return err
	}
	u.update(user)

	if u.Servers == nil {
		u.Servers = make(map[string]*Server)
	}
	if u.Libraries == nil {
		u.Libraries = make(map[string]*Library)
	}
	if u.Series == nil {
		u.Series = make(map[int]*Series)
	}
	return nil
}

// AttachServers is a convenience function that accepts the standard map[string]*Server db response
// to call AttachServer multiple times
func (u *User) AttachServers(servers map[string]*Server) {
	for s := range servers {
		u.AttachServer(servers[s])
	}
}

func (u *User) GetAll() error {
	err := u.SyncUser()
	if err != nil {
		return err
	}
	err = u.SyncServers()
	if err != nil {
		return err
	}
	return u.GetFromExisting()
}

func (u *User) GetFromExisting() error {
	for i := range u.Servers {
		u.Servers[i].p.Scheme = "https"
		u.Servers[i].p.Host = "plex.cajun.pro"
		u.Servers[i].Scheme = "https"
		u.Servers[i].Host = "plex.cajun.pro"
		err := u.Servers[i].SyncLibraries()
		if err != nil {
			return err
		}
	}
	for i := range u.Libraries {
		err := u.Libraries[i].SyncSeries()
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) update(u2 *plex.User) bool {
	updated := false
	if reflect.DeepEqual(u.p, u2) {
		return false
	}
	if u.ID != u2.ID {
		u.ID = u2.ID
		updated = true
	}
	if u.UUID != u2.UUID {
		u.UUID = u2.UUID
		updated = true
	}
	if u.Username != u2.Username {
		u.Username = u2.Username
		updated = true
	}
	if u.Email != u2.Email {
		u.Email = u2.Email
		updated = true
	}
	if u.Thumb != u2.Thumb {
		u.Thumb = u2.Thumb
		updated = true
	}
	if u.AuthToken != u2.AuthToken {
		u.AuthToken = u2.AuthToken
		updated = true
	}
	if u.ID == config.App.OwnerID {
		u.Owner = true
		// force enabled if owner
		u.Enabled = true
	}
	u.p = u2
	return updated
}

func (u *User) Update(u2 *User) bool {
	updated := false
	if u.UUID != u2.UUID {
		u.UUID = u2.UUID
		updated = true
	}
	if u.Username != u2.Username {
		u.Username = u2.Username
		updated = true
	}
	if u.Email != u2.Email {
		u.Email = u2.Email
		updated = true
	}
	if u.Thumb != u2.Thumb {
		u.Thumb = u2.Thumb
		updated = true
	}
	if u.Owner != u2.Owner {
		u.Owner = u2.Owner
		updated = true
	}
	if !u.Owner && u.Enabled != u2.Enabled {
		u.Enabled = u2.Enabled
	}
	return updated
}

func (u *User) createPlexUser() {
	u.p = &plex.User{
		ID:        u.ID,
		UUID:      u.UUID,
		Username:  u.Username,
		Email:     u.Email,
		Thumb:     u.Thumb,
		AuthToken: u.AuthToken,
	}
}
