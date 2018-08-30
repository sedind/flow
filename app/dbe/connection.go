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
	Store   Store
	Dialect Dialect
	Elapsed int64
	TX      *Tx
}

// NewConnection creates a new connection, and sets it's `Dialect`
// appropriately based on the `Details` passed into it.
func NewConnection(details *Details) (*Connection, error) {
	err := details.Finalize()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := &Connection{
		ID: string(time.Now().Unix()),
	}
	switch details.Dialect {
	case "mysql":
		c.Dialect = newMySQLDialect(details)
	}

	return c, nil
}

// Open creates new datasource connection
func (c *Connection) Open() error {
	if c.Store != nil {
		return nil
	}

	dbc, err := sqlx.Open(c.Dialect.Name(), c.Dialect.URL())
	if err != nil {
		return errors.WithStack(err)
	}
	dbc.SetMaxOpenConns(c.Dialect.Details().Pool)
	dbc.SetMaxIdleConns(c.Dialect.Details().IdlePool)

	c.Store = &DB{dbc}
	return nil

}

// Close destroys an active datasource connection
func (c *Connection) Close() error {
	return errors.Wrap(c.Store.Close(), "could not close connection")
}

// NewTransaction starts a new transaction on the connection
func (c *Connection) NewTransaction() (*Connection, error) {
	if c.TX == nil {
		tx, err := c.Store.Transaction()
		if err != nil {
			return nil, errors.Wrap(err, "could not start a new transaction")
		}

		conn := &Connection{
			ID:      string(time.Now().Unix()),
			Store:   tx,
			Dialect: c.Dialect,
			TX:      tx,
		}
		return conn, nil
	}
	return c, nil
}

// Transaction will start a new transaction on the connection. If the inner function
// returns an error then the transaction will be rolled back, otherwise the transaction
// will automatically commit at the end.
func (c *Connection) Transaction(fn func(tx *Connection) error) error {
	return c.Dialect.Lock(func() error {
		var dbErr error
		txc, err := c.NewTransaction()
		if err != nil {
			return err
		}
		err = fn(txc)
		if err != nil {
			dbErr = txc.TX.Rollback()
			return errors.WithStack(err)
		}

		dbErr = txc.TX.Commit()

		return errors.Wrap(dbErr, "error commiting transaction")
	})
}

// Query creates a new "empty" query for the connection.
func (c *Connection) Query() *Query {
	return NewQuery(c)
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
