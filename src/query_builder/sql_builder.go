package querybuilder

type DMLType string

const (
	Select DMLType = "SELECT"
	Drop   DMLType = "DROP"
	Update DMLType = "UPDATE"
	Insert DMLType = "INSERT"
)

// Builder is an interface used to build SQL queries.
type Builder interface {
	Build() string

	// WriteArgs writes to the []any slice with the given arguments.
	WriteArgs(args ...any)
}

type SqlBuilder struct {
	// table is the table name that the operations are performed on.
	table string

	// args is a slice of arguments added to the query, added during
	// building specific queries such as WHERE conditions or UPDATE columns.
	args []any
}

// NewSqlBuilder creates a new SQL builder for building queries.
func NewSqlBuilder(table string) *SqlBuilder {
	return &SqlBuilder{
		table: table,
	}
}

func (sb *SqlBuilder) Select() *SqlSelect {
	return &SqlSelect{
		builder: sb,
	}
}

// Args returns the any slice used for query arguments.
func (sb *SqlBuilder) Args() []any {
	return sb.args
}
