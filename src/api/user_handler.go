package api

import (
	"encoding/json"
	"net/http"

	dbcon "github.com/bobllor/cloud-project/src/db_gateway"
)

const (
	UserRegisterRoute = "POST /register"
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
	Gateway *dbcon.Gateway
}

func (gu *GetUserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {

}

type PostUserHandler struct {
	Gateway *dbcon.Gateway
}

func (pu *PostUserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	w.WriteHeader(http.StatusAccepted)
}
