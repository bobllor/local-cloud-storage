package file

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type File struct {
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

	// Ext is the file extension. This includes the final dot of
	// the extension. It will be an empty string if no extension is given.
	Ext string

	// Path is the absolute path to the file on the disk. This is intended
	// for the backend use only.
	Path string

	// Size is the size of the file.
	Size int64

	// Modified is the most recent time the file has been modified. This
	// will be the most recent time of change or when it was first created.
	Modified time.Time

	// DeleteTime is the time when the file is set to deleted. The acutal
	// deletion occurs after a certain amount of time has passed
	// since the marked deletion time. This value can be nil.
	DeleteTime *time.Time
}

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
				Name:     info.Name(),
				Type:     fileType,
				Size:     info.Size(),
				Ext:      path.Ext(p),
				Path:     p,
				FileID:   id,
				ParentID: parentID,
				Modified: info.ModTime(),
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
