package api

import "net/http"

type HandlerMap map[string]func(http.ResponseWriter, *http.Request)

type Handler interface{}

// httpError is a helper function used to write an error if one occurred.
// This does not check if err != nil. If err == nil then this will do nothing.
func httpError(w http.ResponseWriter, err error, statusCode int) {
	if err == nil {
		return
	}

	http.Error(w, err.Error(), statusCode)
}
