package routines

import (
	log "github.com/sirupsen/logrus"

	"github.com/robfig/cron/v3"

	"github.com/itscontained/automarkwatched/internal/api"
	. "github.com/itscontained/automarkwatched/internal/config"
	"github.com/itscontained/automarkwatched/internal/database"
	"github.com/itscontained/automarkwatched/internal/ui"
	"github.com/itscontained/automarkwatched/pkg/provider/plex"
)

var (
	db   = database.DB
	Cron = cron.New()
)

func Start() {
	_, err := Cron.AddFunc("*/5 * * * *", Scrobble)
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
		userSeries, e := db.GetUserSeriesScrobbled(users[i].ID)
		if e != nil {
			l.WithError(err).Error("problem getting user series scrobbled")
			continue
		}
		data := make(map[string]map[string]string)
		data2 := make(map[string]int)
		allUnwatchedSeries := make([]plex.Series, 0)
		for s := range userSeries {

			if !userSeries[s].Scrobble {
				log.Error("i should not have gotten an unscrobbled record...")
				continue
			}
			if _, ok := data[userSeries[s].ServerMachineIdentifier]; !ok {
				data[userSeries[s].ServerMachineIdentifier] = make(map[string]string)
				us := db.GetUserServer(&users[i], userSeries[s].ServerMachineIdentifier)
				data[userSeries[s].ServerMachineIdentifier]["access_token"] = us.AccessToken
				gs := db.GetServer(userSeries[s].ServerMachineIdentifier)
				data[userSeries[s].ServerMachineIdentifier]["url"] = gs.URL()
			}
			if _, ok := data2[userSeries[s].LibraryUUID]; !ok {
				li := db.GetLibrary(userSeries[s].LibraryUUID)
				data2[userSeries[s].LibraryUUID] = li.Key
				tmpSeries, err := plex.GetTVSeries(
					data[userSeries[s].ServerMachineIdentifier]["url"],
					data[userSeries[s].ServerMachineIdentifier]["access_token"],
					data2[userSeries[s].LibraryUUID],
					true,
				)
				if err != nil {
					log.WithError(err).Error("couldnt pull tv series from plex")
					continue
				}
				allUnwatchedSeries = append(allUnwatchedSeries, tmpSeries...)
			}

		}
		for s := range userSeries {
			isScrobbled := false
			for d := range allUnwatchedSeries {
				if userSeries[s].SeriesRatingKey == allUnwatchedSeries[d].RatingKey {
					if userSeries[s].Scrobble {
						log.Printf("marking %d for  %s watched", allUnwatchedSeries[d].ChildCount, allUnwatchedSeries[d].Title)
						isScrobbled = true
					}
					break
				}

			}
			if !isScrobbled {
				continue
			}
			plex.Scrobble(
				data[userSeries[s].ServerMachineIdentifier]["url"],
				data[userSeries[s].ServerMachineIdentifier]["access_token"],
				userSeries[s].SeriesRatingKey,
			)
		}
	}
}
