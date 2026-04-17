package dbgateway

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bobllor/cloud-project/src/hasher"
	"github.com/bobllor/cloud-project/src/session"
	"github.com/bobllor/cloud-project/src/user"
	"github.com/bobllor/cloud-project/src/utils"
	"github.com/google/uuid"
)

type UserGateway struct {
	database       *sql.DB
	userFieldCount int
	deps           *utils.Deps
	util           DBUtility
}

func NewUserGateway(db *sql.DB, deps *utils.Deps) *UserGateway {
	return &UserGateway{
		database:       db,
		userFieldCount: user.ColumnSize,
		deps:           deps,
		util: DBUtility{
			log: deps.Log,
		},
	}
}

// AddUser adds a new user into the database. It will return the UserAccount
// that was created in the database, or an error if one occurred.
//
// The password is stored as the PHC string from the password hashing function.
func (ug *UserGateway) AddUser(username string, password string) (*user.UserAccount, error) {
	baseQuery := fmt.Sprintf("INSERT INTO %s VALUES", user.TableName)

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

// ValidateUser validates if the credentials are correct for the user. The username and
// password is compared and will return a boolean and the user info. If an error occurs,
// then an error will be returned instead.
//
// If validation is true, then the user will always be returned.
func (ug *UserGateway) ValidateUser(username string, password string) (bool, *user.UserAccount, error) {
	user, err := ug.GetUserByUsername(username)
	if err != nil {
		return false, nil, err
	}
	if user == nil {
		return false, nil, nil
	}

	storedHash, err := hasher.ParsePHC(user.PasswordHash)
	if err != nil {
		return false, nil, err
	}

	validCredentials, err := hasher.Compare(password, storedHash)
	if err != nil {
		return false, nil, err
	}

	return validCredentials, user, nil
}

// GetUserByUsername gets the user row based on the username.
// If the user does not exist, then it will return nil.
func (ug *UserGateway) GetUserByUsername(username string) (*user.UserAccount, error) {
	cb := NewClauseBuilder()
	cb.Equal(user.ColumnUsername, username)

	baseQ := fmt.Sprintf("SELECT * FROM %s", user.TableName)

	cbQ, args, err := cb.Build()
	if err != nil {
		return nil, err
	}

	query := baseQ + " " + cbQ

	rows, err := ug.database.Query(query, args...)
	if err != nil {
		return nil, err
	}

	users := []user.UserAccount{}

	err = SelectRows(rows, &users)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		ug.deps.Log.Infof("user %s has no entries", username)
		return nil, nil
	}

	return &users[0], nil
}

// GetUserByID retrieves the user row based on the account ID.
func (ug *UserGateway) GetUserByID(accountID string) (*user.UserAccount, error) {
	cb := NewClauseBuilder()
	cb.Equal(user.ColumnAccountID, accountID)

	baseQuery := fmt.Sprintf("SELECT * FROM %s", user.TableName)

	cbQ, args, err := cb.Build()
	if err != nil {
		return nil, err
	}

	user := user.UserAccount{}

	query := baseQuery + " " + cbQ
	rows, err := ug.database.Query(query, args...)
	if err != nil {
		return nil, err
	}

	err = SelectRow(rows, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserBySessionID retrieves the user row based on the session ID. If the
// session ID is not found, then it will return nil.
// The return will contain all the columns except the password column of the row.
//
// If there is a session ID, then there will be an existing user. This does not apply vice versa.
func (ug *UserGateway) GetUserBySessionID(sessionID string) (*user.UserAccountNoPassword, error) {
	var us []user.UserAccountNoPassword

	// TODO: add the new sql builder in the future.
	// this is hard coded for now, as the new sql builder is in progress for writing
	subquery := fmt.Sprintf("SELECT %s FROM %s WHERE %s=?", session.ColumnSessionID, session.TableName, session.ColumnSessionID)
	whereClause := fmt.Sprintf("WHERE EXISTS (%s)", subquery)
	query := fmt.Sprintf(
		"SELECT %s,%s,%s,%s FROM %s %s",
		user.ColumnAccountID,
		user.ColumnUsername,
		user.ColumnCreatedOn,
		user.ColumnActive,
		user.TableName,
		whereClause,
	)

	rows, err := ug.database.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	err = SelectRows(rows, &us)
	if err != nil {
		return nil, err
	}

	if len(us) == 0 {
		ug.deps.Log.Infof("No user found with provided session ID")
		return nil, nil
	}

	return &us[0], nil
}

// DeleteUserByID sets an account ID for deletion. This is a soft deletion,
// it sets the Active column to false for deletion at a later date.
func (ug *UserGateway) DeleteUserByID(accountID string) error {
	cb := NewClauseBuilder()
	cb.Equal(user.ColumnAccountID, accountID)

	cbq, args, err := cb.Build()
	if err != nil {
		return err
	}

	cd := NewClauseData()

	cd.AddColumns(user.ColumnActive)
	cd.AddArgs(false)

	sq, sargs, err := cd.BuildSetQuery()
	if err != nil {
		return err
	}

	baseQuery := fmt.Sprintf("UPDATE %s %s %s", user.TableName, sq, cbq)

	execArgs := MakeArgs(sargs, args)

	res, err := execQuery(ug.database, baseQuery, execArgs...)
	if err != nil {
		return err
	}

	ug.util.LogResultRows(res)

	return nil
}

// RestoreUserByID removes the soft deletion state from the account ID.
func (ug *UserGateway) RestoreUserByID(accountID string) error {
	return nil
}
