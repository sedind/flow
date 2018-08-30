package dbe

var _ Dialect = &MySQLDialect{}

// MySQLDialect represents mySQL dialect implementation
type MySQLDialect struct {
	details *Details
}

// Name returns dialect name
func (d *MySQLDialect) Name() string {
	return "mysql"
}

// URL returns Dialect Connection URL
func (d *MySQLDialect) URL() string {
	return d.details.URL
}

// Details returns Connection Details object
func (d *MySQLDialect) Details() *Details {
	return d.details
}

// Create executes SQL INSERT Query on provided datastore
func (d *MySQLDialect) Create(store Store, query Query) error {
	return nil
}

// Update executes SQL UPDATE Query on provided datastore
func (d *MySQLDialect) Update(store Store, query Query) error {
	return nil
}

// Destroy executes SQL DELETE Query on provided datastore
func (d *MySQLDialect) Destroy(store Store, query Query) error {
	return nil
}

// SelectOne executes SQL SELECT Query on provided datastore with LIMIT 1
func (d *MySQLDialect) SelectOne(store Store, query Query) error {
	return nil
}

// SelectMany executes SQL SELECT Query on provided datastore
func (d *MySQLDialect) SelectMany(store Store, query Query) error {
	return nil
}

// HasIndex check has index or not
func (d *MySQLDialect) HasIndex(store Store, tableName string, indexName string) bool {
	return false
}

// HasForeignKey - check has foreign key or not
func (d *MySQLDialect) HasForeignKey(store Store, tableName string, foreignKeyName string) bool {
	return false
}

// RemoveIndex removes index from table
func (d *MySQLDialect) RemoveIndex(store Store, tableName string, indexName string) error {
	return nil
}

// HasTable - check has table or not
func (d *MySQLDialect) HasTable(store Store, tableName string) bool {
	return false
}

// HasColumn - check has column or not
func (d *MySQLDialect) HasColumn(tableName string, columnName string) bool {
	return false
}

// CreateDB creates database from connection details
func (d *MySQLDialect) CreateDB(Store) error {
	return nil
}

// DropDB drops database from connection details
func (d *MySQLDialect) DropDB(Store) error {
	return nil
}

//Lock implements locking mechanism
func (d *MySQLDialect) Lock(fn func() error) error {
	return fn()
}

// TruncateAll truncates all objects from database
func (d *MySQLDialect) TruncateAll(Store) error {
	return nil
}

func newMySQLDialect(details *Details) Dialect {
	d := &MySQLDialect{details}
	return d
}
