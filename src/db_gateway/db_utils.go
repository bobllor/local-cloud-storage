package dbgateway

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/bobllor/gologger"
)

// SelectRow iterates over rows to fill a given source of any type.
// This only selects a single row from the database.
// src must be a pointer to a non-slice type.
//
// The columns length found in rows must be equal to the number of fields given
// in src, and the order of the columns must match the order of the type in src.
//
// If there are no entries in rows then the src interface will be consist of nil values.
// This must be checked for on the parent caller. Otherwise, use SelectRows for
// slice length checking.
func SelectRow(rows *sql.Rows, src interface{}) error {
	v := reflect.ValueOf(src)

	if rows == nil {
		return errors.New("rows cannot be nil")
	}
	if v.Kind() != reflect.Ptr || v.Elem().Type().Kind() == reflect.Slice {
		return errors.New("src interface must be a pointer to a non-slice")
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if v.Elem().NumField() != len(columns) {
		return fmt.Errorf(
			"number of columns in query (%d) does not match interface fields (%d)",
			len(columns),
			v.Elem().NumField(),
		)
	}

	for rows.Next() {
		// setting the slice type to the values for scanning
		values := getScannableValues(len(columns), v)

		err := rows.Scan(values...)
		if err != nil {
			return err
		}
	}

	return nil
}

// SelectRows iterates over rows to fill a given source slice of any type.
// src must be a pointer to a slice.
//
// The columns length found in rows must be equal to the number of fields given
// in src, and the order of the columns must match the order of the type in src.
func SelectRows(rows *sql.Rows, src interface{}) error {
	v := reflect.ValueOf(src)

	if rows == nil {
		return errors.New("rows cannot be nil")
	}
	if v.Kind() != reflect.Ptr || v.Elem().Type().Kind() != reflect.Slice {
		return errors.New("src interface must be a pointer to a slice")
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// the type of the slice, ptr slice any -> slice any -> any
	t := reflect.TypeOf(src).Elem().Elem()
	ve := v.Elem()

	if t.NumField() != len(columns) {
		return fmt.Errorf(
			"number of columns in query (%d) does not match interface fields (%d)",
			len(columns),
			t.NumField(),
		)
	}

	for rows.Next() {
		// setting the slice type to the values for scanning
		ele := reflect.New(t)
		values := getScannableValues(len(columns), ele)

		err := rows.Scan(values...)
		if err != nil {
			return err
		}

		// ve is a dereferenced v, this ensures we append to the same
		// src value
		ve.Set(reflect.Append(ve, ele.Elem()))
	}

	return nil
}

// UpdateRow updates a single row from a table based on its column value.
func UpdateRow(db *sql.DB, table string, whereColumn string, whereArg any, clause ClauseData) (sql.Result, error) {
	cb := NewClauseBuilder()
	cb.In(whereColumn, whereArg)

	clQ, sargs, err := clause.BuildSetQuery()
	if err != nil {
		return nil, err
	}

	baseQ := fmt.Sprintf("UPDATE %s %s", table, clQ)

	cbQ, args, err := cb.Build()
	if err != nil {
		return nil, err
	}

	execArgs := MakeArgs(sargs, args)

	query := baseQ + " " + cbQ

	res, err := execQuery(db, query, execArgs...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DropRows is used to drop rows from a table. This is a destructive action and will
// permanently delete the row. This should only be used for very specific cases or testing.
//
// column is used to target the column where the row is in the given slice of args.
func DropRows(db *sql.DB, table string, column string, args ...any) (sql.Result, error) {
	cb := NewClauseBuilder()

	cb.In(column, args...)

	cbQ, newArgs, err := cb.Build()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("DELETE FROM %s", table) + " " + cbQ

	res, err := execQuery(db, query, newArgs...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// getScannableValues creates a new []any that contains the addresses
// of the given reflect.Value fields.
//
// The value is expected to be a pointer type. It will panic if
// it is not a pointer to a non-slice.
func getScannableValues(size int, v reflect.Value) []any {
	values := make([]any, size)

	fieldCount := v.Elem().NumField()

	for i := range fieldCount {
		field := v.Elem().Field(i)

		values[i] = field.Addr().Interface()
	}

	return values
}

// AppendArgs appends to s []any of any arguments.
//
// This is used to build onto an existing []any for
// query args.
func AppendArgs(s *[]any, arg ...any) {
	for _, v := range arg {
		*s = append(*s, v)
	}
}

// MakeArgs creates a new any slice from any given arguments.
func MakeArgs(args ...any) []any {
	out := []any{}

	for _, a := range args {
		av := reflect.ValueOf(a)
		v := getReflectValue(av)

		if v.Kind() != reflect.Slice {
			out = append(out, a)
		} else {
			for i := 0; i < v.Len(); i++ {
				vv := v.Index(i)
				out = append(out, vv.Interface())
			}
		}
	}

	return out
}

// getReflectValue is a recursive call that retrieves the underlying value
// of v.
//
// If v is not a pointer, then it will return v.
func getReflectValue(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Pointer {
		return v
	}

	return getReflectValue(v.Elem())
}

// DBUtility is a utility struct for database related operations.
type DBUtility struct {
	log *gologger.Logger
}

// LogResultRows is used to log the SQL result rows. If an error is encountered during
// the result parsing, then it will not log the affected rows.
// This will log to the Info level with no error, and the Warn level if an error does occur.
//
// res is the sql.Result.
func (d *DBUtility) LogResultRows(res sql.Result) {
	n, err := res.RowsAffected()
	if err != nil {
		d.log.Warnf("failed to check results: %v", err)
	} else {
		d.log.Infof("affected rows: %d", n)
	}
}

// logQueryAndArgs is used to log a SQL query and its arguments.
// Both will be logged to the Debug level.
func (d *DBUtility) LogQueryAndArgs(query string, args []any) {
	d.log.Debugf("query: %s", query)
	d.log.Debugf("args: %v", args)
}
