package dbe

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// DefaultForeignKeyNamer contains the default foreign key name generator method
type DefaultForeignKeyNamer struct {
}

type commonDialect struct {
	store   *sql.DB
	details Details
	DefaultForeignKeyNamer
}

func init() {
	RegisterDialect("common", &commonDialect{})
}

func (commonDialect) Name() string {
	return "common"
}

// SetDetails - sets connection details
func (c commonDialect) SetDetails(details Details) {
	c.details = details
}

// Details - returns connection details
func (c commonDialect) Details() Details {
	return c.details
}

// URL - returns conenction string URL
func (c commonDialect) URL() string {
	return ""
}

func (c commonDialect) SetStore(store *sql.DB) {
	c.store = store
}

func (commonDialect) BindVar(i int) string {
	return "$$$" // ?
}

func (commonDialect) Quote(key string) string {
	return fmt.Sprintf(`"%s"`, key)
}

func (c commonDialect) HasIndex(tableName string, indexName string) bool {
	var count int
	currentDatabase, tableName := c.currentDatabaseAndTable(tableName)
	c.store.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.STATISTICS WHERE table_schema = ? AND table_name = ? AND index_name = ?", currentDatabase, tableName, indexName).Scan(&count)
	return count > 0
}

func (c commonDialect) RemoveIndex(tableName string, indexName string) error {
	_, err := c.store.Exec(fmt.Sprintf("DROP INDEX %v", indexName))
	return err
}

func (c commonDialect) HasForeignKey(tableName string, foreignKeyName string) bool {
	return false
}

func (c commonDialect) HasTable(tableName string) bool {
	var count int
	currentDatabase, tableName := c.currentDatabaseAndTable(tableName)
	c.store.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ? AND table_name = ?", currentDatabase, tableName).Scan(&count)
	return count > 0
}

func (c commonDialect) HasColumn(tableName string, columnName string) bool {
	var count int
	currentDatabase, tableName := c.currentDatabaseAndTable(tableName)
	c.store.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = ? AND table_name = ? AND column_name = ?", currentDatabase, tableName, columnName).Scan(&count)
	return count > 0
}

func (c commonDialect) ModifyColumn(tableName string, columnName string, typ string) error {
	_, err := c.store.Exec(fmt.Sprintf("ALTER TABLE %v ALTER COLUMN %v TYPE %v", tableName, columnName, typ))
	return err
}

func (c commonDialect) CurrentDatabase() (name string) {
	c.store.QueryRow("SELECT DATABASE()").Scan(&name)
	return
}

func (c commonDialect) LimitAndOffsetSQL(limit, offset interface{}) (sql string) {
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

func (c commonDialect) LastInsertIDReturningSuffix(tableName string, columnName string) string {
	return ""
}

func (c commonDialect) DefaultValueStr() string {
	return "DEFAULT VALUES"
}

func (c commonDialect) Lock(fn func() error) error {
	return fn()
}

func (c commonDialect) currentDatabaseAndTable(tableName string) (string, string) {
	if strings.Contains(tableName, ".") {
		splitStrings := strings.SplitN(tableName, ".", 2)
		return splitStrings[0], splitStrings[1]
	}
	return c.CurrentDatabase(), tableName
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
