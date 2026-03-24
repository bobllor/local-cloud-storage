package tests

import (
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/file"
)

func TestRead(t *testing.T) {
	dir := t.TempDir()

	_, err := createFiles(dir)
	assert.Nil(t, err)

	files, err := file.Read(dir)
	assert.Nil(t, err)
	assert.NotEqual(t, len(files), 0)

	for _, file := range files {
		_, err := os.Stat(file.Path)
		assert.Nil(t, err)
	}
}

func TestReadParentIdFolders(t *testing.T) {
	dir := t.TempDir()

	_, err := createFiles(dir)
	assert.Nil(t, err)

	files, err := file.Read(dir)
	assert.Nil(t, err)
	assert.NotEqual(t, len(files), 0)

	// a map of the parents [Dir(File.FilePath)]-[File] of
	// the files
	parentMap := map[string]file.File{}
	// flat map of [File.FileID]-[File] for all
	// the files
	fileIdMap := map[string]file.File{}

	for _, file := range files {
		_, err := os.Stat(file.Path)
		assert.Nil(t, err)

		_, ok := parentMap[file.Path]
		if !ok {
			parentMap[filepath.Dir(file.Path)] = file
		}

		_, ok = fileIdMap[file.FileID]
		if !ok {
			fileIdMap[file.FileID] = file
		}
	}

	for parentPath, file := range parentMap {
		if file.ParentID != nil {
			// if this fails, then the parent doesn't exist
			// in the flat file map
			idParentFile, ok := fileIdMap[*file.ParentID]
			assert.Equal(t, ok, true)

			// checks if the flat parent file is a directory
			assert.Equal(t, idParentFile.Type, "directory")
			// checks if the flat parent path is not root
			assert.NotEqual(t, idParentFile.Path, dir)
			// checks if the flat parent path is the same as the parentMap key
			assert.Equal(t, idParentFile.Path, parentPath)
		} else {
			// nil ID means the parent is root
			assert.Equal(t, parentPath, dir)
			assert.Equal(t, filepath.Dir(file.Path), dir)
		}
	}
}

func TestFailRead(t *testing.T) {
	_, err := file.Read("dir/does/not/exist")
	assert.NotNil(t, err)
}

// createFiles creates random files in the root. It will return
// a string slice containing the paths of the files created.
//
// root is the path where the files will be created in.
//
// An error will be returned if any errors occur during
// the file creation.
func createFiles(root string) ([]string, error) {
	paths := []string{}

	files := []string{
		"text1.txt",
		"text2.txt",
		"text3.txt",
		"some.logs.log",
		"database.db",
	}

	f1 := "folder1"
	f2 := "folder2"

	dirPaths := []string{
		root,
		filepath.Join(root, f1),
		filepath.Join(root, f2),
		filepath.Join(root, f1, f2),
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, file := range files {
		fullPath := filepath.Join(dirPaths[r.Intn(len(dirPaths)-1)], file)
		basePath := path.Dir(fullPath)

		err := os.MkdirAll(basePath, 0o744)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(fullPath, []byte{}, 0o644)
		if err != nil {
			return nil, err
		}

		paths = append(paths, fullPath)
	}

	return paths, nil
}
