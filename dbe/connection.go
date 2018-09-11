package dbe

import (
	"sync/atomic"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Connection represents all of the necessary details for
// talking with a datastore
type Connection struct {
	ID      string
	Details Details
	Store   *sqlx.DB
	TX      *sqlx.Tx
	Elapsed int64
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
	if c.Store != nil {
		return nil
	}

	dbc, err := sqlx.Open(c.Details.Dialect, c.Details.URL)
	if err != nil {
		return errors.WithStack(err)
	}
	dbc.SetMaxOpenConns(c.Details.Pool)
	dbc.SetMaxIdleConns(c.Details.IdlePool)

	c.Store = dbc

	return dbc.Ping()

}

// Close destroys an active datasource connection
func (c *Connection) Close() error {
	return errors.Wrap(c.Store.Close(), "could not close connection")
}

// NewTransaction starts a new transaction on the connection
func (c *Connection) NewTransaction() (*Connection, error) {
	if c.TX == nil {
		tx, err := c.Store.Beginx()
		if err != nil {
			return nil, errors.Wrap(err, "could not start new transaction")
		}
		cn := &Connection{
			ID:      string(time.Now().Unix()),
			Details: c.Details,
			Store:   c.Store,
			TX:      tx,
		}
		return cn, nil
	}
	return c, nil
}

func (c *Connection) copy() *Connection {
	return &Connection{
		ID:      string(time.Now().Unix()),
		Details: c.Details,
		Store:   c.Store,
		TX:      c.TX,
	}
}

func (c *Connection) timeFunc(name string, fn func() error) error {
	now := time.Now()
	err := fn()
	atomic.AddInt64(&c.Elapsed, int64(time.Now().Sub(now)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
