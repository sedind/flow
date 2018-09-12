package dialect

import (
	"fmt"
	"reflect"
)

// Dialect defines set of operations that are speccific to different SQL dialects
type Dialect interface {
	Name() string
	CreateStmt(string, []string) (string, error)
}

// list of registered dialects
var dialectsMap = map[string]Dialect{}

// New gets new Dialect instance
func New(name string) (Dialect, error) {
	if value, ok := dialectsMap[name]; ok {
		dialect := reflect.New(reflect.TypeOf(value).Elem()).Interface().(Dialect)
		return dialect, nil
	}
	return nil, fmt.Errorf("'%s' dialect is not supported", name)
}

// RegisterDialect adds new dialect to stack
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}
