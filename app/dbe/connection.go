package dbe

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Connection represents all of the necessary details for
// talking with a datastore
type Connection struct {
	ID      string
	Details Details
	DB      *sqlx.DB
}

// NewConnection creates a new connection, and sets it's `Dialect`
// appropriately based on the `Details` passed into it.
func NewConnection(details Details) (*Connection, error) {
	err := details.Finalize()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := &Connection{
		ID:      string(time.Now().Unix()),
		Details: details,
	}

	return c, nil
}

// Open creates new datasource connection
func (c *Connection) Open() error {
	if c.DB != nil {
		return nil
	}

	dbc, err := sqlx.Open(c.Details.Dialect, c.Details.URL)
	if err != nil {
		return errors.WithStack(err)
	}
	dbc.SetMaxOpenConns(c.Details.Pool)
	dbc.SetMaxIdleConns(c.Details.IdlePool)

	c.DB = dbc
	return nil

}

// Close destroys an active datasource connection
func (c *Connection) Close() error {
	return errors.Wrap(c.DB.Close(), "could not close connection")
}
