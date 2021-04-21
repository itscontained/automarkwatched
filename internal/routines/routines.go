package routines

import (
	log "github.com/sirupsen/logrus"

	"github.com/robfig/cron/v3"

	v1 "github.com/itscontained/automarkwatched/api/v1"
	"github.com/itscontained/automarkwatched/internal/api"
	. "github.com/itscontained/automarkwatched/internal/config"
	"github.com/itscontained/automarkwatched/internal/database"
	"github.com/itscontained/automarkwatched/internal/ui"
)

var (
	db   = database.DB
	Cron = cron.New()
)

func Start() {
	_, err := Cron.AddFunc("* * * * *", Scrobble)
	if err != nil {
		log.Error(err)
	}
	Cron.Start()
	for {
		select {
		case <-QuitChan:
			Cron.Stop()
			ui.Stop()
			api.Stop()
			database.Close()
			return
		}
	}
}

func Scrobble() {
	l := log.WithField("job", "Scrobble")
	l.Info("starting job")
	users, err := db.GetUsers()
	if err != nil {
		l.WithError(err).Error("problem getting users")
	}
	for i := range users {
		userSeriesScrobble := make(map[int]*v1.Series)
		userSeriesScrobble, err = db.GetSeriesScrobbled(users[i])
		if err != nil {
			l.WithError(err).Error("problem getting user series scrobbled")
			return
		}
		if err = api.PullServers(users[i]); err != nil {
			l.WithError(err).Error("problem getting user servers")
		}
		for j := range users[i].Servers {
			if err = users[i].Servers[j].SyncLibraries(); err != nil {
				log.WithError(err).Error("problem syncing libraries")
				return
			}
		}
		for j := range users[i].Libraries {
			unwatched, e := users[i].Libraries[j].Unwatched()
			if e != nil {
				l.WithError(e).Error("problem getting unwatched series")
				return
			}
			for _, k := range unwatched {
				if _, ok := userSeriesScrobble[k.RatingKey]; ok {
					if !userSeriesScrobble[k.RatingKey].Scrobble {
						log.Error("i should not have gotten an unscrobbled record...")
						continue
					}
					err = users[i].Servers[userSeriesScrobble[k.RatingKey].ServerID].Scrobble(k.RatingKey)
					log.Infof("scrobbled episodes for %s", k.Title)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
	l.Info("job finished")
}
