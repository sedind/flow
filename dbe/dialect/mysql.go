package dialect

func init() {
	RegisterDialect("mysql", &MySQL{})
}

// MySQL implements dialect speciffic to MySQL DB
type MySQL struct {
	Common
}
