package querybuilder

import (
	"fmt"
	"strings"
)

// WhereClause is a struct that is used to create WHERE clauses, supporting
// chaining and grouping of conditions.
type WhereClause struct {
	// builder is a Builder interface that can build a query string.
	builder Builder

	// node holds the group of conditions and its values.
	node *WhereNode

	// head is the pointer to the first node in the WhereClause. This should
	// never be modified or changed, as it is only used as a reference.
	head *WhereNode

	// logicalOperators is a slice of logical operators (AND/OR) for use with
	// multiple conditions.
	logicalOperators []LogicalOperator
}

// NewWhereClause creates a new WhereClause.
func NewWhereClause(b Builder) *WhereClause {
	return &WhereClause{
		builder: b,
	}
}

// Build builds a string by adding the WHERE query to the parent query.
func (w *WhereClause) Build() string {
	mQuery := w.builder.Build()
	wQuery := w.buildWhereString()

	query := fmt.Sprintf("%s WHERE %s", mQuery, wQuery)

	return query
}

// Equal creates a new "equal to" condition. Example: `Column = value`.
func (w *WhereClause) Equal(column string, arg any) *WhereClause {
	w.addNewNode(column, OperatorEqual, arg)

	return w
}

// In creates a new "in (list)" condition. Example: `Column IN (val1, val2...)`.
func (w *WhereClause) In(column string, args ...any) *WhereClause {
	w.addNewNode(column, OperatorIn, args...)

	return w
}

// Exists creates a new "exists (subquery)" condition. This is a special case in that
// an existing query must be passed along with any arguments it comes with.
//
// The arguments are appended to the w.Args.
func (w *WhereClause) Exists(subquery string, args ...any) *WhereClause {
	w.addNewNode(subquery, OperatorExists, args...)

	return w
}

// addNewNode adds a new node to the WhereNode. This handles both initially nil values,
// and existing values with nodes to maintain the linked list chain.
//
// Arguments are also appended to w.args.
func (w *WhereClause) addNewNode(column string, operator ComparisonOperator, args ...any) {
	node := NewWhereNode(column, operator, args...)
	if w.node == nil {
		w.node = node
		w.head = w.node
	} else {
		w.node.next = node
		w.node = w.node.next
	}

	w.builder.WriteArgs(args...)
}

// And adds the AND logical operator to start a new condition. This
// expects a previous condition to exist.
func (w *WhereClause) And() *WhereClause {
	w.logicalOperators = append(w.logicalOperators, OperatorAnd)

	return w
}

// Or adds the OR logical operator to start a new condition. This
// expects a previous condition to exist.
func (w *WhereClause) Or() *WhereClause {
	w.logicalOperators = append(w.logicalOperators, OperatorOr)

	return w
}

// String creates a string representation of the WhereNode data.
func (w *WhereClause) String() string {
	temp := w.head

	out := []string{}

	for temp != nil {
		s := fmt.Sprintf(
			"[column:'%s',args:%v,operator:'%v',hasNext:%v]",
			temp.column,
			temp.args,
			temp.operator,
			temp.next != nil,
		)

		temp = temp.next
		out = append(out, s)
	}

	return strings.Join(out, ";")
}

// buildWhereString builds the WHERE string by accessing the nodes.
func (w *WhereClause) buildWhereString() string {
	temp := w.head
	isFirstNode := true
	logicOpIndex := 0

	out := []string{}

	for temp != nil {
		if !isFirstNode && logicOpIndex < len(w.logicalOperators) {
			logicOp := w.logicalOperators[logicOpIndex]
			out = append(out, string(logicOp))

			logicOpIndex += 1
		} else if isFirstNode {
			isFirstNode = false
		}

		s := temp.Build()
		temp = temp.next

		out = append(out, s)
	}

	return strings.Join(out, " ")
}
