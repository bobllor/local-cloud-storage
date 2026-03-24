package dbcon

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

// NewDatabase opens the SQL database and returns a [sql.DB].
// It will return an error if any errors occur.
//
// The database is pinged in this call, if it fails then
// an error will be returned.
func NewDatabase(config mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewConfig creates a new [mysql.Config] from the arguments.
func NewConfig(user string, passwd string, net string, addr string, dbName string) mysql.Config {
	c := mysql.Config{}

	c.User = user
	c.Passwd = passwd
	c.Net = net
	c.Addr = addr
	c.DBName = dbName
	c.AllowNativePasswords = true

	return c
}
