package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

var ErrUnauthorized = Error{
	Code:    http.StatusUnauthorized,
	Message: "unauthorized",
}

var ErrDatabase = Error{
	Code:    http.StatusInternalServerError,
	Message: "database read/write problem",
}

var ErrPlexAPI = Error{
	Code:    http.StatusBadGateway,
	Message: "plex api problem",
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, err Error) {
	log.Error(err)
	render.Status(r, err.Code)
	render.JSON(w, r, &err)
}

func ServerError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func SendError(w http.ResponseWriter, e error) {
	http.Error(w, e.Error(), http.StatusInternalServerError)
}

type ServerStatus struct {
	Date       time.Time `json:"date"`
	Configured bool      `json:"configured"`
}

func GetServerStatus(ownerId int) ServerStatus {
	configured := false
	if ownerId > 0 {
		configured = true
	}
	return ServerStatus{
		Date:       time.Now(),
		Configured: configured,
	}
}
