package dbgateway

import (
	"testing"

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
//
// The row in UserAccount uses the ID: 89672a64-f3ff-490c-8f2d-7e5cf5d4aa70

var userAccountID = "89672a64-f3ff-490c-8f2d-7e5cf5d4aa70"

func TestFileQuery(t *testing.T) {
	fDb, err := getTestFileGateway()
	assert.Nil(t, err)

	files, err := fDb.QueryFile(userAccountID)
	assert.Nil(t, err)

	assert.NotEqual(t, len(files), 0)
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
		files[i].OwnerID = userAccountID
	}

	err = fDb.AddFile(files)
	assert.Nil(t, err)
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
