package dialect

import (
	"fmt"
)

var _ Dialect = Common{}

// Common describes common operations for supported dialects
type Common struct {
}

// Name for current dialect
func (c Common) Name() string {
	return "common"
}

// CreateStmt createse SQL INSER statement
func (c Common) CreateStmt(tableName string, columns string, columnNames string) (string, error) {

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, columnNames)
	return query, nil
}
