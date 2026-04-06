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

const (
	testUserName  = "test.username"
	testPassword  = "anothertestpassword"
	testPhcString = "$argon2id$v=19$m=65536,t=2,p=4$QTdpUkJ3c3J0amlOT2huV2VBR2duZw$vzICl8p5CVfpGfypDV4yIVULsYatAmir6B8nHWtcPtE"
)

func TestGetUserID(t *testing.T) {
	udb, err := newTestUserGateway()
	assert.Nil(t, err)

	user, err := udb.GetUserByID(testUserAccountID)
	assert.Nil(t, err)

	assert.Equal(t, user.AccountID, testUserAccountID)
	assert.True(t, user.Active)
	assert.Equal(t, user.Username, "test.username")

	_, err = hasher.ParsePHC(user.PasswordHash)
	assert.Nil(t, err)
}

func TestGetUserByUsername(t *testing.T) {
	udb, err := newTestUserGateway()
	assert.Nil(t, err)

	user, err := udb.GetUserByUsername(testUserName)
	assert.Nil(t, err)

	assert.Equal(t, user.Username, testUserName)
	assert.Equal(t, user.PasswordHash, testPhcString)
}

func TestCheckCredentials(t *testing.T) {
	udb, err := newTestUserGateway()
	assert.Nil(t, err)

	status, err := udb.CheckCredentials(testUserName, testPassword)
	assert.Nil(t, err)

	assert.True(t, status)
}

func TestCheckCredentialsInvalid(t *testing.T) {
	udb, err := newTestUserGateway()
	assert.Nil(t, err)

	status, err := udb.CheckCredentials("userdoesnotexist", testPassword)
	assert.NotNil(t, err)
	assert.False(t, status)

	status, err = udb.CheckCredentials(testUserName, "invalidpassword")
	assert.Nil(t, err)
	assert.False(t, status)
}

func TestAddUser(t *testing.T) {
	udb, err := newTestUserGateway()
	assert.Nil(t, err)
	username := "a user here"
	password := "somepasswordhere"

	acc, err := udb.AddUser(username, password)
	assert.Nil(t, err)

	_, err = DropRows(udb.database, user.TableName, user.ColumnAccountID, acc.AccountID)
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
	_, err = DropRows(udb.database, user.TableName, user.ColumnAccountID, acc.AccountID)
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
