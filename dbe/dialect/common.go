package dialect

import (
	"fmt"
	"strings"
)

var _ Dialect = Common{}

// Common describes common operations for supported dialects
type Common struct {
}

// CreateStmt createse SQL INSER statement
func (c Common) CreateStmt(tableName string, columns []string) (string, error) {
	var colNames []string

	for _, col := range columns {
		colNames = append(colNames, ":"+col)
	}

	cols := strings.Join(columns, ",")
	names := strings.Join(colNames, ",")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, cols, names)
	return query, nil
}
