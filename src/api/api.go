package api

import (
	dbcon "github.com/bobllor/cloud-project/src/db_gateway"
)

// DBAPI is used for handling API requests.
type DBAPI struct {
	FileGateway dbcon.FileGateway
}

func NewDBAPI(fileGateWay dbcon.FileGateway) *DBAPI {
	api := &DBAPI{
		FileGateway: fileGateWay,
	}

	return api
}
