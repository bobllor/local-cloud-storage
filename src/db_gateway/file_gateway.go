package dbgateway

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bobllor/cloud-project/src/file"
)

const (
	FileTable string = "File"
)

type FileGateway struct {
	database       *sql.DB
	fileFieldCount int
}

func NewFileGateway(database *sql.DB) *FileGateway {
	f := &FileGateway{
		database:       database,
		fileFieldCount: file.FileColumnSize,
	}

	return f
}

// QueryFile queries the File table and returns a File slice.
// If an error occurs then it will return an error, and abort
// the scanning process if it is occurring.
func (f *FileGateway) QueryFile(fileOwnerID string) ([]file.File, error) {
	files := []file.File{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", FileTable, file.FileOwnerIDCol)
	q, err := f.database.Query(query, fileOwnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %v", FileTable, err)
	}

	defer q.Close()
	for q.Next() {
		f := file.File{}
		modifiedTimeSl := make([]uint8, 0)
		deletedTimeSl := make([]uint8, 0)
		scanErr := q.Scan(
			&f.OwnerID,
			&f.Name,
			&f.Type,
			&f.FileID,
			&f.ParentID,
			&f.Path,
			&f.Size,
			&modifiedTimeSl,
			&deletedTimeSl,
		)

		dateFormat := time.DateTime
		modifiedDate, err := time.Parse(dateFormat, string(modifiedTimeSl))
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s column: %v", file.ModifiedDateCol, err)
		}

		var deletedDate *time.Time
		if len(deletedTimeSl) > 0 {
			dateTemp, err := time.Parse(dateFormat, string(deletedTimeSl))
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s column: %v", file.DeletedDateCol, err)
			}

			deletedDate = &dateTemp
		}

		f.ModifiedTime = modifiedDate
		f.DeletedOn = deletedDate

		if scanErr != nil {
			return nil, scanErr
		}

		files = append(files, f)
	}

	return files, nil
}

// AddFile adds slice of [File] structs to the File database.
// If an error occurs it will return an error.
//
// This does not write the files to the disk.
func (f *FileGateway) AddFile(files []file.File) error {
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
