package api

import (
	"errors"
	"fmt"
	"net/http"

	dbgateway "github.com/bobllor/cloud-project/src/db_gateway"
	"github.com/bobllor/cloud-project/src/utils"
	"github.com/bobllor/gologger"
)

const (
	FileGetFileParentRoute = "GET /storage/folder/{parentID}"
	FileGetFileRootRoute   = "GET /storage"
)

type FileHandler struct {
	gateway *dbgateway.Gateway
	deps    *utils.Deps
}

func NewFileHandler(gw *dbgateway.Gateway, logger *gologger.Logger) *FileHandler {
	return &FileHandler{
		gateway: gw,
		deps:    utils.NewDeps(logger),
	}
}

// GetFiles retrieves a slice of Files based on the session ID and the given
// parent folder ID.
//
// If a parent folder ID is given, it will retrieve those parent folder files.
// If parent folder is nil, then it will retrieve the files with a nil parent.
//
// This requires the auth middleware session ID.
func (fh *FileHandler) GetFiles(w http.ResponseWriter, r *http.Request) {
	parentKey := "parentID"

	parentID := r.PathValue(parentKey)
	fh.deps.Log.Debugf("Request query: %v", parentID)

	sesID := GetSessionFromCookie(r)
	if sesID == "" {
		fh.deps.Log.Info("No cookie found with request")
		err := fmt.Errorf("no cookie found for %s in files", CookieSessionKey)
		WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	files, err := fh.gateway.File.GetFilesBySessionAndParentFolder(sesID, parentID)
	if err == dbgateway.FileDoesNotExistErr {
		fh.deps.Log.Infof("Given parent ID %s does not exist: %v", parentID, err)
		WriteErrorResponse(w, errors.New("invalid parent ID"), http.StatusBadRequest)
		return
	}
	if err != nil {
		fh.deps.Log.Criticalf("Failed to retrieve files with session ID and parent folder ID: %v", err)
		WriteErrorResponse(w, errors.New("an unknwon error has occurred"), http.StatusInternalServerError)
		return
	}

	res := NewApiResponse(files)
	n, err := WriteResponse(w, res)
	if err != nil {
		fh.deps.Log.Criticalf("Failed to write response to client: %v", err)
		WriteErrorResponse(w, errors.New("an unknown error has occurred"), http.StatusInternalServerError)
		return
	}

	fh.deps.Log.Debugf("Response bytes: %d", n)
}
