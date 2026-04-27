package api

import (
	"encoding/json"
	"net/http"

	"github.com/bobllor/cloud-project/src/session"
)

const (
	CookieSessionKey = "lcsSessionID"
)

type HandlerMap map[string]func(http.ResponseWriter, *http.Request)

type Handler interface{}

// WriteErrorResponse is a helper function used to write an error to
// the ResponseWriter.
func WriteErrorResponse(w http.ResponseWriter, msg string, statusCode int, reason ReasonCode) {
	errRes := NewApiResponseError(statusCode, msg, reason)

	b, err := json.Marshal(errRes)
	// if err is not nil, then default to a basic value.
	// TODO: maybe find a fix for this. remove this line later
	if err != nil {
		http.Error(w, err.Error(), statusCode)
	} else {
		http.Error(w, string(b), statusCode)
	}

}

// WriteResponse writes a response to the ResponseWriter. If an error occurs
// while writing the response, it will return the error and the response will not
// be written.
//
// The bytes written will be returned.
func WriteResponse(w http.ResponseWriter, v any) (int, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}

	i, err := w.Write(b)
	if err != nil {
		return 0, err
	}

	return i, nil
}

// WriteHeaders writes the headers for CORS.
func WriteHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Vary", "Origin")
}

// GetSessionFromCookie retrieves the session ID from the request headers.
// If the cookie does not exist, then it will return an empty string.
func GetSessionFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(CookieSessionKey)
	if err != nil {
		return ""
	}

	return cookie.Value
}

// SetCookieSession sets the cookie for the session.
// If the session already exists in the cookie, then it will overwrite the
// cookie's value.
func SetCookieSession(w http.ResponseWriter, s *session.Session) {
	c := http.Cookie{
		Name:     CookieSessionKey,
		Value:    s.SessionID,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, &c)
}

// SetCookie sets the cookie with the given key and value to the headers.
func SetCookie(w http.ResponseWriter, key string, value string) {

}
