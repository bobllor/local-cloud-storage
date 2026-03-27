package dbgateway

import (
	"strings"
	"testing"

	"github.com/bobllor/assert"
)

func TestQueryParamBuilder(t *testing.T) {
	params := 8
	repeat := 3

	query := QueryParamBuilder(params, repeat)
	querySplit := strings.Split(query, ")")

	// has to subtract -1 due to an invisible string at the end.
	assert.Equal(t, len(querySplit)-1, repeat)

	assert.Equal(t, len(strings.Split(querySplit[0], ",")), params)
}
