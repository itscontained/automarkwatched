package v1

import (
	"errors"

	"github.com/DirtyCajunRice/go-plex"

	"github.com/itscontained/automarkwatched/internal/config"
	px "github.com/itscontained/automarkwatched/internal/plex"
)

var (
	ErrNotApplicationOwned = errors.New("not application owner server")
	ErrNoUserAuthToken     = errors.New("auth token not set")
	ErrNoUserAccessToken   = errors.New("access token not set")
	ErrNoAppOwner          = errors.New("app owner not set")
	App = px.App
)

// User is a wrapper struct for a plex.User
type User struct {
	*plex.User
	Owner   bool               		`json:"owner" db:"owner" goqu:"skipupdate"`
	Exists  bool               		`json:"exists" db:"-"`
	Enabled bool               		`json:"enabled" db:"enabled"`
	Servers map[string]*UserServer 		`json:"-" db:"-"`
	Libraries map[string] *UserLibrary `json:"-" db:"-"`
	Series map[int]*UserSeries `json:"-" db:"-"`
}

func (u *User) SyncUser() error {
	if u.AuthToken == "" {
		return ErrNoUserAuthToken
	}
	user, err := App.User(u.AuthToken)
	if err != nil {
		return err
	}
	u.User = user
	return nil
}

func (u *User) SyncServers() error {
	if u.AuthToken == "" {
		return ErrNoUserAuthToken
	}
	servers, err := u.User.Servers()
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
		u.Servers = make(map[string]*UserServer)
	}
	if _, ok := u.Servers[server.MachineIdentifier]; !ok {
		u.Servers[server.MachineIdentifier] = &UserServer{
			UserID:                  u.ID,
			ServerMachineIdentifier: server.MachineIdentifier,
			AccessToken:             server.AccessToken,
			Enabled:                 true,
			User:                    u,
			Server:                  ,
		}
		return nil
	}
	u.Servers[server.MachineIdentifier].comparePlexServer(server)
	return nil
}

func (u *User) GetRecursive() error {
	err := u.SyncUser()
	if err != nil {
		return err
	}
	err = u.SyncServers()
	if err != nil {
		return err
	}

	for i := range u.Servers {
		err = u.Servers[i].SyncLibraries()
		if err != nil {
			return err
		}
		for j := range u.Servers[i].Libraries {
			err = u.Servers[i].Libraries[j].SyncSeries()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *User) GetSeries() (map[int]*Series, error) {
	series := make(map[int]*Series)
	for i := range u.Servers {
		err := u.Servers[i].SyncLibraries()
		if err != nil {
			return nil, err
		}
		for j := range u.Servers[i].Libraries {
			err = u.Servers[i].Libraries[j].SyncSeries()
			if err != nil {
				return nil, err
			}
			for k := range u.Servers[i].Libraries[j].Series {
				series[k] = u.Servers[i].Libraries[j].Series[k]
			}
		}
	}
	return series, nil
}

func (u *User) MergeUserSeries(userSeries map[int]*UserSeries) (map[int]*UserSeries, error) {
	series, err := u.GetSeries()
	if err != nil {
		return nil, err
	}
	for s := range userSeries {
		userSeries[s].Series = series[s]

	}
	return userSeries, nil
}

func (u *User) MergeUserServerAccess(userServers map[string]*UserServer) {
	for machineIdentifier := range userServers {
		u.Servers[machineIdentifier].AccessToken = userServers[machineIdentifier].AccessToken
	}
}
