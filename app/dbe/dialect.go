package dbe

// Dialect -
type Dialect interface {
	Name() string
	URL() string
	Details() *Details
	Create(store Store, query Query) error
	Update(store Store, query Query) error
	Destroy(store Store, query Query) error
	SelectOne(store Store, query Query) error
	SelectMany(store Store, query Query) error
	HasIndex(store Store, tableName string, indexName string) bool
	HasForeignKey(store Store, tableName string, foreignKeyName string) bool
	RemoveIndex(store Store, tableName string, indexName string) error
	HasTable(store Store, tableName string) bool
	HasColumn(tableName string, columnName string) bool
	CreateDB(Store) error
	DropDB(Store) error
	Lock(func() error) error
	TruncateAll(Store) error
}
