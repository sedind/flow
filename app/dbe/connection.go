package dbe

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// Connection represents all of the necessary details for
// talking with a datastore
type Connection struct {
	ID      string
	Dialect Dialect
}

// NewConnection creates a new connection, and sets it's `Dialect`
// appropriately based on the `ConnectionDetails` passed into it.
func NewConnection(details *Details) (*Connection, error) {
	err := details.Finalize()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := &Connection{
		ID: string(time.Now().Unix()),
	}

	dialect, err := newDialect(details)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.Dialect = dialect

	dbc, err := sql.Open(c.Dialect.Name(), c.Dialect.URL())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := dbc.Ping(); err != nil {
		dbc.Close()
		return nil, errors.WithStack(err)
	}

	dbc.SetMaxOpenConns(c.Dialect.Details().Pool)
	dbc.SetMaxIdleConns(c.Dialect.Details().IdlePool)

	store := &DB{dbc}
	c.Dialect.SetStore(store)

	return c, nil
}

// Close destroys an active datasource connection
func (c *Connection) Close() error {
	return errors.Wrap(c.Dialect.Store().Close(), "could not close connection")
}
