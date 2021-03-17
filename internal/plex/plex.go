package plex

import "github.com/DirtyCajunRice/go-plex"

var App *plex.App

func Setup(name, id string) {
	App = plex.New(id, plex.Product(name))
}
