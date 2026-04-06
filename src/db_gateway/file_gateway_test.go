package dbgateway

import (
	"io"
	"log"
	"testing"
	"time"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/config"
	"github.com/bobllor/cloud-project/src/file"
	"github.com/bobllor/cloud-project/src/tests"
	"github.com/bobllor/cloud-project/src/utils"
	"github.com/bobllor/gologger"
)

// IMPORTANT: These tests require the test database to exist.
// See the Docker setup documentation.
//
// By default when the test Docker setup is ran, one row is added to both
// the UserAccount and File table by default.
// As rows are added into the table, it will grow over the course of the test cases.
// Majority of the tests with modifications only affects the default rows.
// Be aware of it!

// Constant variables that are the column data of the first (and by default) entries in the test DB.
const (
	testUserAccountID   = "89672a64-f3ff-490c-8f2d-7e5cf5d4aa70"
	testFileID          = "randomfileidhere"
	testDefaultFileName = "test1.txt"
)

func TestGetAllFiles(t *testing.T) {
	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	files, err := fDb.GetAllFiles(testUserAccountID)
	assert.Nil(t, err)

	assert.NotEqual(t, len(files), 0)
}

func TestGetFile(t *testing.T) {
	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	conditions := []WhereCondition{
		{
			Column:             file.ColumnFileID,
			Args:               []any{testFileID},
			LogicalOperator:    OperatorAnd,
			ComparisonOperator: Equal,
		},
	}

	qFiles, err := fDb.GetFiles(testUserAccountID, conditions)
	assert.Nil(t, err)

	assert.Equal(t, len(qFiles), 1)
}

func TestAddFile(t *testing.T) {
	dir := t.TempDir()

	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	_, err = tests.CreateFiles(dir)
	assert.Nil(t, err)

	files, err := file.Read(dir)
	assert.Nil(t, err)

	fileIDs := []string{}
	// File.OwnerID is nil, this is changed to the existing account ID by default.
	for i := range files {
		files[i].OwnerID = testUserAccountID

		fileIDs = append(fileIDs, files[i].FileID)
	}

	err = fDb.AddFile(files)
	assert.Nil(t, err)

	qFiles, err := fDb.GetAllFiles(testUserAccountID)
	assert.Nil(t, err)

	// only 1 row exists by default, afterwards it adds however many from files
	assert.NotEqual(t, len(qFiles), 1)
	assert.NotEqual(t, len(qFiles), 0)

	_, err = DropRows(fDb.database, file.TableName, file.ColumnFileID, utils.ConvertToAny(fileIDs)...)
	assert.Nil(t, err)
}

func TestUpdateFileByID(t *testing.T) {
	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	newValue := "this.is.a.text.file.txt"

	cd := ClauseData{
		Columns: []string{file.ColumnFileName},
		Args:    []any{newValue},
	}

	err = fDb.UpdateFileByID(testUserAccountID, testFileID, cd)
	assert.Nil(t, err)

	files, err := fDb.GetFiles(testUserAccountID, getConditionByID(testFileID))
	assert.Nil(t, err)

	assert.Equal(t, files[0].Name, newValue)

	err = resetDefaultFileRow(fDb, file.ColumnFileName, testDefaultFileName)
	assert.Nil(t, err)
}

func TestDeleteFiles(t *testing.T) {
	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	err = fDb.DeleteFiles(testUserAccountID, []string{testFileID})
	assert.Nil(t, err)

	conditions := []WhereCondition{
		{
			Column:             file.ColumnFileID,
			Args:               []any{testFileID},
			LogicalOperator:    OperatorAnd,
			ComparisonOperator: Equal,
		},
	}

	qFiles, err := fDb.GetFiles(testUserAccountID, conditions)
	assert.Nil(t, err)

	assert.NotNil(t, qFiles[0].DeletedOn)
	assert.Equal(t, qFiles[0].FileID, testFileID)

	now := time.Now()
	qDate := qFiles[0].DeletedOn

	assert.Equal(t, qDate.Year(), now.Year())
	assert.Equal(t, qDate.Month(), now.Month())
	assert.Equal(t, qDate.Day(), now.Day())

	err = resetDefaultFileRow(fDb, file.ColumnDeletedOn, nil)
	assert.Nil(t, err)
}

func TestRestoreFiles(t *testing.T) {
	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	conditions := []WhereCondition{
		{
			Column:             file.ColumnFileID,
			Args:               []any{testFileID},
			LogicalOperator:    OperatorAnd,
			ComparisonOperator: Equal,
		},
	}

	err = fDb.DeleteFiles(testUserAccountID, []string{testFileID})
	assert.Nil(t, err)

	qFiles, err := fDb.GetFiles(testUserAccountID, conditions)
	assert.Nil(t, err)

	if qFiles[0].DeletedOn == nil {
		t.Fatal("failed to set file to deleted with DeletedOn")
	}

	err = fDb.RestoreFiles(testUserAccountID, []string{testFileID})
	assert.Nil(t, err)

	qFiles, err = fDb.GetFiles(testUserAccountID, conditions)
	assert.Nil(t, err)

	// whoops my assert library fails this. TODO: need to fix!
	if qFiles[0].DeletedOn != nil {
		t.Fatal("failed restoring deletion to file on column DeletedOn")
	}
}

func TestUpdateModifiedFile(t *testing.T) {
	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	conditions := []WhereCondition{
		{
			Column:             file.ColumnFileID,
			Args:               []any{testFileID},
			LogicalOperator:    OperatorAnd,
			ComparisonOperator: Equal,
		},
	}

	baseFiles, err := fDb.GetFiles(testUserAccountID, conditions)
	assert.Nil(t, err)

	baseDate := baseFiles[0].ModifiedOn

	err = fDb.UpdateModifiedFiles(testUserAccountID, []string{testFileID})
	assert.Nil(t, err)

	newFiles, err := fDb.GetFiles(testUserAccountID, conditions)
	assert.Nil(t, err)

	newDate := newFiles[0].ModifiedOn

	assert.Equal(t, baseDate.Compare(newDate), -1)

	err = resetDefaultFileRow(fDb, file.ColumnModifiedOn, baseDate)
	assert.Nil(t, err)
}

func TestAddDuplicateFileError(t *testing.T) {
	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	f := file.File{
		OwnerID: testUserAccountID,
		FileID:  testFileID,
	}

	err = fDb.AddFile([]file.File{f})
	assert.NotNil(t, err)
}

func TestAddMissingOwnerIDFileError(t *testing.T) {
	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	f := file.File{
		FileID: "fdsa",
	}

	err = fDb.AddFile([]file.File{f})
	assert.NotNil(t, err)
}

func TestUpdateFiles(t *testing.T) {
	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	newName := "this.isa.filename.txt"

	cd := ClauseData{
		Columns: []string{file.ColumnFileName},
		Args:    []any{newName},
	}

	conditions := []WhereCondition{
		{
			Column:             file.ColumnFileID,
			Args:               []any{testFileID},
			ComparisonOperator: Equal,
			LogicalOperator:    OperatorAnd,
		},
	}

	files, err := fDb.GetFiles(testUserAccountID, conditions)
	assert.Nil(t, err)

	baseName := files[0].Name

	err = fDb.UpdateFiles(testUserAccountID, cd, conditions)
	assert.Nil(t, err)

	files, err = fDb.GetFiles(testUserAccountID, conditions)
	assert.Nil(t, err)

	assert.Equal(t, files[0].Name, newName)
	assert.NotEqual(t, files[0].Name, baseName)

	err = resetDefaultFileRow(fDb, file.ColumnFileName, testDefaultFileName)
	assert.Nil(t, err)
}

// getFileDb gets the [FileGateway] for the test database.
// If an error occurs, it will return an error.
//
// This function does not start the test database instance.
func getTestFileGateway() (*FileGateway, error) {
	dbConfig := newTestDBConfig()
	db, err := NewDatabase(dbConfig)
	if err != nil {
		return nil, err
	}

	logger := gologger.NewLogger(log.New(io.Discard, "", log.Ldate|log.Ltime), gologger.Lsilent)
	stdConfig := config.NewConfig(logger)

	fDb := NewFileGateway(db, stdConfig)

	return fDb, nil
}

func getConditionByID(fileID string) []WhereCondition {
	return []WhereCondition{
		{
			Column:             file.ColumnFileID,
			Args:               []any{fileID},
			ComparisonOperator: Equal,
			LogicalOperator:    OperatorAnd,
		},
	}
}

// getClauseData retrieves a default ClauseData that targets
// the a given column and given arguments.
func getClauseData(column string, args ...any) ClauseData {
	return ClauseData{
		Columns: []string{column},
		Args:    args,
	}
}

// resetDefaultFileName resets the default entry's file name to its default value.
// The error must be handled.
func resetDefaultFileRow(fDb *FileGateway, column string, args ...any) error {
	cd := getClauseData(column, args...)

	err := fDb.UpdateFileByID(testUserAccountID, testFileID, cd)
	if err != nil {
		return err
	}
	return nil
}
