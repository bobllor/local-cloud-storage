package dbgateway

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bobllor/cloud-project/src/config"
	"github.com/bobllor/cloud-project/src/hasher"
	"github.com/bobllor/cloud-project/src/user"
	"github.com/google/uuid"
)

type UserGateway struct {
	database       *sql.DB
	userFieldCount int
	config         *config.Config
	util           LogUtility
}

func NewUserGateway(db *sql.DB, config *config.Config) *UserGateway {
	return &UserGateway{
		database:       db,
		userFieldCount: user.UserColumnSize,
		config:         config,
		util: LogUtility{
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
