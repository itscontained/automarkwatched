package config

import (
	"os"

	"github.com/google/uuid"
)

type Application struct {
	Identifier string `db:"identifier"`
	Name       string `db:"-"`
	APIPort    int    `db:"-"`
	WebPort    int    `db:"-"`
	OwnerID    int    `db:"owner_id"`
	Version    string `db:"version"`
}

var (
	App = Application{
		Name:    "AutoMarkWatched",
		APIPort: 5309,
		WebPort: 5310,
		Version: "0.0.0",
	}
	QuitChan = make(chan os.Signal, 1)
)

func init() {
	genUUID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	App.Identifier = genUUID.String()
}
