package dbgateway

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bobllor/cloud-project/src/config"
	"github.com/bobllor/cloud-project/src/file"
	"github.com/bobllor/cloud-project/src/utils"
)

const ()

// NewFileGateway creates a new FileGateway for database related options.
func NewFileGateway(database *sql.DB, config *config.Config) *FileGateway {
	f := &FileGateway{
		database:       database,
		fileFieldCount: file.FileColumnSize,
		config:         config,
		logUtil:        LogUtility{log: config.Log},
	}

	return f
}

type FileGateway struct {
	database       *sql.DB
	fileFieldCount int
	config         *config.Config
	logUtil        LogUtility
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
func (f *FileGateway) GetFiles(fileOwnerID string, conditions []WhereCondition) ([]file.File, error) {
	cb := NewClauseBuilder()

	baseQuery := fmt.Sprintf("SELECT * FROM %s", file.FileTableName)

	cb.Equal(file.FileOwnerIDCol, fileOwnerID)

	err := cb.RegisterConditions(conditions)
	if err != nil {
		return nil, fmt.Errorf("failed to register conditions: %v", err)
	}

	q, args, err := cb.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	query := baseQuery + " " + q

	f.logUtil.LogQueryAndArgs(query, args)

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

// UpdateFileByID updates a single file row by its file ID.
//
// cd is a ClauseData type used to target the columns and values to replace for the row.
func (f *FileGateway) UpdateFileByID(fileOwnerID string, fileID string, cd ClauseData) error {
	conditions := []WhereCondition{
		{
			Column:             file.FileIDCol,
			Args:               []any{fileID},
			ComparisonOperator: Equal,
			LogicalOperator:    OperatorAnd,
		},
	}

	err := f.UpdateFiles(fileOwnerID, cd, conditions)
	if err != nil {
		return err
	}

	return nil
}

// UpdateFiles updates the Files table based ClauseData and conditions.
//
// Errors will be returned if one occurs.
// Certain columns are forbidden from being changed, and will return an error
// if these are found.
func (f *FileGateway) UpdateFiles(fileOwnerID string, cd ClauseData, conditions []WhereCondition) error {
	cb := NewClauseBuilder()

	setQ, err := cd.BuildSetQuery()
	if err != nil {
		return fmt.Errorf("failed to validate ClauseData: %v", err)
	}

	baseQuery := fmt.Sprintf("UPDATE %s", file.FileTableName) + " " + setQ

	cb.Equal(file.FileOwnerIDCol, fileOwnerID)

	err = cb.RegisterConditions(conditions)
	if err != nil {
		return fmt.Errorf("failed to register conditions: %v", err)
	}

	whereQ, args, err := cb.Build()
	if err != nil {
		return fmt.Errorf("failed to build WHERE clause: %v", err)
	}

	query := baseQuery + " " + whereQ

	execArgs := []any{}

	execArgs = append(execArgs, cd.Args...)
	execArgs = append(execArgs, args...)

	f.logUtil.LogQueryAndArgs(query, execArgs)

	res, err := execQuery(f.database, query, execArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v (args: %v)", err, execArgs)
	}

	f.logUtil.LogResultRows(res)

	return nil
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

	query = query + " " + paramStr

	f.logUtil.LogQueryAndArgs(query, flatFiles)

	res, err := execQuery(f.database, query, flatFiles...)
	if err != nil {
		return fmt.Errorf("failed to insert into %s: %v", file.FileTableName, err)
	}

	f.logUtil.LogResultRows(res)

	return nil
}

// UpdateModifiedFiles updates the modified date column to the current time.
func (f *FileGateway) UpdateModifiedFiles(fileOwnerID string, fileIDs []string) error {
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

	f.logUtil.LogQueryAndArgs(query, finalArgs)

	res, err := execQuery(f.database, query, finalArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	f.logUtil.LogResultRows(res)

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

	finalArgs := []any{time.Now().Format(time.DateTime)}
	finalArgs = append(finalArgs, args...)

	f.logUtil.LogQueryAndArgs(query, finalArgs)

	res, err := execQuery(f.database, query, finalArgs...)
	if err != nil {
		return err
	}

	f.logUtil.LogResultRows(res)

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

	f.logUtil.LogQueryAndArgs(query, args)

	res, err := execQuery(f.database, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	f.logUtil.LogResultRows(res)

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
