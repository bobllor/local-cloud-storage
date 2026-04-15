package api

import (
	"encoding/json"
	"net/http"

	dbcon "github.com/bobllor/cloud-project/src/db_gateway"
	"github.com/bobllor/cloud-project/src/utils"
	"github.com/bobllor/gologger"
)

const (
	UserPostRegisterRoute = "POST /register"
	UserPostLoginRoute    = "POST /login"
)

const (
	CookieSessionKey = "LCSSessionID"
)

// TODO: add string checker for empty/invalid characters (username/password)
// TODO: add test cases for Login and other things

func NewUserHandler(gw *dbcon.Gateway, logger *gologger.Logger) *UserHandler {
	uh := &UserHandler{
		Post: PostUserHandler{
			Gateway: gw,
			deps:    utils.NewDeps(logger),
		},
		Get: GetUserHandler{
			Gateway: gw,
			deps:    utils.NewDeps(logger),
		},
	}

	return uh
}

// UserHandler contains handlers used for handling user related
// logic.
type UserHandler struct {
	Post PostUserHandler
	Get  GetUserHandler
	deps *utils.Deps
}

type GetUserHandler struct {
	Gateway *dbcon.Gateway
	deps    *utils.Deps
}

type PostUserHandler struct {
	Gateway *dbcon.Gateway
	deps    *utils.Deps
}

// LoginUser is the handler for handling the login and authentication.
//
// There are two authentication attempts:
//  1. A valid session ID is found in the cookie and validated
//  2. If there is no cookie, then the username and password will be validated
//
// If the user is successfully authenticated, the output of the response
// will contain the status and the session ID will be written to the cookie.
func (pu *PostUserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user RequestUserInfo

	c, err := r.Cookie(CookieSessionKey)
	if err != nil {
		pu.deps.Log.Infof("Cookie key %s not found", CookieSessionKey)
	} else {
		validSession, err := pu.Gateway.Session.ValidateSession(c.Value)
		if err != nil {
			pu.deps.Log.Warnf("Got an error while validating session: %v", err)
		}

		if validSession {
			// NOTE: this might be a bad idea. look at this another time
			err := WriteResponse(w, NewApiResponse(validSession))
			if err != nil {
				WriteErrorResponse(w, err, http.StatusInternalServerError)
			} else {
				return
			}
		}
	}

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	validUser, err := pu.Gateway.User.ValidateUser(user.Username, user.Password)
	if err != nil {
		pu.deps.Log.Criticalf("Error occurred during user validation: %v", err)
		WriteErrorResponse(w, err, http.StatusBadRequest)
	}

	pu.deps.Log.Debugf("Validate user result: %v", validUser)

	res := NewApiResponse(validUser)
	if validUser {
		// TODO: set cookies here
	}

	err = WriteResponse(w, res)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

// RegisterUser is the handler for registering users. A new session.Session struct
// will be marshalled as the response.
func (pu *PostUserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user RequestUserInfo

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	acc, err := pu.Gateway.User.AddUser(user.Username, user.Password)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	s, err := pu.Gateway.Session.UpsertSession(acc.AccountID)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	err = WriteResponse(w, s)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}
