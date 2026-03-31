package file

import (
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const (
	// FileColumnSize is the amount of columns used for the Files table.
	// It is equal to the public fields of the [File] struct.
	FileColumnSize int    = 9
	FileTableName  string = "File"
	FileOwnerIDCol string = "FileOwnerID"
	FileNameCol    string = "FileName"
	FileTypeCol    string = "FileType"
	FileIDCol      string = "FileID"
	ParentIDCol    string = "ParentID"
	FilePathCol    string = "FilePath"
	FileSizeCol    string = "FileSize"
	ModifiedOnCol  string = "ModifiedOn"
	DeletedOnCol   string = "DeletedOn"
)

// Read returns a File slice for all files found in root.
// An error will be returned if there is an issue while reading root.
//
// This is only intended for local disk access, and is not intended to be
// used for adding files into the database outside of some situations.
// Adding into the database is delegated to the API call.
func Read(root string) ([]File, error) {
	fs, err := walk(root)
	if err != nil {
		return nil, err
	}

	return fs, nil
}

// walk is used to traverse root and return a File slice for
// all the files in root.
//
// If any error occurs during the file reading, it will abort the process
// and return an error.
func walk(root string) ([]File, error) {
	fs := []File{}

	folderIDMap := map[string]*string{
		root: nil,
	}

	// folder name of root is the account ID
	// this only is applicable to local files
	accountID := filepath.Base(root)

	walkFunc := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		id := uuid.NewString()

		// skipping root
		if p != root {
			fileType := "file"

			var parentID *string
			parent := filepath.Dir(p)
			if info.IsDir() {
				fileType = "directory"
				_, ok := folderIDMap[p]
				if !ok {
					folderIDMap[p] = &id
				}
			}

			pID, ok := folderIDMap[parent]
			if ok {
				parentID = pID
			}

			f := File{
				Name:       info.Name(),
				Type:       fileType,
				Size:       info.Size(),
				Path:       p,
				FileID:     id,
				ParentID:   parentID,
				ModifiedOn: info.ModTime(),
				OwnerID:    accountID,
			}

			fs = append(fs, f)
		}

		return nil
	}

	err := filepath.Walk(root, walkFunc)
	if err != nil {
		return nil, err
	}

	return fs, nil
}

// FlattenFiles flattens the a slice of File structs to prepare for use in
// a query.
func FlattenFile(files ...File) []any {
	out := []any{}

	appendFunc := func(v any) {
		out = append(out, v)
	}

	for _, file := range files {
		appendFunc(file.OwnerID)
		appendFunc(file.Name)
		appendFunc(file.Type)
		appendFunc(file.FileID)
		appendFunc(file.ParentID)
		appendFunc(file.Path)
		appendFunc(file.Size)
		appendFunc(file.ModifiedOn)
		appendFunc(file.DeletedOn)
	}

	return out
}

type File struct {
	// OwnerID is the ID of the owner of the file.
	OwnerID string

	// Name is the name of the file. This includes the extension
	// of the file.
	Name string

	// Type is the file type. This is either a "directory" or
	// a "file".
	Type string

	// FileID is a unique ID assigned to the file.
	FileID string

	// ParentID is the parent's unique ID that the file resides in.
	// This can be nil, meaning it resides in the root folder.
	ParentID *string

	// Path is the absolute path to the file on the disk. This is intended
	// for the backend use only.
	Path string

	// Size is the size of the file.
	Size int64

	// ModifiedOn is the most recent time the file has been modified. This
	// will be the most recent time of change or when it was first created.
	ModifiedOn time.Time

	// DeletedOn is the time when the file is set to be deleted. The acutal
	// deletion occurs after a certain amount of time has passed
	// since the marked deletion time. This value can be nil.
	DeletedOn *time.Time
}
