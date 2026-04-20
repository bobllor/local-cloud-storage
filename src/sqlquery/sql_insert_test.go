package sqlquery

import (
	"fmt"
	"testing"
	"time"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/file"
)

func TestInsertBuildColumns(t *testing.T) {
	query, args, err := InsertInto(file.TableName, file.ColumnDeletedOn, file.ColumnFileID).
		Args(time.Now(), "12345", time.Now(), "12345").
		Build()
	assert.Nil(t, err)

	params := BuildPlaceholder(2, 2)

	baseQuery := fmt.Sprintf(
		"INSERT INTO %s (%s,%s) VALUES %s",
		file.TableName,
		file.ColumnDeletedOn,
		file.ColumnFileID,
		params,
	)

	assert.Equal(t, len(args), 4)
	assert.Equal(t, query, baseQuery)
}

func TestInsertIntoAllColumns(t *testing.T) {
	query, args, err := InsertInto(file.TableName, file.ColumnDeletedOn, file.ColumnFileID).
		Args(time.Now(), "12345", time.Now(), "12345").
		Build()
	assert.Nil(t, err)

	params := BuildPlaceholder(2, 2)

	baseQuery := fmt.Sprintf(
		"INSERT INTO %s (%s,%s) VALUES %s",
		file.TableName,
		file.ColumnDeletedOn,
		file.ColumnFileID,
		params,
	)

	assert.Equal(t, len(args), 4)
	assert.Equal(t, query, baseQuery)
}

func TestInsertIntoInvalidBuild(t *testing.T) {

}
