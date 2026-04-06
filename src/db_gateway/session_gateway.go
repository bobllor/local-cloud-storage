package dbgateway

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bobllor/cloud-project/src/config"
	"github.com/bobllor/cloud-project/src/session"
	"github.com/google/uuid"
)

func NewSessionGateway(db *sql.DB, cfg *config.Config) (*SessionGateway, error) {
	sg := &SessionGateway{
		database:          db,
		sessionFieldCount: session.ColumnSize,
		cfg:               cfg,
		util:              DBUtility{log: cfg.Log},
	}

	return sg, nil
}

type SessionGateway struct {
	database          *sql.DB
	sessionFieldCount int
	cfg               *config.Config
	util              DBUtility
}

// GetSession retrieves the session of the account ID.
func (sg *SessionGateway) GetSession(accountID string) (*session.Session, error) {
	cb := NewClauseBuilder()

	cb.Equal(session.ColumnAccountID, accountID)

	cbQ, args, err := cb.Build()
	if err != nil {
		return nil, err
	}

	baseQ := fmt.Sprintf("SELECT * FROM %s", session.TableName)

	accSession := &session.Session{}

	query := baseQ + " " + cbQ

	rows, err := sg.database.Query(query, args...)
	if err != nil {
		return nil, err
	}

	err = SelectRow(rows, accSession)
	if err != nil {
		return nil, err
	}

	// TODO: handle an edge case where no users appear,
	// this applies to other functions similar to this as well

	return accSession, nil
}

// UpsertSession adds a new or updates an existing entry for the Session table, ]
// generating a new session ID associated with the account ID to maintain a session.
// It will return the Session that was added if successful.
//
// Existing account IDs will have get updated with new session ID
// dates, and if it is expired then a new session ID is generated.
func (sg *SessionGateway) UpsertSession(accountID string) (*session.Session, error) {
	sessionID := uuid.NewString()
	query := fmt.Sprintf("INSERT INTO %s", session.TableName)

	// expiration date is 30 days from the current time
	currTime := time.Now().UTC()
	expireTime := currTime.AddDate(0, 0, 30)
	ses := session.Session{
		SessionID: sessionID,
		AccountID: accountID,
		CreatedOn: currTime,
		ExpireOn:  expireTime,
	}

	args := ses.ToArgs()

	placeholder := BuildPlaceholder(len(args), 1)

	duplicateStr := fmt.Sprintf(
		"ON DUPLICATE KEY UPDATE %s=?,%s=?,%s=?",
		session.ColumnColumnSessionID,
		session.ColumnCreatedOn,
		session.ColumnExpireOn,
	)

	args = append(args, sessionID)
	args = append(args, currTime)
	args = append(args, expireTime)

	query = query + " " + "VALUES" + placeholder + " " + duplicateStr

	_, err := execQuery(sg.database, query, args...)
	if err != nil {
		return nil, err
	}

	return &ses, nil
}
