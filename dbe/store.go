package dbe

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Store describes set of functionalities needed to interact with DB storage
type Store interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Select(interface{}, string, ...interface{}) error
	Get(interface{}, string, ...interface{}) error
	NamedExec(string, interface{}) (sql.Result, error)
	Exec(string, ...interface{}) (sql.Result, error)
	PrepareNamed(string) (*sqlx.NamedStmt, error)
	Transaction() (*Tx, error)
	Rollback() error
	Commit() error
	Close() error
}
