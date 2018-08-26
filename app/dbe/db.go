package dbe

import (
	"database/sql"
)

// DB - dbe DB object which wraps native sql.DB
type DB struct {
	*sql.DB
}
