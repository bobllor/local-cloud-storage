package dbgateway

import (
	"fmt"
	"testing"
	"time"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/file"
	"github.com/bobllor/cloud-project/src/tests"
	"github.com/bobllor/cloud-project/src/user"
)

func TestSelectRow(t *testing.T) {
	dbConfig := newTestDBConfig()

	db, err := NewDatabase(dbConfig)
	assert.Nil(t, err)

	query := fmt.Sprintf(
		"SELECT %s, %s, %s FROM %s WHERE %s = ?",
		file.ColumnFileName,
		file.ColumnFileID,
		file.ColumnFileSize,
		file.TableName,
		file.ColumnFileOwnerID,
	)

	type FileFilterTest struct {
		FileName string
		FileID   string
		FileSize int
	}

	ffTest := FileFilterTest{}

	rows, err := db.Query(query, testUserAccountID)
	assert.Nil(t, err)

	err = SelectRow(rows, &ffTest)
	assert.Nil(t, err)

	assert.NotEqual(t, ffTest.FileName, "")
	assert.NotEqual(t, ffTest.FileID, "")
	assert.NotEqual(t, ffTest.FileSize, 0)
}

func TestSelectRowsSlice(t *testing.T) {
	dbConfig := newTestDBConfig()

	db, err := NewDatabase(dbConfig)
	assert.Nil(t, err)

	query := fmt.Sprintf(
		"SELECT %s, %s, %s, %s, %s FROM %s WHERE %s = ?",
		file.ColumnFileName,
		file.ColumnFileID,
		file.ColumnFileSize,
		file.ColumnModifiedOn,
		file.ColumnDeletedOn,
		file.TableName,
		file.ColumnFileOwnerID,
	)

	type FileFilterTest struct {
		FileName   string
		FileID     string
		FileSize   int
		ModifiedOn time.Time
		DeletedOn  *time.Time
	}

	ffTest := []FileFilterTest{}

	rows, err := db.Query(query, testUserAccountID)
	assert.Nil(t, err)

	err = SelectRows(rows, &ffTest)
	assert.Nil(t, err)

	fmt.Println(ffTest)
}

func TestMultipleSelectRows(t *testing.T) {
	fdb, err := getTestFileGateway()
	assert.Nil(t, err)

	root := t.TempDir()

	_, err = tests.CreateFiles(root)
	assert.Nil(t, err)

	fileIDs := []any{}

	files, err := file.Read(root)
	assert.Nil(t, err)

	for i, file := range files {
		files[i].OwnerID = testUserAccountID

		fileIDs = append(fileIDs, file.FileID)
	}

	err = fdb.AddFile(files)
	assert.Nil(t, err)

	type MultipleFileColumns struct {
		FileName string
		FileID   string
	}

	query := fmt.Sprintf(
		"SELECT %s,%s FROM %s",
		file.ColumnFileName,
		file.ColumnFileID,
		file.TableName,
	)

	cb := NewClauseBuilder()
	cb.In(file.ColumnFileID, fileIDs...)

	cbQ, args, err := cb.Build()
	assert.Nil(t, err)

	query = query + " " + cbQ

	rows, err := fdb.database.Query(query, args...)
	assert.Nil(t, err)

	_, err = DropRows(fdb.database, file.TableName, file.ColumnFileID, fileIDs...)
	assert.Nil(t, err)

	data := []MultipleFileColumns{}
	err = SelectRows(rows, &data)
	assert.Nil(t, err)

	assert.Equal(t, len(data), len(files))
}

func TestFailSelectRowsNilRows(t *testing.T) {
	v := []file.File{}

	err := SelectRows(nil, &v)
	assert.NotNil(t, err)
}

func TestFailSelectRowsNonPointer(t *testing.T) {
	v := []user.UserAccount{}

	udb, err := newTestUserGateway()
	assert.Nil(t, err)

	query := fmt.Sprintf("SELECT * FROM %s", user.TableName)

	rows, err := udb.database.Query(query)
	assert.Nil(t, err)

	err = SelectRows(rows, v)
	assert.NotNil(t, err)
}

func TestFailSelectRowsInvalidSize(t *testing.T) {
	v := []user.UserAccount{}

	udb, err := newTestUserGateway()
	assert.Nil(t, err)

	query := fmt.Sprintf("SELECT %s FROM %s", user.ColumnAccountID, user.TableName)

	rows, err := udb.database.Query(query)
	assert.Nil(t, err)

	err = SelectRows(rows, &v)
	assert.NotNil(t, err)
}
