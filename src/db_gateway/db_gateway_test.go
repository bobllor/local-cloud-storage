package dbgateway

import (
	"fmt"
	"testing"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/file"
	"github.com/bobllor/cloud-project/src/tests"
	"github.com/go-sql-driver/mysql"
)

func TestClauseDataBuildSetQuery(t *testing.T) {
	cd := ClauseData{
		Columns: []string{file.ColumnFileID, file.ColumnFileName},
		Args:    []any{tests.DbRowInfo.FileID, "a file text.txt"},
	}

	expectedQuery := fmt.Sprintf("SET %s = ?,%s = ?", file.ColumnFileID, file.ColumnFileName)

	setQ, err := cd.BuildSetQuery()
	assert.Nil(t, err)

	assert.Equal(t, setQ, expectedQuery)
}

// newTestDBConfig creates a test DB config for use in test environments.
func newTestDBConfig() *mysql.Config {
	port := "3307"

	user := "root"
	password := ""
	net := "tcp"
	addr := "127.0.0.1" + ":" + port
	dbName := "TestLocalCloudStorage"

	dbConfig := NewConfig(user, password, net, addr, dbName)

	return dbConfig
}
