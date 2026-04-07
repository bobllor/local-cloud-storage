package api

import (
	dbcon "github.com/bobllor/cloud-project/src/db_gateway"
)

type Api struct {
	User *UserHandler
}

// NewApi creates a new Api struct.
func NewApi(gw *dbcon.Gateway) *Api {
	api := &Api{
		User: NewUserHandler(gw),
	}

	return api
}
