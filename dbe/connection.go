package dbe

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sedind/flow/dbe/dialect"
)

// Connection represents all of the necessary details for
// talking with a datastore
type Connection struct {
	ID      string
	Details Details
	Dialect dialect.Dialect
	Store   Store
	Tx      *Tx
}

// NewConnection creates a new connection, and sets it's `Dialect`
// appropriately based on the `Details` passed into it.
func NewConnection(details Details) (*Connection, error) {
	err := details.Finalize()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dialect, err := dialect.New(details.Dialect)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := &Connection{
		ID:      string(time.Now().Unix()),
		Details: details,
		Dialect: dialect,
	}

	return c, nil
}

// Open creates new datasource connection
func (c *Connection) Open() error {
	if c.Store != nil {
		return nil
	}

	dbc, err := sqlx.Open(c.Details.Dialect, c.Details.URL)
	if err != nil {
		return errors.WithStack(err)
	}
	dbc.SetMaxOpenConns(c.Details.Pool)
	dbc.SetMaxIdleConns(c.Details.IdlePool)

	c.Store = &db{dbc}

	return dbc.Ping()

}

// Close destroys an active datasource connection
func (c *Connection) Close() error {
	return errors.Wrap(c.Store.Close(), "could not close connection")
}

// NewTx starts a new transaction on the connection
func (c *Connection) NewTx() (*Connection, error) {
	if c.Tx == nil {
		tx, err := c.Store.Transaction()
		if err != nil {
			return c, errors.Wrap(err, "could not start new transaction")
		}
		cn := &Connection{
			ID:      string(time.Now().Unix()),
			Details: c.Details,
			Dialect: c.Dialect,
			Store:   tx,
			Tx:      tx,
		}
		return cn, nil
	}
	return c, nil
}

func (c *Connection) copy() *Connection {
	return &Connection{
		ID:      string(time.Now().Unix()),
		Details: c.Details,
		Dialect: c.Dialect,
		Store:   c.Store,
		Tx:      c.Tx,
	}
}
