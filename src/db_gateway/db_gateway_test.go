package dbgateway

import (
	"strings"
	"testing"

	"github.com/bobllor/assert"
)

func TestMultiBuildPlaceholder(t *testing.T) {
	params := 8
	repeat := 3

	query := BuildPlaceholder(params, repeat)
	querySplit := strings.Split(query, ")")

	// has to subtract -1 due to an invisible string at the end.
	assert.Equal(t, len(querySplit)-1, repeat)

	assert.Equal(t, len(strings.Split(querySplit[0], ",")), params)
}

func TestSingleBuildPlaceholder(t *testing.T) {
	params := 5
	repeat := 1

	query := BuildPlaceholder(params, repeat)
	spl := strings.Split(query, ")")

	assert.Equal(t, len(spl)-1, repeat)
	assert.Equal(t, len(strings.Split(query, ",")), params)
}
