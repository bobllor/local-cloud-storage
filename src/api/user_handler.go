package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	dbcon "github.com/bobllor/cloud-project/src/db_gateway"
	"github.com/bobllor/cloud-project/src/utils"
	"github.com/bobllor/gologger"
)

const (
	UserPostRegisterRoute = "POST /api/register"
	UserPostLoginRoute    = "POST /api/login"
	UserGetUserRoute      = "/api/user"
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

// GetUserBySessionID retrieves the user information to the response. This uses
// the cookie in the headers.
//
// This requires a valid session ID in order to retrieve the user. The user
// will only include the necessary information that is required with the frontend.
//
// If the user cannot be found, then the output will be nil in the Response.
// An error Response will only be returned for internal errors.
func (gu *GetUserHandler) GetUserBySessionID(w http.ResponseWriter, r *http.Request) {
	WriteHeaders(w, r)

	gu.deps.Log.Debugf("Request cookies size: %d", len(r.Cookies()))

	id := GetSessionFromCookie(r)
	if id == "" {
		err := errors.New("session ID does not exist in cookie")
		gu.deps.Log.Infof("Cookie %s does not exist for session ID", CookieSessionKey)

		WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	ua, err := gu.Gateway.User.GetUserBySessionID(id)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if ua == nil {
		WriteResponse(w, NewApiResponse(nil))
	} else {
		res := NewApiResponse(ua)
		i, err := WriteResponse(w, res)
		if err != nil {
			gu.deps.Log.Warnf("failed to write response: %v", err)
			WriteErrorResponse(w, err, http.StatusInternalServerError)

			return
		}

		gu.deps.Log.Infof("Successfully written %d byte(s) to response", i)
	}
}

type PostUserHandler struct {
	Gateway *dbcon.Gateway
	deps    *utils.Deps
}

// Login is the handler for handling the login and authentication.
// The username and password will be validated if the cookie is not found
// with a valid session ID.
//
// The middleware does not effect this, but the session ID is to prevent reauth.
//
// If the user is successfully authenticated, the output of the response
// will contain the status and the session ID will be written to the cookie.
func (pu *PostUserHandler) Login(w http.ResponseWriter, r *http.Request) {
	pu.deps.Log.Infof("Login handler accessed (%v)", r.RemoteAddr)
	var user RequestUserInfo

	WriteHeaders(w, r)
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
			_, err := WriteResponse(w, NewApiResponse(validSession))
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

	validUser, ua, err := pu.Gateway.User.ValidateUser(user.Username, user.Password)
	if err != nil {
		pu.deps.Log.Criticalf("Error occurred during user validation: %v", err)
		WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	pu.deps.Log.Debugf("%s result: %v", user.Username, validUser)

	res := NewApiResponse(validUser)
	if validUser {
		// a new login will always create a new session ID
		ses, err := pu.Gateway.Session.UpsertSession(ua.AccountID)
		if err != nil {
			WriteErrorResponse(w, err, http.StatusInternalServerError)
		}

		pu.deps.Log.Infof("Setting session data for user %s", ua.Username)
		SetCookieSession(w, ses)
		_, err = WriteResponse(w, res)
		if err != nil {
			WriteErrorResponse(w, err, http.StatusInternalServerError)
		}
	} else {
		err = fmt.Errorf("user %s does not exist", user.Username)
		WriteErrorResponse(w, err, http.StatusBadRequest)
	}
}

// RegisterUser is the handler for registering users. A new session.Session struct
// will be marshalled as the response.
func (pu *PostUserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user RequestUserInfo

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		pu.deps.Log.Warnf("Failed to decode JSON: %v", err)
		WriteErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	acc, err := pu.Gateway.User.AddUser(user.Username, user.Password)
	if err != nil {
		pu.deps.Log.Warnf("Failed to add user: %v", err)
		WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	s, err := pu.Gateway.Session.UpsertSession(acc.AccountID)
	if err != nil {
		pu.deps.Log.Warnf("Failed to upsert session: %v", err)
		WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	a := NewApiResponse(s)
	_, err = WriteResponse(w, a)
	if err != nil {
		pu.deps.Log.Warnf("Failed to write response: %v", err)
		WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}
