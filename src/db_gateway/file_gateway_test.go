package dbgateway

import (
	"testing"
	"time"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/file"
	"github.com/bobllor/cloud-project/src/tests"
)

// IMPORTANT: These tests require the test database to exist.
// See the Docker setup documentation.
//
// By default when the test Docker setup is ran, one row is added to both
// the UserAccount and File table by default.
// As rows are added into the table, it will grow over the course of the test cases.
// Be aware of it!

// Variable of the default column data inserted into the test database on initialization.
const (
	testUserAccountID = "89672a64-f3ff-490c-8f2d-7e5cf5d4aa70"
	testFileID        = "randomfileidhere"
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

	fileFilters := []FileFilter{
		{
			Column:   file.FileIDCol,
			Args:     []any{testFileID},
			Type:     "IN",
			Operator: "AND",
		},
	}

	qFiles, err := fDb.GetFiles(testUserAccountID, fileFilters)
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

	// File.OwnerID is nil, this is changed to the existing account ID by default.
	for i := range files {
		files[i].OwnerID = testUserAccountID
	}

	err = fDb.AddFile(files)
	assert.Nil(t, err)

	qFiles, err := fDb.GetAllFiles(testUserAccountID)
	assert.Nil(t, err)

	// only 1 row exists by default, afterwards it adds however many from files
	assert.NotEqual(t, len(qFiles), 1)
}

func TestDeleteFiles(t *testing.T) {
	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	err = fDb.DeleteFiles(testUserAccountID, []string{testFileID})
	assert.Nil(t, err)

	fileFilters := []FileFilter{
		{
			Column:   file.FileIDCol,
			Args:     []any{testFileID},
			Type:     "IN",
			Operator: "AND",
		},
	}

	qFiles, err := fDb.GetFiles(testUserAccountID, fileFilters)
	assert.Nil(t, err)

	assert.NotNil(t, qFiles[0].DeletedOn)
	assert.Equal(t, qFiles[0].FileID, testFileID)

	now := time.Now()
	qDate := qFiles[0].DeletedOn

	assert.Equal(t, qDate.Year(), now.Year())
	assert.Equal(t, qDate.Month(), now.Month())
	assert.Equal(t, qDate.Day(), now.Day())
}

func TestRestoreFiles(t *testing.T) {
	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	fileFilters := []FileFilter{
		{
			Column:   file.FileIDCol,
			Args:     []any{testFileID},
			Type:     "EQUAL",
			Operator: "AND",
		},
	}

	err = fDb.DeleteFiles(testUserAccountID, []string{testFileID})
	assert.Nil(t, err)

	qFiles, err := fDb.GetFiles(testUserAccountID, fileFilters)
	assert.Nil(t, err)

	if qFiles[0].DeletedOn == nil {
		t.Fatal("failed to set file to deleted with DeletedOn")
	}

	err = fDb.RestoreFiles(testUserAccountID, []string{testFileID})
	assert.Nil(t, err)

	qFiles, err = fDb.GetFiles(testUserAccountID, fileFilters)
	assert.Nil(t, err)

	// whoops my assert library fails this. TODO: need to fix!
	if qFiles[0].DeletedOn != nil {
		t.Fatal("failed restoring deletion to file on column DeletedOn")
	}
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

// getFileDb gets the [FileGateway] for the test database.
// If an error occurs, it will return an error.
//
// This function does not start the test database instance.
func getTestFileGateway() (*FileGateway, error) {
	port := "3307"

	user := "root"
	password := ""
	net := "tcp"
	addr := "127.0.0.1" + ":" + port
	dbName := "TestLocalCloudStorage"

	config := NewConfig(user, password, net, addr, dbName)

	db, err := NewDatabase(config)
	if err != nil {
		return nil, err
	}

	fDb := NewFileGateway(db)

	return fDb, nil
}
