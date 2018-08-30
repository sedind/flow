package dbe

import "github.com/jmoiron/sqlx"

// DB - dbe DB object which wraps native sql.DB
type DB struct {
	*sqlx.DB
}

// Transaction creates new Transaction on DB connection
func (db *DB) Transaction() (*Tx, error) {
	return NewTx(db)
}

// Rollback - does nothing
func (db *DB) Rollback() error {
	return nil
}

// Commit - does nothing
func (db *DB) Commit() error {
	return nil
}
