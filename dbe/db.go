package dbe

import "github.com/jmoiron/sqlx"

// db struct is sqlx.DB wrapper usedd to implement Store interface
type db struct {
	*sqlx.DB
}

func (db *db) Transaction() (*Tx, error) {
	return newTx(db)
}

func (db *db) Rollback() error {
	return nil
}

func (db *db) Commit() error {
	return nil
}
