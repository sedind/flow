package dbe

import (
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Tx - transaction struct
type Tx struct {
	ID int
	*sqlx.Tx
}

// NewTx creates new Transaction
func NewTx(db *DB) (*Tx, error) {
	t := &Tx{
		ID: rand.Int(),
	}
	tx, err := db.Beginx()
	t.Tx = tx
	return t, errors.Wrap(err, "could not create new Transaction")
}

// Transaction simply returns the current transaction,
// this is defined so it implements the `Store` interface.
func (tx *Tx) Transaction() (*Tx, error) {
	return tx, nil
}

// Close -
func (tx *Tx) Close() error {
	return nil
}
