package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
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
	username := "john.doe"

	uh := NewUserHandler(gw, tests.NewTestLogger())
	sv.RegisterHandlerFunc(UserPostRegisterRoute, uh.PostRegisterUser)

	mSv := httptest.NewServer(sv.Handler)
	defer mSv.Close()
	url := mSv.URL + "/api/register"
	c := mSv.Client()

	t.Run("Normal registration", func(t *testing.T) {
		t.Cleanup(func() {
			_, err := dbcon.DropRows(db, user.TableName, user.ColumnUsername, username)
			assert.Nil(t, err)
		})

		b, err := json.Marshal(map[string]string{"username": username, "password": "apasswordhere"})
		assert.Nil(t, err)

		res, err := c.Post(url, ContentJson, bytes.NewBuffer(b))
		assert.Nil(t, err)
		assert.True(t, res.StatusCode <= http.StatusBadRequest)
		defer res.Body.Close()

		var apres ApiResponse
		err = json.NewDecoder(res.Body).Decode(&apres)
		assert.Nil(t, err)
		assert.NotNil(t, apres)

		assert.Equal(t, len(res.Cookies()), 1)
		assert.NotNil(t, apres.Output)
	})

	t.Run("Duplicate registration", func(t *testing.T) {
		b, err := json.Marshal(map[string]string{"username": tests.DbRowInfo.Username, "password": "whatever"})
		assert.Nil(t, err)

		res, err := c.Post(url, ContentJson, bytes.NewBuffer(b))
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		defer res.Body.Close()

		var apres ApiResponse
		err = json.NewDecoder(res.Body).Decode(&apres)
		assert.Nil(t, err)

		assert.Equal(t, apres.Status, StatusError)
		assert.Equal(t, apres.Error.Code, http.StatusBadRequest)
		assert.Equal(t, apres.Error.Reason, ReasonUserAlreadyExists)
	})
}

func TestLoginUser(t *testing.T) {
	sv := getTestServer(t)
	gw, db := getGatewayDb(t)

	uh := NewUserHandler(gw, tests.NewTestLogger())
	sv.RegisterHandlerFunc(UserPostLoginRoute, uh.PostLogin)

	tsv := httptest.NewServer(sv.Handler)
	defer tsv.Close()
	tc := tsv.Client()

	url := tsv.URL + "/api/login"

	t.Run("User Exists", func(t *testing.T) {
		b, err := json.Marshal(map[string]string{
			"username": tests.DbRowInfo.Username,
			"password": tests.TestPassword,
		})
		assert.Nil(t, err)

		t.Cleanup(func() {
			_, err := dbcon.UpdateRow(
				db,
				session.TableName,
				session.ColumnAccountID,
				tests.DbRowInfo.AccountID,
				dbcon.ClauseData{
					Columns: []string{session.ColumnSessionID},
					Args:    []any{tests.DbRowInfo.SessionID},
				},
			)
			assert.Nil(t, err)
		})

		res, err := tc.Post(url, ContentJson, bytes.NewBuffer(b))
		assert.Nil(t, err)
		defer res.Body.Close()

		var v ApiResponse
		err = json.NewDecoder(res.Body).Decode(&v)
		assert.Nil(t, err)

		assert.Equal(t, v.Status, StatusSuccess)
		assert.Equal(t, v.Output, true)
	})

	t.Run("Login Fail", func(t *testing.T) {
		b, err := json.Marshal(map[string]string{
			"username": "nonexistent.username",
			"password": tests.TestPassword,
		})
		assert.Nil(t, err)

		res, err := tc.Post(url, ContentJson, bytes.NewBuffer(b))
		assert.Nil(t, err)
		defer res.Body.Close()

		var v ApiResponse
		err = json.NewDecoder(res.Body).Decode(&v)
		assert.Nil(t, err)

		assert.Equal(t, v.Status, StatusError)
		assert.Equal(t, v.Error.Code, http.StatusBadRequest)
	})
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
