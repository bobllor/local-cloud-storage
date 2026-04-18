package querybuilder

import "fmt"

// SqlInsert is used to build INSERT queries.
type SqlInsert struct {
	builder *SqlBuilder

	// columns are the columns that are being inserted into the
	// row. It is not required, by default queries will
	// attempt to insert into all rows in order.
	columns []string

	// paramAmount is the number of times the param query is repeated.
	// This is used for batch insertions. For example, if paramAmount = 3
	// then there will be 3 groups of params: "(?),(?),(?)".
	paramAmount int
}

func (s *SqlInsert) Columns(columns ...string) *SqlInsert {
	s.columns = append(s.columns, columns...)

	return s
}

// Args adds arguments to use for the columns.
// The amount of columns used is what is inserted into with args. If
// there are more args than columns, then the excess args will be
// used in the next batch.
//
// It is recommended that the len(args) % len(columns) == 0 for clean
// insertion, otherwise unexpected errors can occur.
func (s *SqlInsert) Args(args ...any) *SqlInsert {
	s.builder.args = append(s.builder.args, args...)

	return s
}

// Build creates the query for the INSERT statement.
func (s *SqlInsert) Build() string {
	params := BuildPlaceholder(len(s.columns), s.paramAmount)

	query := fmt.Sprintf(
		"INSERT INTO %s VALUES %s",
		s.builder.table,
		params,
	)

	return query
}
