package dbgateway

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bobllor/cloud-project/src/file"
	"github.com/bobllor/cloud-project/src/session"
	"github.com/bobllor/cloud-project/src/sqlquery"
	"github.com/bobllor/cloud-project/src/utils"
)

// NewFileGateway creates a new FileGateway for database related options.
func NewFileGateway(database *sql.DB, deps *utils.Deps) *FileGateway {
	f := &FileGateway{
		database:       database,
		fileFieldCount: file.ColumnSize,
		deps:           deps,
	}

	return f
}

type FileGateway struct {
	database       *sql.DB
	fileFieldCount int
	deps           *utils.Deps
}

// GetAllFiles returns a File slice of all File rows belonging to the file owner.
//
// If an error occurs then it will return an error, and abort
// the scanning process if it is occurring.
func (f *FileGateway) GetAllFiles(fileOwnerID string) ([]file.File, error) {
	cb := NewClauseBuilder()
	baseQuery := fmt.Sprintf("SELECT * FROM %s", file.TableName)
	cb.Equal(file.ColumnFileOwnerID, fileOwnerID)

	con, args, err := cb.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build conditions: %v", err)
	}

	query := baseQuery + " " + con

	rows, err := f.database.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %v", file.TableName, err)
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

	baseQuery := fmt.Sprintf("SELECT * FROM %s", file.TableName)

	cb.Equal(file.ColumnFileOwnerID, fileOwnerID)

	err := cb.RegisterConditions(conditions)
	if err != nil {
		return nil, fmt.Errorf("failed to register conditions: %v", err)
	}

	q, args, err := cb.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	query := baseQuery + " " + q

	f.deps.Log.Debugf("Query: %s", query)
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
			Column:             file.ColumnFileID,
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

	setQ, sargs, err := cd.BuildSetQuery()
	if err != nil {
		return fmt.Errorf("failed to validate ClauseData: %v", err)
	}

	baseQuery := fmt.Sprintf("UPDATE %s %s", file.TableName, setQ)

	cb.Equal(file.ColumnFileOwnerID, fileOwnerID)

	err = cb.RegisterConditions(conditions)
	if err != nil {
		return fmt.Errorf("failed to register conditions: %v", err)
	}

	whereQ, args, err := cb.Build()
	if err != nil {
		return fmt.Errorf("failed to build WHERE clause: %v", err)
	}

	query := baseQuery + " " + whereQ

	execArgs := MakeArgs(sargs, args)

	logQuery(f.deps.Log, query)

	f.deps.Log.Debugf("Query: %s", query)
	res, err := execQuery(f.database, query, execArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v (args: %v)", err, execArgs)
	}

	logResultRows(f.deps.Log, res)

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

	query := fmt.Sprintf("INSERT INTO %s VALUES", file.TableName)
	paramStr := BuildPlaceholder(f.fileFieldCount, len(files))

	query = query + " " + paramStr

	logQuery(f.deps.Log, query)
	res, err := execQuery(f.database, query, flatFiles...)
	if err != nil {
		return fmt.Errorf("failed to insert into %s: %v", file.TableName, err)
	}

	logResultRows(f.deps.Log, res)

	return nil
}

// UpdateModifiedFiles updates the modified date column to the current time.
func (f *FileGateway) UpdateModifiedFiles(fileOwnerID string, fileIDs []string) error {
	cb := NewClauseBuilder()

	convIds := utils.ConvertToAny(fileIDs)

	cb.Equal(file.ColumnFileOwnerID, fileOwnerID).And().In(file.ColumnFileID, convIds...)

	qCon, args, err := cb.Build()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	query := fmt.Sprintf("UPDATE %s SET %s = ?", file.TableName, file.ColumnModifiedOn) + " " + qCon

	finalArgs := []any{time.Now().Format(time.DateTime)}
	finalArgs = append(finalArgs, args...)

	logQuery(f.deps.Log, query)
	res, err := execQuery(f.database, query, finalArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	logResultRows(f.deps.Log, res)

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

	cb.Equal(file.ColumnFileOwnerID, fileOwnerID).And().In(file.ColumnFileID, convFileIDs...)

	qCondition, args, err := cb.Build()
	if err != nil {
		return fmt.Errorf("failed to build condition query: %v", err)
	}

	baseQuery := fmt.Sprintf(
		"UPDATE %s SET %s = ?",
		file.TableName,
		file.ColumnDeletedOn,
	)
	query := baseQuery + " " + qCondition

	finalArgs := []any{time.Now().Format(time.DateTime)}
	finalArgs = append(finalArgs, args...)

	logQuery(f.deps.Log, query)

	res, err := execQuery(f.database, query, finalArgs...)
	if err != nil {
		return err
	}

	logResultRows(f.deps.Log, res)

	return nil
}

// RestoreFiles sets a file IDs that are unmark files that were marked for deletion.
func (f *FileGateway) RestoreFiles(fileOwnerID string, fileIDs []string) error {
	cb := NewClauseBuilder()

	cb.Equal(file.ColumnFileOwnerID, fileOwnerID)

	convIDs := utils.ConvertToAny(fileIDs)

	cb.And().In(file.ColumnFileID, convIDs...)

	cond, args, err := cb.Build()
	if err != nil {
		return fmt.Errorf("failed to build conditions: %v", err)
	}

	baseQuery := fmt.Sprintf("UPDATE %s SET %s = NULL", file.TableName, file.ColumnDeletedOn)
	query := baseQuery + " " + cond

	logQuery(f.deps.Log, query)

	res, err := execQuery(f.database, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	logResultRows(f.deps.Log, res)

	return nil
}

// GetFilesBySessionAndParentFolder retrieves the files based on the session ID and the parent
// of the folder.
func (f *FileGateway) GetFilesBySessionAndParentFolder(sessionID string, parentFolderID *string) ([]file.File, error) {
	// joins are raw SQL, not going to make it into an ORM due to how complex it is
	// creates the basic main query for combination with join
	// the WHERE clause is appended later
	mainQuery, _, err := sqlquery.Select(fmt.Sprintf("%s f", file.TableName), "f.*").Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	args := []any{sessionID}
	parentCondition := "= ?"
	if parentFolderID == nil {
		parentCondition = "IS NULL"
	} else {
		args = append(args, parentFolderID)
	}

	query := fmt.Sprintf(
		"%s JOIN %s ON s.%s = f.%s WHERE s.%s = ? AND %s %s",
		mainQuery,
		fmt.Sprintf("%s s", session.TableName),
		session.ColumnAccountID,
		file.ColumnFileOwnerID,
		session.ColumnSessionID,
		file.ColumnParentID,
		parentCondition,
	)

	logQuery(f.deps.Log, query)
	rows, err := f.database.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %v", err)
	}

	files, err := f.getFiles(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to parse File query: %v", err)
	}

	return files, nil
}

// getFiles is a helper function used to scan and return
// a slice of Files.
//
// sql.Rows will automatically be closed at the end of function.
func (f *FileGateway) getFiles(rows *sql.Rows) ([]file.File, error) {
	files := []file.File{}

	for rows.Next() {
		f := file.File{}

		scanErr := rows.Scan(
			&f.OwnerID,
			&f.Name,
			&f.Type,
			&f.FileID,
			&f.ParentID,
			&f.Path,
			&f.Size,
			&f.ModifiedOn,
			&f.DeletedOn,
		)

		if scanErr != nil {
			return nil, scanErr
		}

		files = append(files, f)
	}

	return files, nil
}
