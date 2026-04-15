package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/bobllor/assert"
	dbcon "github.com/bobllor/cloud-project/src/db_gateway"
	"github.com/bobllor/cloud-project/src/server"
	"github.com/bobllor/cloud-project/src/session"
	"github.com/bobllor/cloud-project/src/tests"
	"github.com/bobllor/cloud-project/src/user"
	"github.com/bobllor/cloud-project/src/utils"
)

func TestPostRegisterUser(t *testing.T) {
	sv := getTestServer(t)
	gw, db := getGatewayDb(t)

	uh := NewUserHandler(gw, tests.NewTestLogger())
	sv.RegisterHandler(UserPostRegisterRoute, uh.Post.RegisterUser)

	mSv := httptest.NewServer(sv.Handler)
	defer mSv.Close()
	url := mSv.URL
	c := mSv.Client()

	b, err := json.Marshal(map[string]string{"username": "john.doe", "password": "apasswordhere"})
	assert.Nil(t, err)

	res, err := c.Post(url+"/register", ContentJson, bytes.NewBuffer(b))
	assert.Nil(t, err)
	assert.True(t, res.StatusCode < 300 && res.StatusCode >= 200)

	var ses session.Session
	err = json.NewDecoder(res.Body).Decode(&ses)
	assert.Nil(t, err)
	assert.NotNil(t, ses)

	_, err = dbcon.DropRows(db, user.TableName, user.ColumnAccountID, ses.AccountID)
	assert.Nil(t, err)

	defer res.Body.Close()
}

// getTestServer creates a new Server test instance.
func getTestServer(t *testing.T) *server.Server {
	addr := ":8080"

	serv, err := server.NewServer(addr)
	assert.Nil(t, err)

	return serv
}

// getGatewayDb creates a test dbcon.Gateway and a sql.DB for use.
// If an error occurs, then it will fatal and exit.
func getGatewayDb(t *testing.T) (*dbcon.Gateway, *sql.DB) {
	dbcfg := dbcon.NewConfig(
		tests.DbMetaInfo.User,
		tests.DbMetaInfo.Password,
		tests.DbMetaInfo.Net,
		tests.DbMetaInfo.Addr,
		tests.DbMetaInfo.DbName,
	)

	tdb, err := dbcon.NewDatabase(dbcfg)
	assert.Nil(t, err)

	deps := utils.NewTestDeps()

	fg := dbcon.NewFileGateway(tdb, deps)
	ug := dbcon.NewUserGateway(tdb, deps)
	sg := dbcon.NewSessionGateway(tdb, deps)

	gw := dbcon.NewGateway(fg, ug, sg)

	return gw, tdb
}
