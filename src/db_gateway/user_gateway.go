package dbgateway

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bobllor/cloud-project/src/config"
	"github.com/bobllor/cloud-project/src/hasher"
	"github.com/bobllor/cloud-project/src/session"
	"github.com/bobllor/cloud-project/src/user"
	"github.com/google/uuid"
)

type UserGateway struct {
	database       *sql.DB
	userFieldCount int
	config         *config.Config
	util           DBUtility
}

func NewUserGateway(db *sql.DB, config *config.Config) *UserGateway {
	return &UserGateway{
		database:       db,
		userFieldCount: user.UserColumnSize,
		config:         config,
		util: DBUtility{
			log: config.Log,
		},
	}
}

// AddUser adds a new user into the database. It will return the UserAccount
// that was created in the database, or an error if one occurred.
//
// The password is stored as the PHC string from the password hashing function.
func (ug *UserGateway) AddUser(username string, password string) (*user.UserAccount, error) {
	baseQuery := fmt.Sprintf("INSERT INTO %s VALUES", user.UserTableName)

	accountID := uuid.NewString()
	raw, err := hasher.Hash(password, nil, hasher.DefaultArgon2Params)
	if err != nil {
		return nil, err
	}

	hashRes := raw.Encode()

	acc := &user.UserAccount{
		AccountID:    accountID,
		PasswordHash: hashRes.PHC,
		Username:     username,
		CreatedOn:    time.Now().UTC(),
		Active:       true,
	}

	args := acc.ToArgs()
	placeholders := BuildPlaceholder(len(args), 1)

	query := baseQuery + " " + placeholders

	ug.util.LogQueryAndArgs(query, args)

	res, err := execQuery(ug.database, query, args...)
	if err != nil {
		return nil, err
	}

	ug.util.LogResultRows(res)

	return acc, err
}

// GetUser retrieves the full row of the user.
func (ug *UserGateway) GetUser(accountID string) (*user.UserAccount, error) {
	cb := NewClauseBuilder()
	cb.Equal(user.ColumnAccountID, accountID)

	baseQuery := fmt.Sprintf("SELECT * FROM %s", user.UserTableName)

	cbQ, args, err := cb.Build()
	if err != nil {
		return nil, err
	}

	user := []user.UserAccount{}

	query := baseQuery + " " + cbQ
	rows, err := ug.database.Query(query, args...)
	if err != nil {
		return nil, err
	}

	err = SelectRows(rows, &user)
	if err != nil {
		return nil, err
	}

	if len(user) == 0 {
		return nil, fmt.Errorf("failed to retrieve row with the ID %s", accountID)
	}

	return &user[0], nil
}

// NewSession creates a new entry for the Session table, generating
// a new session ID associated with the account ID to maintain a session.
func (ug *UserGateway) NewSession(accountID string) error {
	sessionID := uuid.NewString()
	query := fmt.Sprintf("REPLACE INTO %s", session.TableName)

	session := session.Session{
		SessionID: sessionID,
		AccountID: accountID,
		CreatedOn: time.Now().UTC(),
	}

	args := session.ToArgs()

	placeholder := BuildPlaceholder(len(args), 1)

	query = query + " " + "VALUES" + placeholder

	_, err := execQuery(ug.database, query, args...)
	if err != nil {
		return err
	}

	return nil
}
