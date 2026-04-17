package user

import (
	"time"
)

const (
	ColumnSize         = 5
	TableName          = "UserAccount"
	ColumnAccountID    = "AccountID"
	ColumnUsername     = "Username"
	ColumnPasswordHash = "PasswordHash"
	ColumnCreatedOn    = "CreatedOn"
	ColumnActive       = "Active"
)

type UserAccount struct {
	AccountID    string
	Username     string
	PasswordHash string
	CreatedOn    time.Time
	Active       bool
}

// UserAccountNoPassword is the UserAccount struct
// but does not include the password hash.
type UserAccountNoPassword struct {
	AccountID string    `json:"account_id"`
	Username  string    `json:"username"`
	CreatedOn time.Time `json:"created_on"`
	Active    bool      `json:"active"`
}

// ToArgs converts the struct into an any slice.
// This is used for query arguments.
func (ua *UserAccount) ToArgs() []any {
	args := []any{}

	args = append(args, ua.AccountID)
	args = append(args, ua.Username)
	args = append(args, ua.PasswordHash)
	// converted to date string for adding as a date
	args = append(args, ua.CreatedOn.UTC())
	args = append(args, ua.Active)

	return args
}
