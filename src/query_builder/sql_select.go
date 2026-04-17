package querybuilder

import (
	"fmt"
	"strings"
)

type ComparisonOperator string

const (
	OperatorEqual   ComparisonOperator = "="
	OperatorGt      ComparisonOperator = ">"
	OperatorGte     ComparisonOperator = ">="
	OperatorLt      ComparisonOperator = "<"
	OperatorLte     ComparisonOperator = "<="
	OperatorNe      ComparisonOperator = "<>"
	OperatorBetween ComparisonOperator = "BETWEEN"
	OperatorLike    ComparisonOperator = "LIKE"
	OperatorIn      ComparisonOperator = "IN"
	OperatorExists  ComparisonOperator = "EXISTS"
)

type LogicalOperator string

const (
	OperatorAnd LogicalOperator = "AND"
	OperatorOr  LogicalOperator = "OR"
)

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

// Column adds a single column to select in the query.
func (s *SqlSelect) Column(column string) *SqlSelect {
	s.columns = append(s.columns, column)

	return s
}

// Columns adds the columns to select in the query.
func (s *SqlSelect) Columns(columns ...string) *SqlSelect {
	for _, column := range columns {
		s.columns = append(s.columns, column)
	}

	return s
}

// AllColumns selects all columns of the given table.
func (s *SqlSelect) AllColumns() *SqlSelect {
	s.columns = append(s.columns, "*")

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
