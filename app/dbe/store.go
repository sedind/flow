package dbe

import "database/sql"

//Store - represents DB store interface
type Store interface {
	Close() error
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}
