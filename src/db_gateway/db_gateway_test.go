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

func TestSingleQueryParamBuilder(t *testing.T) {
	params := 5
	repeat := 1

	query := QueryParamBuilder(params, repeat)
	spl := strings.Split(query, ")")

	assert.Equal(t, len(spl)-1, repeat)
	assert.Equal(t, len(strings.Split(query, ",")), params)
}
