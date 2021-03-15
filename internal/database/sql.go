package database

import (
	"errors"
	"fmt"
	"reflect"

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
	DB.Sqlx.MustExec(userServerSchema)
	DB.Sqlx.MustExec(librarySchema)
	DB.Sqlx.MustExec(userLibrarySchema)
	DB.Sqlx.MustExec(seriesSchema)
	DB.Sqlx.MustExec(userSeriesSchema)
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

func (db *database) GetUsers() (users []v1.User, err error) {
	err = db.Goqu.From("users").ScanStructs(&users)
	return
}

func (db *database) AddUser(user v1.User) (err error) {
	_, err = db.Goqu.Insert("users").Rows(user).Executor().Exec()
	return
}

func (db *database) UpdateUser(user v1.User) (err error) {
	_, err = db.Goqu.Update("users").Set(user).Where(goqu.Ex{"id": user.ID}).Executor().Exec()
	return
}

// server methods
func (db *database) GetServer(machineIdentifier string) *v1.Server {
	var s v1.Server
	ok, err := db.Goqu.From("servers").Where(goqu.Ex{"machine_identifier": machineIdentifier}).ScanStruct(&s)
	if !ok || err != nil {
		log.WithField("machine_identifier", machineIdentifier).Error("could not get server from database")
		return nil
	}
	return &s
}

func (db *database) GetServers() (map[string]*v1.Server, error) {
	servers := make([]*v1.Server, 0)
	if err := db.Goqu.From("servers").ScanStructs(&servers); err != nil {
		return nil, err
	}
	serverMap := make(map[string]*v1.Server)
	for _, v := range servers {
		serverMap[v.MachineIdentifier] = v
	}
	return serverMap, nil
}

func (db *database) AddServers(servers map[string]*v1.Server) {
	var s []*v1.Server
	for i := range servers {
		s = append(s, servers[i])
	}
	_, err := db.Goqu.Insert("servers").Rows(s).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem adding servers to database")
	}
}

func (db *database) UpdateServer(server *v1.Server) {
	ex := goqu.Ex{"machine_identifier": server.MachineIdentifier}
	_, err := db.Goqu.Update("servers").Set(server).Where(ex).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem updating server in database")
	}
}

// user server methods
func (db *database) GetUserServer(user *v1.User, machineIdentifier string) *v1.UserServer {
	ex := goqu.Ex{"user_id": user.ID, "server_machine_identifier": machineIdentifier}
	var userServer v1.UserServer
	ok, err := db.Goqu.From("user_servers").Where(ex).ScanStruct(&userServer)
	if !ok || err != nil {
		log.WithFields(log.Fields{
			"user_id":                   user.ID,
			"server_machine_identifier": machineIdentifier,
		}).Error("could not get user server from database")
		return nil
	}
	return &userServer
}

func (db *database) GetUserServers(UserID int) map[string]*v1.UserServer {
	userServers := make([]*v1.UserServer, 0)
	err := db.Goqu.From("user_servers").Where(goqu.Ex{"user_id": UserID}).ScanStructs(&userServers)
	if err != nil {
		log.WithError(err).Error("problem getting user servers")
		return nil
	}
	userServerMap := make(map[string]*v1.UserServer)
	for i := range userServers {
		userServerMap[userServers[i].ServerMachineIdentifier] = userServers[i]
	}
	return userServerMap
}

func (db *database) AddUserServers(userServers map[string]*v1.UserServer) {
	var us []*v1.UserServer
	for i := range userServers {
		us = append(us, userServers[i])
	}
	_, err := db.Goqu.Insert("user_servers").Rows(us).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem adding user servers to database")
	}
}

func (db *database) UpdateUserServer(userServer v1.UserServer) (err error) {
	ex := goqu.Ex{"user_id": userServer.UserID, "server_machine_identifier": userServer.ServerMachineIdentifier}
	_, err = db.Goqu.Update("user_servers").Set(userServer).Where(ex).Executor().Exec()
	return
}

// library methods
func (db *database) GetLibrary(uuid string) *v1.Library {
	var library v1.Library
	ok, err := db.Goqu.From("libraries").Where(goqu.Ex{"uuid": uuid}).ScanStruct(&library)
	if !ok || err != nil {
		log.WithField("uuid", uuid).Error("could not get library from database")
		return nil
	}
	return &library
}

func (db *database) GetLibraries() map[string]*v1.Library {
	libraries := make([]*v1.Library, 0)
	err := db.Goqu.From("libraries").ScanStructs(&libraries)
	if err != nil {
		log.WithError(err).Error("problem getting libraries from database")
	}
	libraryMap := make(map[string]*v1.Library)
	for i := range libraries {
		libraryMap[libraries[i].UUID] = libraries[i]
	}
	return libraryMap
}

func (db *database) AddLibraries(libraries map[string]*v1.Library) {
	var l []*v1.Library
	for i := range libraries {
		l = append(l, libraries[i])
	}
	_, err := db.Goqu.Insert("libraries").Rows(l).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem adding library to database")
	}
	return
}

func (db *database) UpdateLibrary(library *v1.Library) {
	_, err := db.Goqu.Update("libraries").Set(library).Where(goqu.Ex{"uuid": library.UUID}).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem updating libraries in database")
	}
	return
}

// user library methods
func (db *database) GetUserLibrary(UserID int, machineIdentifier, uuid string) (userLibrary *v1.UserLibrary, ok bool, err error) {
	ex := goqu.Ex{"user_id": UserID, "server_machine_identifier": machineIdentifier, "library_uuid": uuid}
	ok, err = db.Goqu.From("user_libraries").Where(ex).ScanStruct(&userLibrary)
	return
}

func (db *database) GetUserLibraries(user *v1.User) map[string]*v1.UserLibrary {
	var ul []*v1.UserLibrary
	err := db.Goqu.From("user_libraries").Where(goqu.Ex{"user_id": user.ID}).ScanStructs(&ul)
	if err != nil {
		log.WithError(err).Error("problem getting user libraries from database")
		return nil
	}
	userLibraries := make(map[string]*v1.UserLibrary)
	for i := range ul {
		userLibraries[ul[i].LibraryUUID] = ul[i]
	}
	return userLibraries
}

func (db *database) AddUserLibraries(userlibraries map[string]*v1.UserLibrary) {
	var ul []*v1.UserLibrary
	for i := range userlibraries {
		ul = append(ul, userlibraries[i])
	}
	_, err := db.Goqu.Insert("user_libraries").Rows(ul).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem adding user library to database")
	}
}

func (db *database) UpdateUserLibrary(userLibrary *v1.UserLibrary) {
	ex := goqu.Ex{
		"user_id":                   userLibrary.UserID,
		"server_machine_identifier": userLibrary.ServerMachineIdentifier,
		"library_uuid":              userLibrary.LibraryUUID,
	}
	_, err := db.Goqu.Update("user_libraries").Set(userLibrary).Where(ex).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem updating user library in database")
	}
}

// series methods
func (db *database) GetOneSeries(ratingKey int) *v1.Series {
	var series v1.Series
	ok, err := db.Goqu.From("series").Where(goqu.Ex{"rating_key": ratingKey}).ScanStruct(&series)
	if !ok || err != nil {
		log.WithField("rating_key", ratingKey).Error("could not get a series from database")
		return nil
	}
	return &series
}

func (db *database) GetSeries() map[int]*v1.Series {
	var series []*v1.Series
	err := db.Goqu.From("series").ScanStructs(&series)
	if err != nil {
		log.WithError(err).Error("problem getting series in database")
		return nil
	}
	s := make(map[int]*v1.Series)
	for i := range series {
		s[series[i].RatingKey] = series[i]
	}
	return s
}

func (db *database) GetUserSeries(user *v1.User) map[int]*v1.UserSeries {
	var userSeries []*v1.UserSeries
	err := db.Goqu.From("user_series").Where(goqu.Ex{"user_id": user.ID}).ScanStructs(&userSeries)
	if err != nil {
		log.WithError(err).Error("problem getting user series from database")
		return nil
	}
	us := make(map[int]*v1.UserSeries)
	for i := range userSeries {
		us[userSeries[i].SeriesRatingKey] = userSeries[i]
	}
	return us
}

func (db *database) GetUserSeriesScrobbled(id int) (userSeries []v1.UserSeries, err error) {
	err = db.Goqu.From("user_series").Where(goqu.Ex{"user_id": id, "scrobble": true}).ScanStructs(&userSeries)
	return
}

func (db *database) AddUserSeries(userSeries map[int]*v1.UserSeries) {
	for i := range userSeries {
		_, err := db.Goqu.Insert("user_series").Rows(userSeries[i]).Executor().Exec()
		if err != nil {
			log.WithError(err).Error("problem adding user series to database")
		}
	}
}

func (db *database) UpdateOneUserSeriesScrobble(user *v1.User, ratingKey int, scrobble bool) {
	where := goqu.Ex{"user_id": user.ID, "series_rating_key": ratingKey}
	_, err := db.Goqu.Update("user_series").Where(where).Set(goqu.Record{"scrobble": scrobble}).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem updating user series scrobble in database")
	}
}

func (db *database) Update(table string, s interface{}) (err error) {
	v := reflect.ValueOf(s)
	fmt.Printf("%#v\n", v)
	q, _, _ := db.Goqu.Update("app").Set(v).ToSQL()
	log.Print(q)
	_, err = db.Goqu.Update(table).Set(v).Executor().Exec()
	return
}

func (db *database) UpdateToSQL(table string, s interface{}) (statement string) {
	statement, _, _ = db.Goqu.Update(table).Set(s).ToSQL()
	return
}

func (db *database) UpdateAppConfig() (err error) {
	_, err = db.Goqu.Update("app").Set(App).Where(goqu.Ex{"identifier": App.Identifier}).Executor().Exec()
	return
}

func (db *database) AddSeries(series map[int]*v1.Series) {
	for i := range series {
		_, err := db.Goqu.Insert("series").Rows(series[i]).Executor().Exec()
		if err != nil {
			log.WithError(err).Error("problem adding series to database")
		}
	}
}

func (db *database) UpdateSeries(series *v1.Series) {
	_, err := db.Goqu.Update("series").Set(series).Where(goqu.Ex{"ratingKey": series.RatingKey}).Executor().Exec()
	if err != nil {
		log.WithError(err).Error("problem updating series in database")
	}
}

func (db *database) UpdateUserSeries(user v1.User, series []v1.Series) (err error) {
	errorCount := 0
	for _, s := range series {
		err = db.Update("user_series", s)
		if err != nil {
			log.WithError(err).Error("could not update scrobble state")
			errorCount++
		}
	}
	if errorCount > 0 {
		err = errors.New(fmt.Sprintf("could not update %d/%d scrobble states", errorCount, len(series)))
	}
	return
}
