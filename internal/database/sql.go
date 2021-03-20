package database

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	v1 "github.com/itscontained/automarkwatched/api/v1"
	. "github.com/itscontained/automarkwatched/internal/config"
)

type database struct {
	Sqlx   *sqlx.DB
	Goqu   *goqu.Database
	Driver string
	Source string
}

var DB = &database{
	Driver: "sqlite3",
	Source: "db.sqlite",
}

func Init() {
	// create a new sqlx connection to the db
	x, err := sqlx.Connect(DB.Driver, DB.Source+"?_foreign_keys=1")
	if err != nil {
		log.Fatalln(err)
	}
	DB.Sqlx = x
	DB.Sqlx.MapperFunc(strcase.ToSnake)

	// create a goqu db object that uses the sqlx connection
	DB.Goqu = goqu.New("sqlite3", DB.Sqlx)

	// Create tables
	DB.Sqlx.MustExec(appSchema)
	DB.Sqlx.MustExec(userSchema)
	DB.Sqlx.MustExec(serverSchema)
	DB.Sqlx.MustExec(librarySchema)
	DB.Sqlx.MustExec(seriesSchema)
}

func Close() {
	if err := DB.Sqlx.Close(); err != nil {
		log.WithError(err).Error("problem gracefully closing database")
	}
}

func (db *database) GetAppConfig() error {
	var app Application
	found, err := db.Goqu.From("app").ScanStruct(&app)
	if err != nil {
		return err
	}
	if !found {
		_, err = db.Goqu.Insert("app").Rows(App).Executor().Exec()
		if err != nil {
			return err
		}
		log.WithField("identifier", App.Identifier).Debugf("saved new app identifier")
		return nil
	}
	App.Identifier = app.Identifier
	App.OwnerID = app.OwnerID
	return nil
}

// User Methods
func (db *database) GetOwnerID() (id int, err error) {
	_, err = db.Goqu.From("users").Select("id").Where(goqu.C("owner").Eq(true)).ScanVal(&id)
	return
}

func (db *database) GetUser(id int) *v1.User {
	var u v1.User
	ok, err := db.Goqu.From("users").Where(goqu.Ex{"id": id}).ScanStruct(&u)
	if !ok || err != nil {
		log.WithField("user_id", id).Error("could not get user from database")
		return nil
	}
	return &u
}

func (db *database) GetUsers() (map[int]*v1.User, error) {
	var users []*v1.User
	err := db.Goqu.From("users").ScanStructs(&users)
	if err != nil {
		return nil, err
	}
	u := make(map[int]*v1.User)
	for i := range users {
		u[users[i].ID] = users[i]
	}
	return u, nil
}

func (db *database) AddUser(user *v1.User) (err error) {
	_, err = db.Goqu.Insert("users").Rows(user).Executor().Exec()
	return
}

func (db *database) UpdateUser(user *v1.User) (err error) {
	_, err = db.Goqu.Update("users").Set(user).Where(goqu.Ex{"id": user.ID}).Executor().Exec()
	return
}

// server methods
func (db *database) GetServer(user *v1.User, machineIdentifier string) *v1.Server {
	var s v1.Server
	ex := goqu.Ex{"machine_identifier": machineIdentifier, "user_id": user.ID}
	ok, err := db.Goqu.From("servers").Where(ex).ScanStruct(&s)
	if !ok || err != nil {
		log.WithField("machine_identifier", machineIdentifier).Error("could not get server from database")
		return nil
	}
	return &s
}

func (db *database) GetServers(user *v1.User) (map[string]*v1.Server, error) {
	servers := make([]*v1.Server, 0)
	if err := db.Goqu.From("servers").Where(goqu.Ex{"user_id": user.ID}).ScanStructs(&servers); err != nil {
		return nil, err
	}
	serverMap := make(map[string]*v1.Server)
	for i := range servers {
		serverMap[servers[i].MachineIdentifier] = servers[i]
	}
	return serverMap, nil
}

func (db *database) AddServers(servers map[string]*v1.Server) error {
	var s []*v1.Server
	for i := range servers {
		s = append(s, servers[i])
	}
	_, err := db.Goqu.Insert("servers").Rows(s).Executor().Exec()
	if err != nil {
		return err
	}
	return nil
}

func (db *database) UpdateServer(server *v1.Server) {
	ex := goqu.Ex{"machine_identifier": server.MachineIdentifier, "user_id": server.UserID}
	_, err := db.Goqu.Update("servers").Set(server).Where(ex).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem updating server in database")
	}
}

// library methods
func (db *database) GetLibrary(user *v1.User, uuid string) *v1.Library {
	var library v1.Library
	ok, err := db.Goqu.From("libraries").Where(goqu.Ex{"uuid": uuid, "user_id": user.ID}).ScanStruct(&library)
	if !ok || err != nil {
		log.WithField("uuid", uuid).Error("could not get library from database")
		return nil
	}
	return &library
}

func (db *database) GetLibraries(user *v1.User) (map[string]*v1.Library, error) {
	libraries := make([]*v1.Library, 0)
	err := db.Goqu.From("libraries").Where(goqu.Ex{"user_id": user.ID}).ScanStructs(&libraries)
	if err != nil {
		return nil, err
	}
	libraryMap := make(map[string]*v1.Library)
	for i := range libraries {
		libraryMap[libraries[i].UUID] = libraries[i]
	}
	return libraryMap, nil
}

func (db *database) AddLibraries(libraries map[string]*v1.Library) error {
	var l []*v1.Library
	for i := range libraries {
		l = append(l, libraries[i])
	}
	_, err := db.Goqu.Insert("libraries").Rows(l).Executor().Exec()
	if err != nil {
		return err
	}
	return nil
}

func (db *database) UpdateLibrary(library *v1.Library) {
	ex := goqu.Ex{"uuid": library.UUID, "user_id": library.UserID}
	_, err := db.Goqu.Update("libraries").Set(library).Where(ex).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem updating libraries in database")
	}
	return
}

// series methods
func (db *database) GetOneSeries(user *v1.User, ratingKey int) *v1.Series {
	var series v1.Series
	ok, err := db.Goqu.From("series").Where(goqu.Ex{"rating_key": ratingKey, "user_id": user.ID}).ScanStruct(&series)
	if !ok || err != nil {
		log.WithField("rating_key", ratingKey).Error("could not get a series from database")
		return nil
	}
	return &series
}

func (db *database) AddSeries(series map[int]*v1.Series) error {
	for i := range series {
		_, err := db.Goqu.Insert("series").Rows(series[i]).Executor().Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *database) GetSeries(user *v1.User) (map[int]*v1.Series, error) {
	var series []*v1.Series
	err := db.Goqu.From("series").Where(goqu.Ex{"user_id": user.ID}).ScanStructs(&series)
	if err != nil {
		return nil, err
	}
	s := make(map[int]*v1.Series)
	for i := range series {
		s[series[i].RatingKey] = series[i]
	}
	return s, nil
}

func (db *database) GetSeriesScrobbled(user *v1.User) (map[int]*v1.Series, error) {
	var s []*v1.Series
	err := db.Goqu.From("series").Where(goqu.Ex{"user_id": user.ID, "scrobble": true}).ScanStructs(&s)
	if err != nil {
		return nil, err
	}
	series := make(map[int]*v1.Series)
	for i := range s {
		series[s[i].RatingKey] = s[i]
	}
	return series, nil
}

func (db *database) UpdateOneSeriesScrobble(user *v1.User, ratingKey int, scrobble bool) {
	where := goqu.Ex{"user_id": user.ID, "rating_key": ratingKey}
	_, err := db.Goqu.Update("series").Where(where).Set(goqu.Record{"scrobble": scrobble}).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem updating user series scrobble in database")
	}
}

func (db *database) UpdateSeriesScrobble(user *v1.User, ratingKeys []int, scrobble bool) error {
	for i := range ratingKeys {
		where := goqu.Ex{"user_id": user.ID, "rating_key": ratingKeys[i]}
		_, err := db.Goqu.Update("series").Where(where).Set(goqu.Record{"scrobble": scrobble}).Executor().Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *database) UpdateAppConfig() (err error) {
	_, err = db.Goqu.Update("app").Set(App).Where(goqu.Ex{"identifier": App.Identifier}).Executor().Exec()
	return
}

func (db *database) UpdateSeries(series *v1.Series) {
	_, err := db.Goqu.Update("series").Set(series).Where(goqu.Ex{"rating_key": series.RatingKey}).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem updating series in database")
	}
}
