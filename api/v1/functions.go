package v1

import "github.com/itscontained/automarkwatched/pkg/provider/plex"

func NewUser(authToken string) *User {
	return &User{
		User: &plex.User{
			AuthToken: authToken,
		},
		Owner:   false,
		Enabled: true,
		Exists:  false,
		Servers: make(map[string]*Server),
	}
}
