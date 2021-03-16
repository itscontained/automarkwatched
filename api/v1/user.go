package v1

import (
	"errors"

	"github.com/itscontained/automarkwatched/internal/config"
	"github.com/itscontained/automarkwatched/pkg/provider/plex"
)

var (
	ErrNotApplicationOwned = errors.New("not application owner server")
	ErrNoUserAuthToken     = errors.New("auth token not set")
	ErrNoUserAccessToken   = errors.New("access token not set")
	ErrNoAppOwner          = errors.New("app owner not set")
)

// User is a wrapper struct for a plex.User
type User struct {
	*plex.User
	Owner   bool               `json:"owner" db:"owner" goqu:"skipupdate"`
	Exists  bool               `json:"exists" db:"-"`
	Enabled bool               `json:"enabled" db:"enabled"`
	Servers map[string]*Server `json:"-" db:"-"`
}

func (u *User) SyncUser() error {
	user, err := plex.GetUser(u.AuthToken)
	if err != nil {
		return err
	}
	u.User = &user
	return nil
}

func (u *User) SyncServers() error {
	if u.AuthToken == "" {
		return ErrNoUserAuthToken
	}
	servers, err := plex.GetServers(u.AuthToken)
	if err != nil {
		return err
	}
	for _, server := range servers {
		err = u.AttachPlexServer(&server)
		if err != nil || err == ErrNoAppOwner {
			return err
		}
	}
	return nil
}

func (u *User) AttachServer(server *Server) error {
	err := u.AttachPlexServer(server.Server)
	if err != nil && err != ErrNoAppOwner {
		return err
	}
	u.Servers[server.MachineIdentifier].Enabled = server.Enabled
	return nil
}

func (u *User) AttachPlexServer(server *plex.Server) error {
	if config.App.OwnerID == 0 {
		return ErrNoAppOwner
	}
	if u.ID == config.App.OwnerID && server.OwnerId == 0 {
		server.OwnerId = u.ID
	}
	if server.OwnerId != config.App.OwnerID {
		return ErrNotApplicationOwned
	}
	if u.Servers == nil {
		u.Servers = make(map[string]*Server)
	}
	if _, ok := u.Servers[server.MachineIdentifier]; !ok {
		u.Servers[server.MachineIdentifier] = &Server{
			Server:  server,
			Enabled: true,
		}
		return nil
	}
	u.Servers[server.MachineIdentifier].comparePlexServer(server)
	return nil
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
