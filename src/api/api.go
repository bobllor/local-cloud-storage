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

type ApiHandler struct {
	UserHandler    *UserHandler
	FileHandler    *FileHandler
	SessionHandler *SessionHandler
	gateway        *dbcon.Gateway
	log            *gologger.Logger
}

// NewApiHandler creates a new Api struct.
func NewApiHandler(gw *dbcon.Gateway, logger *gologger.Logger) *ApiHandler {
	api := &ApiHandler{
		UserHandler:    NewUserHandler(gw, logger),
		FileHandler:    NewFileHandler(gw, logger),
		SessionHandler: NewSessionHandler(gw, logger),
		gateway:        gw,
		log:            logger,
	}

	return api
}

// RequestMiddleware wraps a function with a middleware used to log and write headers by default.
//
// This does not handle auth, use ah.CreateAuthMiddleware for auth based middleware.
func (ah *ApiHandler) CreateRequestMiddleware(f func(http.ResponseWriter, *http.Request)) http.Handler {
	next := http.HandlerFunc(f)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: add more logging and things yeah.
		ah.middlewareLogger(r)
		WriteHeaders(w, r)

		next.ServeHTTP(w, r)
	})
}

// CreateAuthHandler creates a new handler from a given handler function, wrapped in an
// authentication-based closure.
//
// Headers are automatically written if wrapped with this method.
func (ah *ApiHandler) CreateAuthMiddleware(f func(http.ResponseWriter, *http.Request)) http.Handler {
	next := http.HandlerFunc(f)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ah.middlewareLogger(r)

		sessionCookie, err := r.Cookie(CookieSessionKey)
		WriteHeaders(w, r)
		if err != nil {
			newErr := errors.New("unauthorized access")
			ah.log.Infof("Unauthorized access, no cookie found for %v", r.RemoteAddr)
			WriteErrorResponse(w, newErr, http.StatusUnauthorized)

			return
		}

		validSession, err := ah.gateway.Session.ValidateSession(sessionCookie.Value)
		if err != nil {
			ah.log.Criticalf("Validating session database query failed: %v", err)
			WriteErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		if !validSession {
			sessionErr := errors.New("session ID is invalid")
			ah.log.Infof("Invalid session ID, failed validation for %v", r.RemoteAddr)
			WriteErrorResponse(w, sessionErr, http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	})
}

// middlewareLogger is a helper function used to log the request information.
// This is used for general logging of the request.
func (ah *ApiHandler) middlewareLogger(r *http.Request) {
	ah.log.Infof("%s: accessed on agent %s", r.RemoteAddr, r.UserAgent())
}
