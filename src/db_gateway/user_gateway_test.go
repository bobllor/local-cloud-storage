package dbgateway

import (
	"testing"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/hasher"
	"github.com/bobllor/cloud-project/src/tests"
	"github.com/bobllor/cloud-project/src/user"
	"github.com/bobllor/cloud-project/src/utils"
)

const (
	testPassword = "anothertestpassword"
)

func TestGetUserID(t *testing.T) {
	udb, err := newTestUserGateway()
	assert.Nil(t, err)

	user, err := udb.GetUserByID(tests.DbRowInfo.AccountID)
	assert.Nil(t, err)

	assert.Equal(t, user.AccountID, tests.DbRowInfo.AccountID)
	assert.True(t, user.Active)
	assert.Equal(t, user.Username, "test.username")

	_, err = hasher.ParsePHC(user.PasswordHash)
	assert.Nil(t, err)
}

func TestGetUserByUsername(t *testing.T) {
	udb, err := newTestUserGateway()
	assert.Nil(t, err)

	user, err := udb.GetUserByUsername(tests.DbRowInfo.Username)
	assert.Nil(t, err)

	assert.Equal(t, user.Username, tests.DbRowInfo.Username)
	assert.Equal(t, user.PasswordHash, tests.DbRowInfo.PhcString)
}

func TestCheckCredentials(t *testing.T) {
	udb, err := newTestUserGateway()
	assert.Nil(t, err)

	status, err := udb.CheckCredentials(tests.DbRowInfo.Username, testPassword)
	assert.Nil(t, err)

	assert.True(t, status)
}

func TestCheckCredentialsInvalid(t *testing.T) {
	udb, err := newTestUserGateway()
	assert.Nil(t, err)

	status, err := udb.CheckCredentials("userdoesnotexist", testPassword)
	assert.NotNil(t, err)
	assert.False(t, status)

	status, err = udb.CheckCredentials(tests.DbRowInfo.Username, "invalidpassword")
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

	deps := utils.NewTestDeps()

	ug := NewUserGateway(db, deps)

	return ug, nil
}
