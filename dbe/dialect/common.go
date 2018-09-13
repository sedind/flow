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

// UpdateStmt createse SQL INSER statement
func (c Common) UpdateStmt(tableName string, columns string, where string) (string, error) {

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, columns, where)
	return query, nil
}

// DeleteStmt creates SQL DELETE statement
func (c Common) DeleteStmt(tableName string, where string) (string, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, where)
	return query, nil
}

// TranslateSQL to supported dialect
func (c Common) TranslateSQL(sql string) string {
	return sql
}
