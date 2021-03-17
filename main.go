package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/itscontained/automarkwatched/internal/api"
	"github.com/itscontained/automarkwatched/internal/config"
	"github.com/itscontained/automarkwatched/internal/database"
	"github.com/itscontained/automarkwatched/internal/plex"
	"github.com/itscontained/automarkwatched/internal/routines"
	"github.com/itscontained/automarkwatched/internal/ui"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{PadLevelText: true})

	database.Init()

	if err := database.DB.GetAppConfig(); err != nil {
		log.Error(err)
	}

	plex.Setup(config.App.Name, config.App.Identifier)

	api.Start()
	ui.Start()
	routines.Start()
}
