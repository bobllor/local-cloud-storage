package session

import "time"

const (
	ColumnSize            = 4
	TableName             = "Session"
	ColumnColumnSessionID = "SessionID"
	ColumnAccountID       = "AccountID"
	ColumnCreatedOn       = "CreatedOn"
	ColumnExpireOn        = "ExpireOn"
)

type Session struct {
	SessionID string
	AccountID string
	CreatedOn time.Time
	ExpireOn  time.Time
}

func (s *Session) ToArgs() []any {
	args := []any{}

	args = append(args, s.SessionID)
	args = append(args, s.AccountID)
	args = append(args, s.CreatedOn)
	args = append(args, s.ExpireOn)

	return args
}
