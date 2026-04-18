package querybuilder

type DMLType string

const (
	Select DMLType = "SELECT"
	Drop   DMLType = "DROP"
	Update DMLType = "UPDATE"
	Insert DMLType = "INSERT"
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
	Build() string

	// WriteArgs writes to an any slice with the given arguments.
	WriteArgs(args ...any)
}

type SqlBuilder struct {
	// table is the table name that the operations are performed on.
	table string

	// args is a slice used for query parameters. New arguments are
	// appended on for any method that is adding arguments.
	args []any
}

// NewSqlBuilder creates a new SQL builder for building queries.
func NewSqlBuilder(table string) *SqlBuilder {
	return &SqlBuilder{
		table: table,
	}
}

func (sb *SqlBuilder) Select() *SqlSelect {
	s := &SqlSelect{
		builder: sb,
	}

	return s
}

func (sb *SqlBuilder) Insert(paramAmount int) *SqlInsert {
	si := &SqlInsert{
		builder:     sb,
		paramAmount: paramAmount,
	}

	return si
}

// Args returns the any slice used for query arguments.
func (sb *SqlBuilder) Args() []any {
	return sb.args
}
