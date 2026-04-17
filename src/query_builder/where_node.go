package querybuilder

import "fmt"

// WhereNode is a node representing a single condition. It contains an
// internal node pointer that is used to create the condition chain
// for a query.
type WhereNode struct {
	// column is the column or subquery that is being used.
	column string

	// args is a slice of any type.
	args []any

	// operator is the operator used for conditionals.
	operator ComparisonOperator

	// next is the next WhereNode for chaining.
	next *WhereNode
}

// NewWhereNode creates a new WhereNode.
func NewWhereNode(column string, compOperator ComparisonOperator, args ...any) *WhereNode {
	return &WhereNode{
		column:   column,
		operator: compOperator,
		args:     args,
	}
}

// SetColumn sets the column value. If a column value is not needed,
// then an empty string can be set.
func (w *WhereNode) SetColumn(column string) {
	w.column = column
}

// AppendArgs appends any args into the arg values.
func (w *WhereNode) AppendArgs(args ...any) {
	for _, arg := range args {
		w.args = append(w.args, arg)
	}
}

// Build builds the query string for the node. This does not include the WHERE
// clause, it only contains the conditions.
func (w *WhereNode) Build() string {
	query := w.buildQuery()

	return query
}

// buildQuery builds the query string for the node and handles
// operator specific strings.
func (w *WhereNode) buildQuery() string {
	params := BuildPlaceholder(len(w.args), 1)
	var query string
	switch w.operator {
	case OperatorIn:
		query = fmt.Sprintf("%s IN %s", w.column, params)
	case OperatorEqual:
		query = fmt.Sprintf("%s = %s", w.column, params)
	case OperatorExists:
		// special case, the subquery is considered the column
		// args are not used here, args is added in the WhereClause
		query = fmt.Sprintf("EXISTS (%s)", w.column)
	}

	return query
}
