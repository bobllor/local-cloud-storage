package dbgateway

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/bobllor/gologger"
)

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
