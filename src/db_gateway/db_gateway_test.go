package dbgateway

import (
	"fmt"
	"testing"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/file"
)

func TestClauseDataBuildSetQuery(t *testing.T) {
	cd := ClauseData{
		Columns: []string{file.FileIDCol, file.FileNameCol},
		Args:    []any{testFileID, "a file text.txt"},
	}

	expectedQuery := fmt.Sprintf("SET %s = ?,%s = ?", file.FileIDCol, file.FileNameCol)

	setQ, err := cd.BuildSetQuery()
	assert.Nil(t, err)

	assert.Equal(t, setQ, expectedQuery)
}
