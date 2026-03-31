package dbgateway

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bobllor/cloud-project/src/file"
	"github.com/bobllor/cloud-project/src/utils"
)

func NewFileGateway(database *sql.DB) *FileGateway {
	f := &FileGateway{
		database:       database,
		fileFieldCount: file.FileColumnSize,
	}

	return f
}

type FileGateway struct {
	database       *sql.DB
	fileFieldCount int
}

// GetAllFiles returns a File slice of all File rows belonging to the file owner.
//
// If an error occurs then it will return an error, and abort
// the scanning process if it is occurring.
func (f *FileGateway) GetAllFiles(fileOwnerID string) ([]file.File, error) {
	cb := NewClauseBuilder()
	baseQuery := fmt.Sprintf("SELECT * FROM %s", file.FileTableName)
	cb.Equal(file.FileOwnerIDCol, fileOwnerID)

	con, args, err := cb.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build conditions: %v", err)
	}

	query := baseQuery + " " + con

	rows, err := f.database.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %v", file.FileTableName, err)
	}

	files, err := f.getFiles(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows: %v", err)
	}

	return files, nil
}

// GetFiles returns File rows based on the clause and conditions.
//
// batcher is a type that holds information for dynamic queries.
func (f *FileGateway) GetFiles(batcher *Batcher) ([]file.File, error) {
	cb := NewClauseBuilder()

	baseQuery := fmt.Sprintf("SELECT * FROM %s", file.FileTableName)

	err := cb.RegisterBatcher(batcher)
	if err != nil {
		return nil, err
	}

	q, args, err := cb.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	query := baseQuery + " " + q

	// TODO: add logging here
	fmt.Println(query)
	fmt.Println(args)

	rows, err := f.database.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}

	files, err := f.getFiles(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows: %v", err)
	}

	return files, nil
}

// AddFile adds slice of File structs to the File database.
// If an error occurs it will return an error.
//
// This does not write the files to the disk.
func (f *FileGateway) AddFile(files []file.File) error {
	if len(files) == 0 {
		return fmt.Errorf("no arguments given for AddFile")
	}

	flatFiles := file.FlattenFile(files...)

	query := fmt.Sprintf("INSERT INTO %s VALUES", file.FileTableName)
	paramStr := BuildPlaceholder(f.fileFieldCount, len(files))

	res, err := execQuery(f.database, query+" "+paramStr, flatFiles...)
	if err != nil {
		return fmt.Errorf("failed to insert into %s: %v", file.FileTableName, err)
	}

	rowCount, err := res.RowsAffected()
	// TODO: add logging here
	fmt.Printf("Rows inserted: %v", rowCount)

	return nil
}

// UpdateModifiedFile updates the modified date column to the current time.
func (f *FileGateway) UpdateModifiedFile(fileOwnerID string, fileIDs []string) error {
	cb := NewClauseBuilder()

	convIds := utils.ConvertToAny(fileIDs)

	cb.Equal(file.FileOwnerIDCol, fileOwnerID).And().In(file.FileIDCol, convIds...)

	qCon, args, err := cb.Build()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	query := fmt.Sprintf("UPDATE %s SET %s = ?", file.FileTableName, file.ModifiedOnCol) + " " + qCon

	finalArgs := []any{time.Now().Format(time.DateTime)}
	finalArgs = append(finalArgs, args...)

	res, err := execQuery(f.database, query, finalArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	// TODO: logging
	fmt.Println(res.RowsAffected())

	return nil
}

// DeleteFile sets a slice of file IDs to be marked for deletion.
// This does not delete the files immediately.
func (f *FileGateway) DeleteFiles(fileOwnerID string, fileIDs []string) error {
	cb := NewClauseBuilder()

	if len(fileIDs) == 0 {
		return fmt.Errorf("failed to delete files, got empty file IDs")
	}

	convFileIDs := utils.ConvertToAny(fileIDs)

	cb.Equal(file.FileOwnerIDCol, fileOwnerID).And().In(file.FileIDCol, convFileIDs...)

	qCondition, args, err := cb.Build()
	if err != nil {
		return fmt.Errorf("failed to build condition query: %v", err)
	}

	baseQuery := fmt.Sprintf(
		"UPDATE %s SET %s = ?",
		file.FileTableName,
		file.DeletedOnCol,
	)
	query := baseQuery + " " + qCondition

	// TODO: add logging
	fmt.Println(query)

	finalArgs := []any{time.Now().Format(time.DateTime)}
	finalArgs = append(finalArgs, args...)

	res, err := execQuery(f.database, query, finalArgs...)
	if err != nil {
		return err
	}

	fmt.Println(res.RowsAffected())

	return nil
}

// RestoreFiles sets a file IDs that are unmark files that were marked for deletion.
func (f *FileGateway) RestoreFiles(fileOwnerID string, fileIDs []string) error {
	cb := NewClauseBuilder()

	cb.Equal(file.FileOwnerIDCol, fileOwnerID)

	convIDs := utils.ConvertToAny(fileIDs)

	cb.And().In(file.FileIDCol, convIDs...)

	cond, args, err := cb.Build()
	if err != nil {
		return fmt.Errorf("failed to build conditions: %v", err)
	}

	baseQuery := fmt.Sprintf("UPDATE %s SET %s = NULL", file.FileTableName, file.DeletedOnCol)
	query := baseQuery + " " + cond

	rows, err := execQuery(f.database, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	// TODO: add logging
	fmt.Println(rows.RowsAffected())

	return nil
}

// getFiles is a helper function used to scan and return
// a slice of Files.
//
// sql.Rows will automatically be closed at the end of function.
func (f *FileGateway) getFiles(rows *sql.Rows) ([]file.File, error) {
	files := []file.File{}

	for rows.Next() {
		f := file.File{}
		// datetime is returned as a []uint38, uint38 -> string -> date
		modifiedTimeSl := make([]uint8, 0)
		deletedTimeSl := make([]uint8, 0)

		scanErr := rows.Scan(
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
			return nil, fmt.Errorf("failed to parse %s column: %v", file.ModifiedOnCol, err)
		}

		var deletedDate *time.Time
		// NULL is an empty slice
		if len(deletedTimeSl) > 0 {
			dateTemp, err := time.Parse(dateFormat, string(deletedTimeSl))
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s column: %v", file.DeletedOnCol, err)
			}

			deletedDate = &dateTemp
		}

		f.ModifiedOn = modifiedDate
		f.DeletedOn = deletedDate

		if scanErr != nil {
			return nil, scanErr
		}

		files = append(files, f)
	}

	return files, nil
}
