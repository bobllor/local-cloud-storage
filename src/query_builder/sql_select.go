package querybuilder

import (
	"fmt"
	"strings"
)

// SqlSelect is a struct to perform SELECT queries.
type SqlSelect struct {
	builder *SqlBuilder

	// columns are the columns that are being selected. This
	// must have a value.
	columns []string

	// where is used to build WHERE conditionals.
	// Initially this will be nil until the method Where is
	// explictly called in order to create the conditions.
	where *WhereClause
}

// Build builds string for the base SELECT query.
func (s *SqlSelect) Build() string {
	query := fmt.Sprintf(
		"%s %s FROM %s",
		Select,
		strings.Join(s.columns, ","),
		s.builder.table,
	)

	return query
}

// Columns adds the columns to select in the query.
func (s *SqlSelect) Columns(columns ...string) *SqlSelect {
	for _, column := range columns {
		s.columns = append(s.columns, column)
	}

	return s
}

// Where creates a new WhereClause for conditional usage.
func (s *SqlSelect) Where() *WhereClause {
	s.where = NewWhereClause(s)

	return s.where
}

func (s *SqlSelect) WriteArgs(args ...any) {
	s.builder.args = append(s.builder.args, args...)
}
