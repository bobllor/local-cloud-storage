package dbgateway

import (
	"database/sql"

	"github.com/bobllor/gologger"
)

// LogUtility is used to remove repeatable logging for SQL queries.
type LogUtility struct {
	log *gologger.Logger
}

// LogResultRows is used to log the SQL result rows. If an error is encountered during
// the result parsing, then it will not log the affected rows.
// This will log to the Info level with no error, and the Warn level if an error does occur.
//
// res is the sql.Result.
func (l *LogUtility) LogResultRows(res sql.Result) {
	n, err := res.RowsAffected()
	if err != nil {
		l.log.Warnf("failed to check results: %v", err)
	} else {
		l.log.Infof("affected rows: %d", n)
	}
}

// logQueryAndArgs is used to log a SQL query and its arguments.
// Both will be logged to the Debug level.
func (l *LogUtility) LogQueryAndArgs(query string, args []any) {
	l.log.Debugf("query: %s", query)
	l.log.Debugf("args: %v", args)
}
