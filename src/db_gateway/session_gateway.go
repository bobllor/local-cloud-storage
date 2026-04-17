package dbgateway

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	querybuilder "github.com/bobllor/cloud-project/src/query_builder"
	"github.com/bobllor/cloud-project/src/session"
	"github.com/bobllor/cloud-project/src/utils"
	"github.com/google/uuid"
)

const (
	ExpireTimeDays = 14
)

func NewSessionGateway(db *sql.DB, deps *utils.Deps) *SessionGateway {
	sg := &SessionGateway{
		database:          db,
		sessionFieldCount: session.ColumnSize,
		deps:              deps,
	}

	return sg
}

type SessionGateway struct {
	database          *sql.DB
	sessionFieldCount int
	deps              *utils.Deps
}

// GetSessionByAccountID retrieves the session of the account ID. If a session is not found,
// then it will return nil.
func (sg *SessionGateway) GetSessionByAccountID(accountID string) (*session.Session, error) {
	accSession := []session.Session{}

	sb := querybuilder.NewSqlBuilder(session.TableName)
	query := sb.Select().AllColumns().Where().Equal(session.ColumnAccountID, accountID).Build()

	sg.deps.Log.Debugf("Query: %s", query)
	rows, err := sg.database.Query(query, sb.Args()...)
	if err != nil {
		return nil, err
	}

	err = SelectRows(rows, &accSession)
	if err != nil {
		return nil, err
	}

	if len(accSession) == 0 {
		return nil, nil
	}

	return &accSession[0], nil
}

// GetSessionBySessionID retrieves the session by session ID. If a session is not found,
// then it will return nil.
func (sg *SessionGateway) GetSessionBySessionID(sessionID string) (*session.Session, error) {
	if !sg.validateID(sessionID) {
		return nil, nil
	}

	sb := querybuilder.NewSqlBuilder(session.TableName)
	query := sb.Select().AllColumns().Where().Equal(session.ColumnSessionID, sessionID).Build()

	sg.deps.Log.Debugf("Query: %s", query)
	rows, err := sg.database.Query(query, sb.Args()...)
	if err != nil {
		return nil, err
	}

	var ses []session.Session
	err = SelectRows(rows, &ses)
	if err != nil {
		return nil, err
	}

	if len(ses) == 0 {
		sg.deps.Log.Info("No session ID found")
		return nil, nil
	}

	return &ses[0], nil
}

// UpsertSession adds a new or updates an existing entry for the Session table,
// generating a new session ID associated with the account ID to maintain a session.
// It will return the Session that was added if successful.
//
// Existing account IDs will have get updated with new session ID
// dates, and if it is expired then a new session ID is generated.
func (sg *SessionGateway) UpsertSession(accountID string) (*session.Session, error) {
	sessionID := uuid.NewString()
	query := fmt.Sprintf("INSERT INTO %s", session.TableName)

	currTime := time.Now().UTC()
	expireTime := currTime.AddDate(0, 0, ExpireTimeDays)
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
		session.ColumnSessionID,
		session.ColumnCreatedOn,
		session.ColumnExpireOn,
	)

	AppendArgs(&args, sessionID, currTime, expireTime)

	query = query + " " + "VALUES" + placeholder + " " + duplicateStr

	_, err := execQuery(sg.database, query, args...)
	if err != nil {
		sg.deps.Log.Warnf("Failed to execute query: %s | Values: %v", query, args)
		return nil, err
	}

	return &ses, nil
}

// ValidateSession validates a session with the user's session ID.
// If the sessionID is invalid, it does not exist, or it does not match the stored database
// version, then it will return false.
//
// Any errors will be returned during the database query.
func (sg *SessionGateway) ValidateSession(sessionID string) (bool, error) {
	// false conditions:
	//	- any DB errors (error must be handled)
	//	- sessionID are empty strings or invalid formatting
	//	- sessionID does not match stored sessionID
	//	- stored expiration date is < current time
	//	- session row is not found with account ID

	if !sg.validateID(sessionID) {
		return false, nil
	}

	// TODO: add cache access here, probably redis or if you are lazy a hash map

	ses, err := sg.GetSessionBySessionID(sessionID)
	if err != nil {
		return false, err
	}
	if ses == nil {
		return false, nil
	}

	if ses.SessionID != sessionID {
		return false, nil
	}
	if ses.ExpireOn.UTC().Before(time.Now().UTC()) {
		return false, nil
	}

	return true, nil
}

// validateID validates the ID string formatting. It returns a true
// if it is valid, otherwise it will return false.
// This does not check the database.
func (sg *SessionGateway) validateID(id string) bool {
	if strings.TrimSpace(id) == "" {
		return false
	}

	return true
}
