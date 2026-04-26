package dbgateway

import (
	"database/sql"
	"fmt"

	"github.com/bobllor/cloud-project/src/file"
	"github.com/go-sql-driver/mysql"
)

const (
	dbDriver = "mysql"
)

type Gateway struct {
	File    *FileGateway
	User    *UserGateway
	Session *SessionGateway
}

// NewGateway creates a new Gateway containing the database gateways.
func NewGateway(fg *FileGateway, ug *UserGateway, sg *SessionGateway) *Gateway {
	gw := &Gateway{
		File:    fg,
		User:    ug,
		Session: sg,
	}

	return gw
}

// NewDatabase opens the SQL database and returns a sql.DB.
// It will return an error if any errors occur.
//
// The database is pinged in this call, if it fails then
// an error will be returned.
func NewDatabase(config *mysql.Config) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, config.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewConfig creates a new *mysql.Config from the arguments.
func NewConfig(user string, passwd string, net string, addr string, dbName string) *mysql.Config {
	c := mysql.NewConfig()

	c.User = user
	c.Passwd = passwd
	c.Net = net
	c.Addr = addr
	c.DBName = dbName

	c.AllowNativePasswords = true
	c.ParseTime = true

	return c
}

// NewClauseData creates a new ClauseData for building clauses for
// column preparations.
func NewClauseData() *ClauseData {
	return &ClauseData{}
}

// ClauseData is used to build placeholders for clauses. This is used for
// preparing columns.
type ClauseData struct {
	// Columns are any columns that the query is being performed on.
	Columns []string
	// Args are used for the columns. This must be the same size as the columns.
	Args []any
}

// AddColumns adds the columns to the data.
// Note that for every column that is added, it is expected to have an equivalent
// argument.
//
// It will return the number of columns added.
func (cd *ClauseData) AddColumns(columns ...string) int {
	cd.Columns = append(cd.Columns, columns...)

	return len(columns)
}

// AddArgs adds any given arguments to the data.
// Note that for every argument that is added, it is expected to have an equivalent
// column.
//
// It will return the number of arguments added.
func (cd *ClauseData) AddArgs(args ...any) int {
	cd.Args = append(cd.Args, args...)

	return len(args)
}

// BuildSetQuery builds the query for SET operations. It returns the query
// string: "SET Column1=value1, Column2=value2, ...".
// The arguments is also returned.
//
// An error will occur if the length of Args and Columns are not the same.
func (cd *ClauseData) BuildSetQuery() (string, []any, error) {
	err := cd.Validate()
	if err != nil {
		return "", nil, fmt.Errorf("failed to validate ClauseData: %v", err)
	}

	baseQuery := "SET"
	setPlaceholder := BuildSetPlaceholder(cd.Columns)
	query := baseQuery + " " + setPlaceholder

	return query, cd.Args, nil
}

// Validate is used to validate ClauseData. An error will be returned if
// it fails to validate.
func (cd *ClauseData) Validate() error {
	emptyError := "cannot have empty slice for %s"
	if len(cd.Columns) == 0 {
		return fmt.Errorf(emptyError, "Columns")
	}
	if len(cd.Args) == 0 {
		return fmt.Errorf(emptyError, "Args")
	}

	if len(cd.Args) != len(cd.Columns) {
		return fmt.Errorf("sizes columns and args are not equal (%d != %d)", len(cd.Columns), len(cd.Args))
	}

	return nil
}

// WhereCondition is used to build conditions for the WHERE clause.
type WhereCondition struct {
	// Column is the column name being used as the condition.
	Column string

	// Args is any arguments used with the condition. The size depends on which
	// ComparisonOperator is used, but there must be a minimum one argument.
	Args []any

	// ComparisonOperator is used to determine which condition to add to the clause.
	// The choice of operator will affect how many Args are used in the function.
	ComparisonOperator ComparisonOperator

	// LogicalOperator is the logical condition used to connect two clauses. Valid values are
	// AND or OR.
	LogicalOperator LogicalOperator
}

// exec executes a query on a database and returns the Result. This
// is only used for INSERT and UPDATE.
func execQuery(db *sql.DB, query string, args ...any) (sql.Result, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}

	res, err := tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction for %s: %v", file.TableName, err)
	}

	return res, err
}
