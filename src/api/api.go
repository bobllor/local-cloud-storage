package api

import (
	"errors"
	"net/http"

	dbcon "github.com/bobllor/cloud-project/src/db_gateway"
	"github.com/bobllor/gologger"
)

const (
	ContentJson = "application/json"
)

const (
	CookieSessionKey = "lcsSessionID"
)

type ApiHandler struct {
	UserHandler *UserHandler
	gateway     *dbcon.Gateway
}

// NewApiHandler creates a new Api struct.
func NewApiHandler(gw *dbcon.Gateway, logger *gologger.Logger) *ApiHandler {
	api := &ApiHandler{
		UserHandler: NewUserHandler(gw, logger),
		gateway:     gw,
	}

	return api
}

// CreateAuthHandler creates a new handler from a given handler function, wrapped in an
// authentication-based closure.
func (ah *ApiHandler) CreateAuthHandler(f func(http.ResponseWriter, *http.Request)) http.Handler {
	next := http.HandlerFunc(f)

	// TODO: add logging
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(CookieSessionKey)
		WriteHeaders(w, r)
		if err != nil {
			newErr := errors.New("unauthorized")
			WriteErrorResponse(w, newErr, http.StatusUnauthorized)

			return
		}

		validSession, err := ah.gateway.Session.ValidateSession(sessionCookie.Value)
		if err != nil {
			WriteErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		if !validSession {
			sessionErr := errors.New("session ID is invalid")
			WriteErrorResponse(w, sessionErr, http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	})
}
