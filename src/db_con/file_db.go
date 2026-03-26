package dbcon

import (
	"database/sql"
	"fmt"

	"github.com/bobllor/cloud-project/src/file"
)

const (
	FileTable string = "File"
)

type FilesDB struct {
	database       *sql.DB
	fileFieldCount int
}

func NewFilesDatabase(database *sql.DB) *FilesDB {
	f := &FilesDB{
		database:       database,
		fileFieldCount: file.FileColumnSize,
	}

	return f
}

// QueryFiles queries the File table and returns a File slice.
// If an error occurs then it will return an error, and abort
// the scanning process if it is occurring.
func (f *FilesDB) QueryFiles() ([]file.File, error) {
	files := []file.File{}

	// TODO: add WHERE filter with user ID when added
	query := fmt.Sprintf("SELECT FileName, FileType, FileID, FilePath, FileSize FROM %s", FileTable)
	q, err := f.database.Query(query)
	if err != nil {
		return nil, err
	}

	defer q.Close()
	for q.Next() {
		f := file.File{}
		scanErr := q.Scan(&f.Name, &f.Type, &f.FileID, &f.Path, &f.Size)
		if scanErr != nil {
			return nil, scanErr
		}

		files = append(files, f)
	}

	return files, nil
}

// AddFile adds a File to the Files database.
// If an error occurs it will return an error.
//
// This does not write the files to the disk.
func (f *FilesDB) AddFile(files []file.File) error {
	if len(files) == 0 {
		return fmt.Errorf("no arguments given for AddFile")
	}
	tx, err := f.database.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction on %s: %v", FileTable, err)
	}
	defer tx.Rollback()

	flatFiles := file.FlattenFile(files...)

	query := fmt.Sprintf("INSERT INTO %s VALUES", FileTable)
	paramStr := QueryParamBuilder(f.fileFieldCount, len(files))

	res, err := tx.Exec(query+" "+paramStr, flatFiles...)
	if err != nil {
		return fmt.Errorf("failed to insert into %s: %v", FileTable, err)
	}

	rowCount, _ := res.RowsAffected()
	// TODO: add logging here
	fmt.Printf("Rows inserted: %v", rowCount)

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction for %s: %v", FileTable, err)
	}

	return nil
}
