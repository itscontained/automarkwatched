package plex

import (
	"github.com/DirtyCajunRice/go-plex"

	v1 "github.com/itscontained/automarkwatched/api/v1"
	"github.com/itscontained/automarkwatched/internal/api"
)

var App *plex.App

func Setup(name, id string) {
	App = plex.New(id, plex.Product(name))
	api.A = App
	v1.App = App
}
