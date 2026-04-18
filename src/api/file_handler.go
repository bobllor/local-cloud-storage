package api

import (
	"fmt"
	"net/http"
)

const (
	FileGetFileRoute = "/api/files"
)

type FileHandler struct {
}

type GetFileHandler struct {
}

// GetFiles retrieves a slice of Files based on the session ID and the given
// parent folder ID.
//
// If a parent folder ID is given, it will retrieve those parent folder files.
// If parent folder is nil, then it will retrieve the files with a nil parent.
//
// This requires the session ID.
func (g *GetFileHandler) GetFiles(w http.ResponseWriter, r *http.Request) {
	sesId := GetSessionFromCookie(r)
	if sesId == "" {
		err := fmt.Errorf("no cookie found for %s in files", CookieSessionKey)
		WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}
}
