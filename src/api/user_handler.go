package api

import (
	"fmt"
	"net/http"

	dbcon "github.com/bobllor/cloud-project/src/db_gateway"
)

func NewUserHandler(gw *dbcon.Gateway) *UserHandler {
	uh := &UserHandler{
		Post: PostUserHandler{
			Gateway: gw,
		},
		Get: GetUserHandler{
			Gateway: gw,
		},
	}

	return uh
}

// UserHandler contains handlers used for handling user related
// logic.
type UserHandler struct {
	Post PostUserHandler
	Get  GetUserHandler
}

type GetUserHandler struct {
	Gateway  *dbcon.Gateway
	Handlers HandlerMap
}

func (gu *GetUserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {

}

type PostUserHandler struct {
	Gateway  *dbcon.Gateway
	Handlers HandlerMap
}

func (pu *PostUserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
}
