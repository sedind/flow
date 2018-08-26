package dbe

import "github.com/pkg/errors"

// Dialect interface contains behaviors that differ across SQL database
type Dialect interface {
	// Name - returns dialect's name
	Name() string

	// SetStore - sets DB store for dialect
	SetStore(store Store)

	// Store returns store object
	Store() Store

	// URS connection URL
	URL() string

	//Details gets connection details
	Details() Details

	// BindVar return the placeholder for actual values in SQL statements, in many dbs it is "?", Postgres using $1
	BindVar(i int) string

	// Quote - quotes field name to avoid SQL parsing excwptions by using reserved words as a field name
	Quote(key string) string

	// HasIndex check has index or not
	HasIndex(tableName string, indexName string) bool

	// HasForeignKey - check has foreign key or not
	HasForeignKey(tableName string, foreignKeyName string) bool

	// RemoveIndex - removes index from table
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

func newDialect(details *Details) (Dialect, error) {
	switch details.Dialect {
	case "mysql":
		return &DialectMySQL{
			DialectCommon: DialectCommon{
				details: *details,
			},
		}, nil

	}
	return nil, errors.Errorf("unsupported dialect %s", details.Dialect)
}
