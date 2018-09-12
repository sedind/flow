package model

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/sedind/inflect"

	"github.com/pkg/errors"
)

var tableMap = map[string]string{}
var tableMapMutex = sync.RWMutex{}

const modelTag = "db"

// Value holds content of a Model object
type Value interface{}

// TableNamer interface allows for the customize table mapping
// between a name and the database. For example the value
// `User{}` will automatically map to "users". Implementing `TableNamer`
// would allow this to change to be changed to whatever you would like.
type TableNamer interface {
	TableName() string
}

// Model is used to wrap user Model objects
// and use them in Query objects to execute DB queries
type Model struct {
	Value
	tableName string
	As        string
}

// New Creates new DBE Model for given model object
func New(model interface{}) *Model {
	return &Model{
		Value: model,
	}
}

// ID returns Unique Identifier of the Model
func (m *Model) ID() interface{} {
	fVal, err := m.fieldByName("ID")
	if err != nil {
		return 0
	}
	return fVal.Interface()
}

// TableName returns the corresponding name of the underlying database table
// for a given `Model`. See also `TableNamer` to change the default name of the table.
func (m *Model) TableName() string {
	//check if Model content is type string
	if s, ok := m.Value.(string); ok {
		return s
	}

	// check if Model content implements TableNamer interface
	if n, ok := m.Value.(TableNamer); ok {
		return n.TableName()
	}

	if m.tableName != "" {
		return m.tableName
	}

	t := reflect.TypeOf(m.Value)
	name := m.typeName(t)

	tableMapMutex.Lock()
	defer tableMapMutex.Unlock()

	if tableMap[name] == "" {
		m.tableName = inflect.Tableize(name)
		tableMap[name] = m.tableName
	}

	return tableMap[name]
}

// TouchCreatedAt sets current time to CreatedAt field
func (m *Model) TouchCreatedAt() {
	fbn, err := m.fieldByName("CreatedAt")
	if err == nil {
		fbn.Set(reflect.ValueOf(time.Now()))
	}
}

// TouchUpdatedAt sets current time to UpdatedAt field
func (m *Model) TouchUpdatedAt() {
	fbn, err := m.fieldByName("UpdatedAt")
	if err == nil {
		fbn.Set(reflect.ValueOf(time.Now()))
	}
}

// Columns gets model database Columns using struct db tag
func (m *Model) Columns() []string {
	alias := m.As
	if alias == "" {
		alias = m.TableName()
	}
	cols := []string{}
	t := reflect.TypeOf(m.Value)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(modelTag)
		if tag != "" && tag != "-" {
			col := fmt.Sprintf("%s.%s", alias, tag)
			cols = append(cols, col)
		}
	}

	return cols
}

// WhereID constructs string for WHERE clause in query
// in table_name.id = idValue
func (m *Model) WhereID() string {
	id := m.ID()
	var value string
	switch id.(type) {
	case int, int64:
		value = fmt.Sprintf("%s.id = %d", m.TableName(), id)
	default:
		value = fmt.Sprintf("%s.id ='%s'", m.TableName(), id)
	}
	return value
}

func (m *Model) typeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		el := t.Elem()
		if el.Kind() == reflect.Ptr {
			el = el.Elem()
		}

		// check if elem implements TableNamer interface
		tableNamer := (*TableNamer)(nil)
		if el.Implements(reflect.TypeOf(tableNamer).Elem()) {
			v := reflect.New(el)
			out := v.MethodByName("TableName").Call([]reflect.Value{})
			name := out[0].String()
			if tableMap[el.Name()] == "" {
				tableMap[el.Name()] = name
			}
		}
		return el.Name()
	default:
		return t.Name()
	}
}

func (m *Model) fieldByName(s string) (reflect.Value, error) {
	el := reflect.ValueOf(m.Value).Elem()
	fVal := el.FieldByName(s)
	if !fVal.IsValid() {
		return fVal, errors.Errorf("Model does not have a field %s", s)
	}
	return fVal, nil
}
