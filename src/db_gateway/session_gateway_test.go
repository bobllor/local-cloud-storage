package dbgateway

import (
	"io"
	"log"
	"testing"
	"time"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/config"
	"github.com/bobllor/cloud-project/src/session"
	"github.com/bobllor/cloud-project/src/user"
	"github.com/bobllor/gologger"
)

const (
	testSessionID = "7ca90f85-b1e0-4214-8ff6-4e3720cc8078"
)

func TestGetSession(t *testing.T) {
	sg, err := newTestSessionGateway()
	assert.Nil(t, err)

	s, err := sg.GetSession(testUserAccountID)
	assert.Nil(t, err)

	assert.Equal(t, s.AccountID, testUserAccountID)
	assert.Equal(t, s.SessionID, testSessionID)
}

func TestUpsertSessionNew(t *testing.T) {
	ug, err := newTestUserGateway()
	assert.Nil(t, err)
	sg, err := newTestSessionGateway()
	assert.Nil(t, err)

	usr, err := ug.AddUser("username.here", "password")
	assert.Nil(t, err)

	addS, err := sg.UpsertSession(usr.AccountID)
	assert.Nil(t, err)

	sesh, err := sg.GetSession(addS.AccountID)
	assert.Nil(t, err)

	_, err = DropRows(ug.database, user.TableName, user.ColumnAccountID, usr.AccountID)
	assert.Nil(t, err)

	assert.Equal(t, sesh.AccountID, addS.AccountID)
	assert.Equal(t, sesh.SessionID, addS.SessionID)
	assert.Equal(t, sesh.CreatedOn.Truncate(time.Minute), addS.CreatedOn.Truncate(time.Minute))
	assert.Equal(t, sesh.ExpireOn.Truncate(time.Minute), addS.ExpireOn.Truncate(time.Minute))
}

func TestUpsertSessionReplace(t *testing.T) {
	sg, err := newTestSessionGateway()
	assert.Nil(t, err)

	baseS, err := sg.GetSession(testUserAccountID)
	assert.Nil(t, err)

	sesh, err := sg.UpsertSession(testUserAccountID)
	assert.Nil(t, err)

	newS, err := sg.GetSession(testUserAccountID)
	assert.Nil(t, err)

	// reset the updated values
	_, err = UpdateRow(
		sg.database,
		session.TableName,
		session.ColumnAccountID,
		testUserAccountID,
		ClauseData{
			Columns: []string{session.ColumnColumnSessionID, session.ColumnCreatedOn, session.ColumnExpireOn},
			Args:    []any{baseS.SessionID, baseS.CreatedOn, baseS.ExpireOn},
		},
	)
	assert.Nil(t, err)

	assert.NotEqual(t, baseS.SessionID, newS.SessionID)
	assert.Equal(t, baseS.CreatedOn.Compare(newS.CreatedOn), -1)
	assert.Equal(t, baseS.ExpireOn.Compare(newS.ExpireOn), -1)

	assert.Equal(t, newS.SessionID, sesh.SessionID)
	assert.True(t, newS.CreatedOn.Truncate(time.Minute).Equal(sesh.CreatedOn.Truncate(time.Minute)))
	assert.True(t, newS.ExpireOn.Truncate(time.Minute).Equal(sesh.ExpireOn.Truncate(time.Minute)))
}

func newTestSessionGateway() (*SessionGateway, error) {
	dbCfg := newTestDBConfig()
	db, err := NewDatabase(dbCfg)
	if err != nil {
		return nil, err
	}

	logger := gologger.NewLogger(log.New(io.Discard, "", log.Ltime), gologger.Lsilent)

	cfg := config.NewConfig(logger)

	sg, err := NewSessionGateway(db, cfg)
	if err != nil {
		return nil, err
	}

	return sg, nil
}
