package dbe

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Columns represents DB table columns
type Columns struct {
	Cols       map[string]*Column
	lock       *sync.RWMutex
	TableName  string
	TableAlias string
}

// NewColumns creates Columns object for given table name
func NewColumns(tableName string) Columns {
	return NewColumnsWithAlias(tableName, "")
}

// NewColumnsWithAlias creates Columns object for given table name and table alias
func NewColumnsWithAlias(tableName string, tableAlias string) Columns {
	return Columns{
		lock:       &sync.RWMutex{},
		Cols:       map[string]*Column{},
		TableName:  tableName,
		TableAlias: tableAlias,
	}
}

// Add new columns
func (c *Columns) Add(names ...string) {
	c.lock.Lock()

	tableAlias := c.alias()

	for _, name := range names {
		name = strings.TrimSpace(name)

		col := c.Cols[name]
		if col == nil {
			cn := name
			if tableAlias != "" {
				cn = fmt.Sprintf("%s.%s", tableAlias, name)
			}

			col = &Column{
				Name: cn,
			}
			c.Cols[name] = col
		}

	}

	c.lock.Unlock()
}

// Remove columns
func (c *Columns) Remove(names ...string) {
	for _, name := range names {
		name = strings.TrimSpace(name)
		delete(c.Cols, name)
	}
}

// alias returns alias used for table
func (c Columns) alias() string {
	tableAlias := c.TableAlias
	if tableAlias == "" {
		tableAlias = c.TableName
	}
	return tableAlias
}

func (c Columns) String() string {
	cols := []string{}
	for _, col := range c.Cols {
		cols = append(cols, col.Name)
	}
	sort.Strings(cols)
	return strings.Join(cols, ", ")
}

// ParamString returns columns as parameter strings
func (c Columns) ParamString() string {

	cols := []string{}
	for _, col := range c.Cols {
		name := col.Name
		if strings.Contains(name, ".") {
			tmp := strings.Split(name, ".")
			name = tmp[1]
		}
		cols = append(cols, ":"+name)
	}
	sort.Strings(cols)
	return strings.Join(cols, ", ")
}

func (c Columns) UpdateString() string {
	cols := []string{}
	for _, c := range c.Cols {
		cols = append(cols, c.UpdateString())
	}
	sort.Strings(cols)
	return strings.Join(cols, ", ")
}
