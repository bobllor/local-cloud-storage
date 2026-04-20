package sqlquery

import "strings"

type DMLType string

const (
	SelectDML     DMLType = "SELECT"
	DeleteDML     DMLType = "DELETE"
	UpdateDML     DMLType = "UPDATE"
	InsertIntoDML DMLType = "INSERT INTO"
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

// Builder is an interface used to build SQL queries.
type Builder interface {
	// Build builds the full query string. It will return the
	// query string, its arguments if any, and an error if one occurs
	Build() (string, []any, error)

	// Write writes to the args slice.
	Write(args ...any)
}

// SqlArgs is adds arguments to the builder after the columns are added.
type SqlArgs struct {
	builder Builder
}

// QueryBuilder is used to build a SQL query.
//
// It uses the elements for building the SQL query:
//   - Main: DML + columns + FROM + table
//   - Join: JOIN + table + ON + condition
//   - Filter: WHERE + condition
//   - Group: GROUP BY + column
//   - Group filter: HAVING + condition
//   - Sort: ORDER BY + condition
//
// Main is passed as a string query to QueryBuilder.Build.
type QueryBuilder struct {
	FilterBuilder *ConditionClause
}

// Build builds the query and returns the query string. Validation
// errors will return an error if it fails to validate.
//
// mainQuery is the starting query string from the main builder.
// It is expected to consist of a DML, columns, data (if applicable), and
// a table.
func (qb *QueryBuilder) Build(mainQuery string) (string, error) {
	holder := []string{mainQuery}

	// TODO: create a validation checker for the query build:
	//	- no table name
	//	- invalid main query string

	if qb.FilterBuilder != nil {
		holder = append(holder, qb.FilterBuilder.query)
	}

	return strings.Join(holder, " "), nil
}

// Args adds the given arguments into the builder used for
// the parameters of a query for columns. The builder is returned.
func (sa *SqlArgs) Args(args ...any) Builder {
	sa.builder.Write(args...)

	return sa.builder
}
