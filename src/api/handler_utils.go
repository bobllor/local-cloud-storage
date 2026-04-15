package api

import (
	"encoding/json"
	"net/http"
)

type HandlerMap map[string]func(http.ResponseWriter, *http.Request)

type Handler interface{}

// WriteErrorResponse is a helper function used to write an error if one occurred to
// the ResponseWriter.
// This does not check if err != nil. If err == nil then this will do nothing.
func WriteErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	if err == nil {
		return
	}

	errRes := NewApiResponseError(statusCode, err.Error())

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
func WriteResponse(w http.ResponseWriter, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}
