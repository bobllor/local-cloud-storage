package dbcon

import "github.com/bobllor/cloud-project/src/file"

type FilesDB struct {
	inquirer DBInquirer
}

func NewFilesDatabase(inquirer DBInquirer) *FilesDB {
	f := &FilesDB{
		inquirer: inquirer,
	}

	return f
}

// QueryFiles queries the File table and returns a File slice.
// If an error occurs then it will return an error, and abort
// the scanning process if it is occurring.
func (f *FilesDB) QueryFiles() ([]file.File, error) {
	files := []file.File{}

	// TODO: add WHERE filter with user ID when added
	q, err := f.inquirer.Query(`SELECT 
		FileName, FileType, FileID, Extension, FilePath, FileSize 
		FROM Files`,
	)
	if err != nil {
		return nil, err
	}

	defer q.Close()
	for q.Next() {
		f := file.File{}
		scanErr := q.Scan(&f.Name, &f.Type, &f.FileID, &f.Ext, &f.Path, &f.Size)
		if scanErr != nil {
			return nil, scanErr
		}

		files = append(files, f)
	}

	return files, nil
}
