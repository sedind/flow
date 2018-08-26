package dbe

import (
	"database/sql"
	"fmt"
	"reflect"
)

// Dialect interface contains behaviors that differ across SQL database
type Dialect interface {
	// Name - returns dialect's name
	Name() string

	// SetDetails - sets connection details
	SetDetails(details Details)

	// Details - returns connection details
	Details() Details

	// URL - returns conenction string URL
	URL() string

	// SetStore - sets DB store for dialect
	SetStore(store *sql.DB)

	// BindVar return the placeholder for actual values in SQL statements, in many dbs it is "?", Postgres using $1
	BindVar(i int) string

	// Quote - quotes field name to avoid SQL parsing excwptions by using reserved words as a field name
	Quote(key string) string

	// HasIndex check has index or not
	HasIndex(tableName string, indexName string) bool

	// HasForeignKey - check has foreign key or not
	HasForeignKey(tableName string, foreignKeyName string) bool

	// Remove index
	RemoveIndex(tableName string, indexName string) error

	// HasTable - check has table or not
	HasTable(tableName string) bool

	// HasColumn - check has column or not
	HasColumn(tableName string, columnName string) bool

	// ModifyColumn - modify column's type
	ModifyColumn(tableName string, columnName string, typ string) error

	// LimitAndOffsetSQL - returns generated SQL with Limit and Offset
	LimitAndOffsetSQL(limit, offset interface{}) string

	// LastInsertIDReturningSuffix - most dbs support LastInsertId, butp postgress needs to use 'RETURNING'
	LastInsertIDReturningSuffix(tableName string, columnName string) string

	// DefaultValueStr
	DefaultValueStr() string

	// BuildKeyName returns a valid key name (foreign key, index key) for the given table, field and reference
	BuildKeyName(kind, tableName string, fields ...string) string

	// CurrentDatabase - returns current database name
	CurrentDatabase() string

	// Lock - lock database
	Lock(func() error) error
}

var dialectsMap = map[string]Dialect{}

func newDialect(details Details) (Dialect, error) {
	if val, ok := dialectsMap[details.Dialect]; ok {
		dialect := reflect.New(reflect.TypeOf(val).Elem()).Interface().(Dialect)
		dialect.SetDetails(details)
		return dialect, nil
	}

	err := fmt.Errorf("'%s' dialect is not supported", details.Dialect)
	return nil, err
}

// RegisterDialect - registers new dialect
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}
