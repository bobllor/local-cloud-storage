package tests

import (
	"strings"
	"testing"

	"github.com/bobllor/assert"
	dbcon "github.com/bobllor/cloud-project/src/db_con"
)

func TestQueryParamBuilder(t *testing.T) {
	params := 8
	repeat := 3

	query := dbcon.QueryParamBuilder(params, repeat)
	querySplit := strings.Split(query, ")")

	// has to subtract -1 due to an invisible string at the end.
	assert.Equal(t, len(querySplit)-1, repeat)

	assert.Equal(t, len(strings.Split(querySplit[0], ",")), params)
}
