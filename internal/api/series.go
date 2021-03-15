package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"

	v1 "github.com/itscontained/automarkwatched/api/v1"
)

func seriesRoutes(r chi.Router) {
	r.Get("/series/sync", syncSeries)
	r.Patch("/series/scrobble", updateSeriesScrobble)
	r.Patch("/series/scrobble/{seriesRatingKey}", updateOneSeriesScrobble)
}

func syncSeries(w http.ResponseWriter, r *http.Request) {
	user := GetContextUser(r)
	l := log.WithField("endpoint", "sync")
	l.Info("starting sync")
	userLibraries := getUserLibraries(user)

	savedSeries := db.GetSeries()
	missing := make(map[int]*v1.Series)
	allSeries := make(map[int]*v1.Series)
	modified := 0
	for i := range userLibraries {
		if err := userLibraries[i].SyncSeries(); err != nil {
			log.WithError(err).Error("sync series problem")
			return
		}
		for s := range userLibraries[i].Library.Series {
			rk := userLibraries[i].Library.Series[s].RatingKey
			allSeries[rk] = userLibraries[i].Library.Series[s]
			allSeries[rk].ServerMachineIdentifier = userLibraries[i].ServerMachineIdentifier
			if _, ok := savedSeries[rk]; !ok {
				missing[rk] = userLibraries[i].Library.Series[s]
				missing[rk].Library.Server = user.Servers[userLibraries[i].ServerMachineIdentifier]
				continue
			}
			if !compareSeries(userLibraries[i].Library.Series[s], savedSeries[rk]) {
				savedSeries[rk].Series = userLibraries[i].Library.Series[s].Series
				db.UpdateSeries(savedSeries[rk])
				modified++
			}
		}
	}
	if len(missing) > 0 {
		db.AddSeries(missing)
	}
	untouched := len(allSeries) - len(missing) - modified
	msg := fmt.Sprintf("Untouched: %d, Added: %d, Modified: %d", untouched, len(missing), modified)
	resp := Response{
		Code:    http.StatusOK,
		Message: msg,
	}
	render.JSON(w, r, &resp)
	go setUserSeries(user, allSeries)
}

func setUserSeries(user *v1.User, series map[int]*v1.Series) {
	savedUserSeries := GetUserSeries(user)

	missing := make(map[int]*v1.UserSeries)
	for i := range series {
		if _, ok := savedUserSeries[series[i].RatingKey]; !ok {
			missing[series[i].RatingKey] = &v1.UserSeries{
				UserID:                  user.ID,
				ServerMachineIdentifier: series[i].Library.Server.MachineIdentifier,
				LibraryUUID:             series[i].Library.UUID,
				SeriesRatingKey:         series[i].RatingKey,
				Scrobble:                false,
				Enabled:                 true,
				User:                    user,
				Series:                  series[i],
			}
		}
	}

	if len(missing) > 0 {
		db.AddUserSeries(missing)
	}
}

func GetUserSeries(user *v1.User) map[int]*v1.UserSeries {
	userSeries := db.GetUserSeries(user)
	if userSeries == nil {
		return nil
	}
	mappedUserSeries := make(map[int]*v1.UserSeries)
	for i := range userSeries {
		userSeries[i].Series = db.GetOneSeries(i)
		userSeries[i].User = user
		mappedUserSeries[i] = userSeries[i]
	}
	return mappedUserSeries
}

func compareSeries(s1, s2 *v1.Series) bool {
	mismatched := make([]string, 0)
	if s1.Title != s2.Title {
		mismatched = append(mismatched, "title")
	}
	if s1.RatingKey != s2.RatingKey {
		mismatched = append(mismatched, "agent")
	}
	if s1.Enabled != s2.Enabled {
		mismatched = append(mismatched, "address")
	}
	if len(mismatched) > 0 {
		return false
	}
	return true
}

func GetAll(user *v1.User) map[int]*v1.UserSeries {
	userSeries := GetUserSeries(user)
	userLibraries := getUserLibraries(user)
	for i := range userSeries {
		userSeries[i].Series.Library = userLibraries[userSeries[i].LibraryUUID].Library
		userSeries[i].Series.Library.Series[userSeries[i].SeriesRatingKey] = userSeries[i].Series
	}
	return userSeries
}

func updateOneSeriesScrobble(w http.ResponseWriter, r *http.Request) {
	user := GetContextUser(r)
	ratingKeyString := chi.URLParam(r, "seriesRatingKey")
	if ratingKeyString == "" {
		log.Error("no key defined in url param")
		return
	}
	ratingKey, err := strconv.Atoi(ratingKeyString)
	if err != nil {
		log.Error(err)
		return
	}
	scrobbleString := r.URL.Query().Get("scrobble")
	scrobble, e := strconv.ParseBool(scrobbleString)
	if e != nil {
		log.Error(e)
		return
	}
	db.UpdateOneUserSeriesScrobble(user, ratingKey, scrobble)
	render.NoContent(w, r)
}

func updateSeriesScrobble(w http.ResponseWriter, r *http.Request) {
	var seriesIDs []int
	if err := render.DecodeJSON(r.Body, &seriesIDs); err != nil {
		log.WithError(err).Error("problem parsing json")
		return
	}
}
