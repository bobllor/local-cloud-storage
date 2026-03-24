package api

import (
	dbcon "github.com/bobllor/cloud-project/src/db_con"
)

// DBAPI is used for handling API requests.
type DBAPI struct {
	FilesDB dbcon.FilesDB
}

func NewDBAPI(filesDB dbcon.FilesDB) *DBAPI {
	api := &DBAPI{
		FilesDB: filesDB,
	}

	return api
}
