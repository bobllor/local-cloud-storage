package dbgateway

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bobllor/cloud-project/src/hasher"
	"github.com/bobllor/cloud-project/src/session"
	"github.com/bobllor/cloud-project/src/sqlquery"
	"github.com/bobllor/cloud-project/src/user"
	"github.com/bobllor/cloud-project/src/utils"
	"github.com/google/uuid"
)

type UserGateway struct {
	database       *sql.DB
	userFieldCount int
	deps           *utils.Deps
}

func NewUserGateway(db *sql.DB, deps *utils.Deps) *UserGateway {
	return &UserGateway{
		database:       db,
		userFieldCount: user.ColumnSize,
		deps:           deps,
	}
}

// AddUser adds a new user into the database. It will return the UserAccount
// that was created in the database, or an error if one occurred.
//
// The password is stored as the PHC string from the password hashing function.
func (ug *UserGateway) AddUser(username string, password string) (*user.UserAccount, error) {
	accountID := uuid.NewString()
	raw, err := hasher.Hash(password, nil, hasher.DefaultArgon2Params)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
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

	query, args, err := sqlquery.InsertInto(
		user.TableName,
		user.ColumnAccountID,
		user.ColumnUsername,
		user.ColumnPasswordHash,
		user.ColumnCreatedOn,
		user.ColumnActive,
	).Args(args...).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build INSERT INTO query: %v", err)
	}

	ug.deps.Log.Debugf("Query: %s | Args: %d", query, len(args))
	res, err := execQuery(ug.database, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	logResultRows(ug.deps.Log, res)

	return acc, nil
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
	// TODO: temp hold while testing
	query, args, err := sqlquery.Select(user.TableName).Where().Equal(user.ColumnUsername, username).Build()
	if err != nil {
		return nil, err
	}

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
	query, args, err := sqlquery.Select(user.TableName).Where().Equal(user.ColumnAccountID, accountID).Build()
	if err != nil {
		return nil, err
	}

	user := user.UserAccount{}

	ug.deps.Log.Debugf("Query: %s", query)

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

	sQuery, sArgs, err := sqlquery.Select(session.TableName, session.ColumnAccountID).
		Where().Equal(session.ColumnSessionID, sessionID).Build()
	if err != nil {
		return nil, err
	}

	query, args, err := sqlquery.Select(
		user.TableName,
		user.ColumnAccountID,
		user.ColumnUsername,
		user.ColumnCreatedOn,
		user.ColumnActive,
	).Where().Exists(sQuery, sArgs...).Build()
	if err != nil {
		return nil, err
	}

	ug.deps.Log.Debugf("Query: %s", query)
	rows, err := ug.database.Query(query, args...)
	if err != nil {
		ug.deps.Log.Criticalf("Failed to query data in SQL: %v", err)
		return nil, err
	}

	var us []user.UserAccountNoPassword
	err = SelectRows(rows, &us)
	if err != nil {
		ug.deps.Log.Criticalf("Failed to parse SQL rows: %v", err)
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

	logResultRows(ug.deps.Log, res)

	return nil
}

// RestoreUserByID removes the soft deletion state from the account ID.
func (ug *UserGateway) RestoreUserByID(accountID string) error {
	return nil
}
