package dbcon

import (
	"database/sql"
	"fmt"

	"github.com/bobllor/cloud-project/src/file"
)

const (
	FilesTable string = "Files"
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
	q, err := f.database.Query(`SELECT 
		FileName, FileType, FileID, FilePath, FileSize 
		FROM Files`,
	)
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
		return fmt.Errorf("failed to start transaction on %s: %v", FilesTable, err)
	}
	defer tx.Rollback()

	flatFiles := file.FlattenFile(files...)

	query := fmt.Sprintf("INSERT INTO %s VALUES", FilesTable)
	paramStr := QueryParamBuilder(f.fileFieldCount, len(files))

	res, err := tx.Exec(query+" "+paramStr, flatFiles...)
	if err != nil {
		return fmt.Errorf("failed to insert into %s: %v", FilesTable, err)
	}

	rowCount, _ := res.RowsAffected()
	// TODO: add logging here
	fmt.Printf("Rows inserted: %v", rowCount)

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction for %s: %v", FilesTable, err)
	}

	return nil
}
