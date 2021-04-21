package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"

	v1 "github.com/itscontained/automarkwatched/api/v1"
)

func seriesRoutes(r chi.Router) {
	r.Get("/series", getSeries)
	r.Post("/series", syncSeries)
	r.Patch("/series", scrobbleSeries)
}

func getSeries(w http.ResponseWriter, r *http.Request) {
	user := GetContextUser(r)
	if err := PullSeries(user); err != nil {
		ServerError(w)
		return
	}
	render.JSON(w, r, &user.Libraries)
}

func PullSeries(user *v1.User) error {
	series, err := db.GetSeries(user)
	if err != nil {
		return err
	}
	if user.Series == nil {
		user.Series = make(map[int]*v1.Series)
	}
	for i := range series {
		user.AttachSeries(series[i])
	}
	return nil
}

func syncSeries(w http.ResponseWriter, r *http.Request) {
	user := GetContextUser(r)
	l := log.WithField("endpoint", "sync")
	l.Info("starting sync")

	if err := SyncAll(user); err != nil {
		SendError(w, err)
		return
	}

	msg, err := SaveSeries(user)
	if err != nil {
		SendError(w, err)
		return
	}
	render.JSON(w, r, &msg)
}

func SaveSeries(user *v1.User) (map[string]int, error) {
	missing := make(map[int]*v1.Series)
	ownerSeries, err := db.GetOwnerSeries()
	if err != nil {
		return nil, err
	}
	series, err := db.GetSeries(user)
	if err != nil {
		return nil, err
	}
	for i := range ownerSeries {
		if _, ok := series[i]; !ok {
			ownerSeries[i].Enabled = true
			ownerSeries[i].UserID = user.ID
			missing[i] = ownerSeries[i]
		}
	}
	if len(missing) > 0 {
		if err = db.AddSeries(missing); err != nil {
			return nil, err
		}
	}
	untouched := len(user.Series) - len(missing)
	msg := map[string]int{
		"Untouched": untouched,
		"Added":     len(missing),
	}
	return msg, nil
}

func scrobbleSeries(w http.ResponseWriter, r *http.Request) {
	user := GetContextUser(r)
	scrobbleString := r.URL.Query().Get("scrobble")
	scrobble, err := strconv.ParseBool(scrobbleString)
	if err != nil {
		SendError(w, err)
		return
	}
	ratingKeyString := r.URL.Query().Get("ratingKey")
	ratingKeyStringSlice := strings.Split(ratingKeyString, ",")
	ratingKeys := make([]int, 0)
	for i := range ratingKeyStringSlice {
		if ratingKeyStringSlice[i] == "" {
			continue
		}
		var v int
		v, err = strconv.Atoi(ratingKeyStringSlice[i])
		if err != nil {
			SendError(w, err)
			return
		}
		ratingKeys = append(ratingKeys, v)
	}

	err = db.UpdateSeriesScrobble(user, ratingKeys, scrobble)
	if err != nil {
		log.Error(err)
		SendError(w, err)
		return
	}
}

func SyncAll(user *v1.User) error {
	if err := SyncAllButSeries(user); err != nil {
		return err
	}
	for i := range user.Libraries {

		if err := user.Libraries[i].SyncSeries(); err != nil {
			return err
		}
	}
	return nil
}

func SyncAllButSeries(user *v1.User) error {
	if err := user.SyncUser(); err != nil {
		return err
	}
	if err := user.SyncServers(); err != nil {
		return err
	}
	for i := range user.Servers {
		if err := user.Servers[i].SyncLibraries(); err != nil {
			return err
		}
	}
	return nil
}

func PullAll(user *v1.User) error {
	if err := PullServers(user); err != nil {
		return err
	}
	if err := PullLibraries(user); err != nil {
		return err
	}
	if err := PullSeries(user); err != nil {
		return err
	}
	return nil
}
