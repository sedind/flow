package dbe

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// DefaultForeignKeyNamer contains the default foreign key name generator method
type DefaultForeignKeyNamer struct {
}

//DialectCommon  implements common dialect operations
type DialectCommon struct {
	store   Store
	details Details
	DefaultForeignKeyNamer
}

// Name returns dialect name
func (dc *DialectCommon) Name() string {
	return "common"
}

// SetStore sets dialect DB store
func (dc *DialectCommon) SetStore(store Store) {
	dc.store = store
}

// Store returns store instance
func (dc *DialectCommon) Store() Store {
	return dc.store
}

//Details gets connection details
func (dc *DialectCommon) Details() Details {
	return dc.details
}

// BindVar return the placeholder for actual values in SQL statements, in many dbs it is "?", Postgres using $1
func (dc *DialectCommon) BindVar(i int) string {
	return "$$$" // ?
}

// Quote - quotes field name to avoid SQL parsing excwptions by using reserved words as a field name
func (dc *DialectCommon) Quote(key string) string {
	return fmt.Sprintf(`"%s"`, key)
}

// HasIndex check has index or not
func (dc *DialectCommon) HasIndex(tableName string, indexName string) bool {
	var count int
	currentDatabase, tableName := dc.currentDatabaseAndTable(tableName)
	dc.store.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE table_schema = ? AND table_name = ? AND index_name = ?", currentDatabase, tableName, indexName).Scan(&count)
	return count > 0
}

// HasForeignKey - check has foreign key or not
func (dc *DialectCommon) HasForeignKey(tableName string, foreignKeyName string) bool {
	return false
}

// RemoveIndex - removes index from table
func (dc *DialectCommon) RemoveIndex(tableName string, indexName string) error {
	_, err := dc.store.Exec(fmt.Sprintf("DROP INDEX %v", indexName))
	return err
}

// HasTable - check has table or not
func (dc *DialectCommon) HasTable(tableName string) bool {
	var count int
	currentDatabase, tableName := dc.currentDatabaseAndTable(tableName)
	dc.store.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ? AND table_name = ?", currentDatabase, tableName).Scan(&count)
	return count > 0
}

// HasColumn - check has column or not
func (dc *DialectCommon) HasColumn(tableName string, columnName string) bool {
	var count int
	currentDatabase, tableName := dc.currentDatabaseAndTable(tableName)
	dc.store.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = ? AND table_name = ? AND column_name = ?", currentDatabase, tableName, columnName).Scan(&count)
	return count > 0
}

// ModifyColumn - modify column's type
func (dc *DialectCommon) ModifyColumn(tableName string, columnName string, typ string) error {
	_, err := dc.store.Exec(fmt.Sprintf("ALTER TABLE %v ALTER COLUMN %v TYPE %v", tableName, columnName, typ))
	return err
}

// LimitAndOffsetSQL - returns generated SQL with Limit and Offset
func (dc *DialectCommon) LimitAndOffsetSQL(limit, offset interface{}) (sql string) {
	if limit != nil {
		if parsedLimit, err := strconv.ParseInt(fmt.Sprint(limit), 0, 0); err == nil && parsedLimit >= 0 {
			sql += fmt.Sprintf(" LIMIT %d", parsedLimit)
		}
	}
	if offset != nil {
		if parsedOffset, err := strconv.ParseInt(fmt.Sprint(offset), 0, 0); err == nil && parsedOffset >= 0 {
			sql += fmt.Sprintf(" OFFSET %d", parsedOffset)
		}
	}
	return
}

// LastInsertIDReturningSuffix - most dbs support LastInsertId, butp postgress needs to use 'RETURNING'
func (dc *DialectCommon) LastInsertIDReturningSuffix(tableName string, columnName string) string {
	return ""
}

// DefaultValueStr -
func (dc *DialectCommon) DefaultValueStr() string {
	return "DEFAULT VALUES"
}

// CurrentDatabase - returns current database name
func (dc *DialectCommon) CurrentDatabase() (name string) {
	dc.store.QueryRow("SELECT DATABASE()").Scan(&name)
	return
}

// Lock -
func (dc *DialectCommon) Lock(fn func() error) error {
	return fn()
}

func (dc *DialectCommon) currentDatabaseAndTable(tableName string) (string, string) {
	if strings.Contains(tableName, ".") {
		splitStrings := strings.SplitN(tableName, ".", 2)
		return splitStrings[0], splitStrings[1]
	}
	return dc.CurrentDatabase(), tableName
}

// BuildKeyName returns a valid key name (foreign key, index key) for the given table, field and reference
func (DefaultForeignKeyNamer) BuildKeyName(kind, tableName string, fields ...string) string {
	keyName := fmt.Sprintf("%s_%s_%s", kind, tableName, strings.Join(fields, "_"))
	keyName = regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(keyName, "_")
	return keyName
}

// IsByteArrayOrSlice returns true of the reflected value is an array or slice
func IsByteArrayOrSlice(value reflect.Value) bool {
	return (value.Kind() == reflect.Array || value.Kind() == reflect.Slice) && value.Type().Elem() == reflect.TypeOf(uint8(0))
}
