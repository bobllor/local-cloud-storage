package dbgateway

import (
	"io"
	"log"
	"testing"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/config"
	"github.com/bobllor/cloud-project/src/hasher"
	"github.com/bobllor/cloud-project/src/user"
	"github.com/bobllor/gologger"
)

func TestGetUser(t *testing.T) {
	udb, err := newTestUserGateway()
	assert.Nil(t, err)

	user, err := udb.GetUser(testUserAccountID)
	assert.Nil(t, err)

	assert.Equal(t, user.AccountID, testUserAccountID)
	assert.Equal(t, user.Active, true)
	assert.Equal(t, user.Username, "test.username")

	_, err = hasher.ParsePHC(user.PasswordHash)
	assert.Nil(t, err)
}

func TestAddUser(t *testing.T) {
	udb, err := newTestUserGateway()
	assert.Nil(t, err)
	username := "a user here"
	password := "somepasswordhere"

	acc, err := udb.AddUser(username, password)
	assert.Nil(t, err)

	_, err = devDropRows(udb.database, user.TableName, user.ColumnAccountID, acc.AccountID)
	assert.Nil(t, err)
}

func TestAddUserComparePassword(t *testing.T) {
	udb, err := newTestUserGateway()
	assert.Nil(t, err)
	username := "a user here"
	password := "somepasswordhere"

	acc, err := udb.AddUser(username, password)
	assert.Nil(t, err)

	// drop row immediately in case of failures below, the rest doesnt need the table data
	_, err = devDropRows(udb.database, user.TableName, user.ColumnAccountID, acc.AccountID)
	assert.Nil(t, err)

	baseRes, err := hasher.ParsePHC(acc.PasswordHash)
	assert.Nil(t, err)

	baseSalt, err := baseRes.DecodeSalt()
	assert.Nil(t, err)

	raw, err := hasher.Hash(password, baseSalt, baseRes.Params)
	assert.Nil(t, err)

	compareRes := raw.Encode()

	assert.Equal(t, compareRes.Hash, baseRes.Hash)
}

func newTestUserGateway() (*UserGateway, error) {
	cnf := newTestDBConfig()
	db, err := NewDatabase(cnf)
	if err != nil {
		return nil, err
	}

	logger := gologger.NewLogger(log.New(io.Discard, "", log.Ldate|log.Ltime), gologger.Lsilent)
	stdConfig := config.NewConfig(logger)

	ug := NewUserGateway(db, stdConfig)

	return ug, nil
}
