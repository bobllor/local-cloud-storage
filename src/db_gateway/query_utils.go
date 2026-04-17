package dbgateway

import (
	"fmt"
	"slices"
	"strings"
)

type ComparisonOperator string

const (
	Equal   ComparisonOperator = "="
	Gt      ComparisonOperator = ">"
	Gte     ComparisonOperator = ">="
	Lt      ComparisonOperator = "<"
	Lte     ComparisonOperator = "<="
	Ne      ComparisonOperator = "<>"
	Between ComparisonOperator = "BETWEEN"
	Like    ComparisonOperator = "LIKE"
	In      ComparisonOperator = "IN"
	Exists  ComparisonOperator = "EXISTS"
)

type LogicalOperator string

const (
	OperatorAnd LogicalOperator = "AND"
	OperatorOr  LogicalOperator = "OR"
)

// NewClauseBuilder creates a new ClauseBuilder.
func NewClauseBuilder() *ClauseBuilder {
	return &ClauseBuilder{}
}

// BuildPlaceholder builds the placeholder string for parameters that are passed into
// queries.
// It will return the string base on the placeholder and repeat counts: "(?,?,...),(...)..."
// This does not include VALUES.
//
// placholderCount is the amount of placeholders that are being used. This is expected
// to be the number of parameters being used.
//
// repeat is how many times to repeat the final placeholder string. This is used for batch operations.
// If batch operations are not needed then it is expected to be 1. If less than 1 is given,
// that it will automatically be converted into 1.
func BuildPlaceholder(placeholderCount int, repeat int) string {
	questions := []string{}
	out := []string{}

	if repeat < 1 {
		repeat = 1
	}

	for range placeholderCount {
		questions = append(questions, "?")
	}

	param := "(" + strings.Join(questions, ",") + ")"

	for range repeat {
		out = append(out, param)
	}

	return strings.Join(out, ",")
}

// BuildSetPlaceholder builds the strings for updating columns in a table.
// The output will be in the form of: "Column1 = value,Column2 = value, ..."
func BuildSetPlaceholder(columns []string) string {
	placeholders := []string{}

	for _, column := range columns {
		placeholder := fmt.Sprintf("%s = ?", column)

		placeholders = append(placeholders, placeholder)
	}

	return strings.Join(placeholders, ",")
}

// ClauseBuilder is a helper in building dynamic WHERE conditions.
//
// When a new clause is registered, the method call will return the same
// ClauseBuilder. This can be used to chain multiple conditions with AND/OR.
// However, it must be correct otherwise when building the conditions it will fail.
type ClauseBuilder struct {
	// clauses are SQL clauses used in the WHERE condition. The operator
	// conditions (AND/OR) are also added in here.
	clauses []string

	// args are the arg parameters for the clause. It is all the args in one
	// and is not equal to the size of clauses, due to a single clause potentially
	// having multiple parameters. The args order is the same as when
	// the clause was added.
	args []any
}

// Build combines the registered clauses and returns the clause and
// a flatten slice of all arguments in order of the clauses.
//
// The return query includes the WHERE.
//
// The clauses will be validated prior to building, if the clauses
// are not built correctly then an error will be returned.
func (c *ClauseBuilder) Build() (string, []any, error) {
	if len(c.clauses) == 0 {
		return "", nil, fmt.Errorf("clause validation failed: cannot build an empty clause")
	}
	err := c.validateEndClause()
	if err != nil {
		return "", nil, err
	}
	err = c.validateClauses()
	if err != nil {
		return "", nil, err
	}

	clauseHolder := []string{}

	clauseHolder = append(clauseHolder, c.clauses...)

	clause := "WHERE " + strings.Join(clauseHolder, " ")

	return clause, c.args, nil
}

// Equal registers a new equal (=) clause to the builder.
func (c *ClauseBuilder) Equal(column string, arg any) *ClauseBuilder {
	clause := fmt.Sprintf("%s %s ?", column, Equal)

	c.register(clause, arg)

	return c
}

// In registers a new IN clause to the builder. This is used for
// checking values in a list for SQL.
func (c *ClauseBuilder) In(column string, args ...any) *ClauseBuilder {
	params := BuildPlaceholder(len(args), 1)

	clause := fmt.Sprintf("%s %s %s", column, In, params)

	c.register(clause, args...)

	return c
}

// RegisterConditions registers WHERE conditions into the ClauseBuilder.
//
// If conditions are registered when no clauses are registered, the first condition
// logical operator will be ignored.
func (cb *ClauseBuilder) RegisterConditions(conditions []WhereCondition) error {
	for _, condition := range conditions {
		// ignores this if empty as it would be invalid
		// logical operators in the first slot will be caught by cb.Build
		if len(cb.args) != 0 {
			switch logicOp := condition.LogicalOperator; logicOp {
			case OperatorAnd:
				cb.And()
			case OperatorOr:
				cb.Or()
			default:
				return fmt.Errorf("logical operator %s not supported", logicOp)
			}
		}

		switch compOp := condition.ComparisonOperator; compOp {
		case Equal:
			cb.Equal(condition.Column, condition.Args[0])
		case In:
			cb.In(condition.Column, condition.Args...)
		default:
			return fmt.Errorf("comparison operator %s not supported", compOp)
		}
	}

	return nil
}

// And adds an AND operator to the clauses to prepare the next
// condition.
func (c *ClauseBuilder) And() *ClauseBuilder {
	c.addOperator(OperatorAnd)

	return c
}

// Or adds an OR operator to the clauses to prepare the next
// condition.
func (c *ClauseBuilder) Or() *ClauseBuilder {
	c.addOperator(OperatorOr)

	return c
}

// register is a helper function that registers the clause and args to
// the equivalent fields.
func (c *ClauseBuilder) register(clause string, args ...any) {
	c.clauses = append(c.clauses, clause)
	c.args = append(c.args, args...)
}

// addOperator adds an operator to clauses to prepare the next
// clause condition.
func (c *ClauseBuilder) addOperator(operator LogicalOperator) {
	c.clauses = append(c.clauses, string(operator))
}

// validateClauses checks if a clause is valid. This validates
// if the clauses and operators are added properly.
// It does not do input validation on queries.
//
// It will return an error if it is not valid.
func (c *ClauseBuilder) validateClauses() error {
	operators := []LogicalOperator{OperatorAnd, OperatorOr}
	mustBeOp := false

	for _, clause := range c.clauses {
		if mustBeOp && !slices.Contains(operators, LogicalOperator(clause)) {
			return fmt.Errorf("expected operator (AND/OR) but got clause %s (%v)", clause, c.clauses)
		} else if !mustBeOp && slices.Contains(operators, LogicalOperator(clause)) {
			return fmt.Errorf("expected non-operator clause but got %s (%v)", clause, c.clauses)
		}

		if mustBeOp {
			mustBeOp = false
		} else {
			mustBeOp = true
		}
	}

	return nil
}

// validateEndClause checks if the ending of the clauses is an operator.
// This is used to drop operators without a follow-up clause.
// When a normal clause is encountered, then it will stop the search.
func (c *ClauseBuilder) validateEndClause() error {
	operators := []LogicalOperator{OperatorAnd, OperatorOr}

	for lp := len(c.clauses) - 1; lp > -1; lp-- {
		clause := c.clauses[lp]

		if slices.Contains(operators, LogicalOperator(clause)) {
			return fmt.Errorf("end clauses cannot end in an operator: got %s (%v)", clause, c.clauses)
		} else {
			break
		}
	}

	return nil
}
