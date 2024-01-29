package database

import (
	"context"
	"database/sql"
)

// SqlGdbc (SQL Go database connection) is a wrapper for SQL database handler ( can be *sql.DB or *sql.Tx)
// It should be able to work with all SQL data that follows SQL standard.
type SqlGdbc interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	// Get and map single row to dest object
	Get(dest interface{}, query string, args ...interface{}) error
	// Get and map multiple rows to dest objects
	Select(dest interface{}, query string, args ...interface{}) error
	// If need transaction support, add this interface
	Transactor
}

type Transactor interface {
	WithinTransaction(context.Context, func(ctx context.Context) error) error
}
