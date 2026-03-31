package dbgateway

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bobllor/cloud-project/src/file"
	"github.com/go-sql-driver/mysql"
)

const (
	dbDriver = "mysql"
)

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

	return c
}

// NewBatcher creates a new Batcher for performing batch operations.
//
// fileOwnerID is the target owner of the File rows being modified. This is required.
func NewBatcher(fileOwnerID string) (*Batcher, error) {
	if strings.TrimSpace(fileOwnerID) == "" {
		return nil, fmt.Errorf("cannot have an empty string for the file owner ID")
	}

	sb := &Batcher{
		FileOwnerID: fileOwnerID,
	}

	return sb, nil
}

// exec executes a query on a database and returns the Result. This
// is only used for INSERT and UPDATE.
func execQuery(db *sql.DB, query string, args ...any) (sql.Result, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	res, err := tx.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %v", query, err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction for %s: %v", file.FileTableName, err)
	}

	return res, err
}

// Batcher is used to build queries for batch operations. Although it is intended to be used
// for batching, it can be used for single rows if needed.
type Batcher struct {
	// FileOwnerID is the ID owner of the rows in the File table being queried for.
	FileOwnerID string
	// SetClauseData is used to build the placeholders for the SET clause. This is optional
	// and can be left empty if SET is not being used.
	SetClauseData SetClauseData
	// Conditions is a slice that is used to create conditions with the WHERE clause.
	Conditions []WhereCondition
}

// AddSetClauseData adds the data to Batcher for the SET clause.
func (b *Batcher) AddSetClauseData(setClauseData SetClauseData) {
	b.SetClauseData = setClauseData
}

// AddWhereConditions adds the data to Batcher for the WhereConditions for the WHERE clause.
func (b *Batcher) AddWhereConditions(conditions []WhereCondition) {
	b.Conditions = conditions
}

// BuildSet builds the full query for SET operations. This includes any of the
// WhereConditions if given.
func (b *Batcher) BuildSetQuery(tableName string) (string, error) {
	err := b.SetClauseData.Validate()
	if err != nil {
		return "", fmt.Errorf("failed to validate SetInfo: %v", err)
	}

	baseQuery := fmt.Sprintf("UPDATE %s SET", tableName)
	setPlaceholder := BuildSetPlaceholder(b.SetClauseData.Columns)

	cb := NewClauseBuilder()

	err = cb.RegisterBatcher(b)
	if err != nil {
		return "", err
	}

	query := baseQuery + " " + setPlaceholder

	return query, nil
}

// SetClauseData is used to build the placeholders for the SET clause. This is optional
// and can be left empty if SET is not being used.
type SetClauseData struct {
	// Columns are any columns that the query is being performed on.
	Columns []string
	// Args are used for the columns. This must be the same size as the columns.
	Args []any
}

// Validate is used to validate SetInfo data. An error will be returned if
// it fails to validate.
func (si *SetClauseData) Validate() error {
	emptyError := "cannot have empty slice for %s"
	if len(si.Columns) == 0 {
		return fmt.Errorf(emptyError, "Columns")
	}
	if len(si.Args) == 0 {
		return fmt.Errorf(emptyError, "Args")
	}

	if len(si.Args) != len(si.Columns) {
		return fmt.Errorf("sizes columns and args are not equal (%d != %d)", len(si.Columns), len(si.Args))
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
