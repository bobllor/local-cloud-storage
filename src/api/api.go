package api

import (
	"net/http"

	dbcon "github.com/bobllor/cloud-project/src/db_gateway"
	"github.com/bobllor/gologger"
)

const (
	ContentJson = "application/json"
)

// NewApi creates a new Api struct.
func NewApi(gw *dbcon.Gateway, logger *gologger.Logger) *Api {
	api := &Api{
		User: NewUserHandler(gw, logger),
	}

	return api
}

type Api struct {
	User     *UserHandler
	Handlers HandlerMap
}

// GetHandlers
func (a *Api) GetHandlers() HandlerMap {
	return a.Handlers
}

// addUserHandlers is used to add the HandlerFunc with a route, used for the HandlerMap.
// This function will panic if a duplicate handler is added.
func (a *Api) addUserHandlers() {
	u := a.User

	a.addHandler(UserPostRegisterRoute, u.Post.RegisterUser)
}

// addHandler adds a handler to a.HandlerMap. If a duplicate route is found, then this method
// will panic.
func (a *Api) addHandler(route string, handleFunc func(w http.ResponseWriter, r *http.Request)) {
	_, ok := a.Handlers[route]
	if !ok {
		a.Handlers[route] = handleFunc
	} else {
		panic("duplicate route found")
	}
}
