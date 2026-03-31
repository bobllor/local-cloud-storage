package dbgateway

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/file"
)

func TestBuildSingleQuery(t *testing.T) {
	cb := ClauseBuilder{}

	arg1 := "afileid"
	expectedArgs := 1
	expectedQuery := fmt.Sprintf("WHERE %s = ?", file.FileIDCol)

	cb.Equal(file.FileIDCol, arg1)

	q, args, err := cb.Build()
	assert.Nil(t, err)

	assert.Equal(t, len(args), expectedArgs)
	assert.Equal(t, q, expectedQuery)
}

func TestBuildAndQuery(t *testing.T) {
	cb := ClauseBuilder{}

	expectedArgs := 2

	arg1 := "fileidhere"
	arg2 := "filename.txt"

	expectedQuery := fmt.Sprintf("WHERE %s = ? AND %s = ?", file.FileIDCol, file.FileNameCol)

	cb.Equal(file.FileIDCol, arg1).And().Equal(file.FileNameCol, arg2)

	q, args, err := cb.Build()
	assert.Nil(t, err)

	assert.Equal(t, len(args), expectedArgs)
	assert.Equal(t, args[0], arg1)
	assert.Equal(t, args[1], arg2)

	assert.Equal(t, q, expectedQuery)
}

func TestInQuery(t *testing.T) {
	cb := NewClauseBuilder()

	expectedArgs := 4
	expectedQuery := fmt.Sprintf("WHERE %s = ? AND %s IN (?,?,?)", file.FileOwnerIDCol, file.FileNameCol)

	cb.Equal(file.FileOwnerIDCol, "file-owner").And().In(file.FileNameCol, "test1.txt", "test2.txt", "test3.txt")

	q, args, err := cb.Build()
	assert.Nil(t, err)

	assert.Equal(t, len(args), expectedArgs)
	assert.Equal(t, q, expectedQuery)
}

func TestEndOperatorClauseError(t *testing.T) {
	cb := NewClauseBuilder()

	cb.Equal(file.FileOwnerIDCol, "fileidhere").
		And().
		Equal(file.FileNameCol, "filenamehere.txt").
		And().And().Or()

	_, _, err := cb.Build()
	assert.NotNil(t, err)
}

func TestEmptyClauseError(t *testing.T) {
	cb := NewClauseBuilder()

	_, _, err := cb.Build()
	assert.NotNil(t, err)
}

func TestSingleBuildPlaceholder(t *testing.T) {
	params := 5
	repeat := 1

	query := BuildPlaceholder(params, repeat)
	spl := strings.Split(query, ")")

	assert.Equal(t, len(spl)-1, repeat)
	assert.Equal(t, len(strings.Split(query, ",")), params)
}

func TestMultiBuildPlaceholder(t *testing.T) {
	params := 8
	repeat := 3

	query := BuildPlaceholder(params, repeat)
	querySplit := strings.Split(query, ")")

	// has to subtract -1 due to an invisible string at the end.
	assert.Equal(t, len(querySplit)-1, repeat)

	assert.Equal(t, len(strings.Split(querySplit[0], ",")), params)
}

func TestSingleSetPlaceholder(t *testing.T) {
	columns := []string{file.FileNameCol}

	query := BuildSetPlaceholder(columns)

	querySplit := strings.Split(query, ",")

	assert.Equal(t, len(querySplit), len(columns))
	assert.Equal(t, strings.Contains(query, columns[0]), true)

	counter := 0

	for _, ch := range query {
		if ch == '?' {
			counter += 1
		}
	}

	assert.Equal(t, counter, len(columns))
}

func TestMultiSetPlaceholder(t *testing.T) {
	columns := []string{file.FileNameCol, file.FileIDCol, file.FileSizeCol}

	query := BuildSetPlaceholder(columns)

	querySplit := strings.Split(query, ",")

	assert.Equal(t, len(querySplit), len(columns))

	for _, col := range columns {
		assert.Equal(t, strings.Contains(query, col), true)
	}

	counter := 0

	for _, ch := range query {
		if ch == '?' {
			counter += 1
		}
	}

	assert.Equal(t, counter, len(columns))
}
