package querybuilder

type SqlInsert struct {
	// columns are the columns that are being inserted into the
	// row. It is not required, by default queries will
	// attempt to insert into all rows in order.
	columns []string
	args    []any
}
