package dbe

import (
	"database/sql"

	"github.com/sedind/flow/app/randx"

	"github.com/pkg/errors"
)

// Connection represents all of the necessary details for
// talking with a datastore
type Connection struct {
	ID      string
	Store   *sql.DB
	Dialect Dialect
}

// NewConnection creates a new connection, and sets it's `Dialect`
// appropriately based on the `ConnectionDetails` passed into it.
func NewConnection(details Details) (*Connection, error) {
	err := details.Finalize()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := &Connection{
		ID: randx.String(32),
	}

	dialect, err := newDialect(details)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.Dialect = dialect

	return c, nil
}

// Open creates a new datasource connection
func (c *Connection) Open() error {
	if c.Store != nil {
		return nil
	}
	db, err := sql.Open(c.Dialect.Name(), c.Dialect.URL())
	if err != nil {
		return errors.WithStack(err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return errors.WithStack(err)
	}

	db.SetMaxOpenConns(c.Dialect.Details().Pool)
	db.SetMaxIdleConns(c.Dialect.Details().IdlePool)

	c.Store = db
	return nil
}

// Close destroys an active datasource connection
func (c *Connection) Close() error {
	return errors.Wrap(c.Store.Close(), "could not close connection")
}
