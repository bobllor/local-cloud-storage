package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/bobllor/assert"
	dbcon "github.com/bobllor/cloud-project/src/db_gateway"
	"github.com/bobllor/cloud-project/src/server"
	"github.com/bobllor/cloud-project/src/tests"
	"github.com/bobllor/cloud-project/src/utils"
)

func TestPostRegisterUser(t *testing.T) {
	sv, err := getTestServer()
	assert.Nil(t, err)

	gw, err := getTestGateway()
	assert.Nil(t, err)

	uh := NewUserHandler(gw)

	sv.RegisterHandler(UserRegisterRoute, uh.Post.RegisterUser)

	mSv := httptest.NewServer(sv.Handler)
	defer mSv.Close()

	url := mSv.URL

	c := mSv.Client()

	b, err := json.Marshal(map[string]string{"username": "hello", "password": "no"})
	assert.Nil(t, err)

	res, err := c.Post(url+"/register", ContentJson, bytes.NewBuffer(b))
	assert.Nil(t, err)
	assert.True(t, res.StatusCode < 300 && res.StatusCode >= 200)

	defer res.Body.Close()
}

func getTestServer() (*server.Server, error) {
	addr := ":8080"

	serv, err := server.NewServer(addr)
	if err != nil {
		return nil, err
	}

	return serv, nil
}

func getTestGateway() (*dbcon.Gateway, error) {
	dbcfg := dbcon.NewConfig(
		tests.DbMetaInfo.User,
		tests.DbMetaInfo.Password,
		tests.DbMetaInfo.Net,
		tests.DbMetaInfo.Addr,
		tests.DbMetaInfo.DbName,
	)

	tdb, err := dbcon.NewDatabase(dbcfg)
	if err != nil {
		return nil, err
	}

	deps := utils.NewTestDeps()

	fg := dbcon.NewFileGateway(tdb, deps)
	ug := dbcon.NewUserGateway(tdb, deps)
	sg := dbcon.NewSessionGateway(tdb, deps)

	gw := dbcon.NewGateway(fg, ug, sg)

	return gw, nil
}
