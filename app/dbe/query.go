package dbe

// Query object is used to build up a query
// to be executed against the `Connection`.
type Query struct {
}

// NewQuery creates a new "empty" query from the current connection.
func NewQuery(c *Connection) *Query {
	return &Query{}
}
