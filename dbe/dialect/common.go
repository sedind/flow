package dialect

import (
	"fmt"
	"strings"
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
func (c Common) CreateStmt(tableName string, columns []string, columnNames []string) (string, error) {

	cols := strings.Join(columns, ",")
	names := strings.Join(columnNames, ",")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, cols, names)
	return query, nil
}
