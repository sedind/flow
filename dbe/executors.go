package dbe

import (
	"fmt"

	"github.com/sedind/flow/validate"
)

// Create add a new given entry to the database, excluding the given columns.
// It updates `created_at` and `updated_at` columns automatically.
func (c *Connection) Create(model interface{}, excludeColumns ...string) error {
	m := &Model{Value: model}
	cols := m.Columns()

	stmt, err := c.Dialect.CreateStmt(m.TableName(), cols)
	if err != nil {
		return err
	}

	fmt.Println(stmt)

	return nil
}

// Delete given model from database
func (c *Connection) Delete(model interface{}) error {
	//m := &Model{Value: model}
	return nil
}

// Save wraps the Create and Update methods. It executes a Create if no ID is provided with the entry;
// or issues an Update otherwise.
func (c *Connection) Save(model interface{}, excludeColumns ...string) error {
	return nil
}

// Update given model to database
func (c *Connection) Update(model interface{}, excludeColumns ...string) error {
	return nil
}

// ValidateAndCreate applies validation rules on the given entry, then creates it
// if the validation succeed, excluding the given columns.
func (c *Connection) ValidateAndCreate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	return nil, nil
}

// ValidateAndSave applies validation rules on the given entry, then save it
// if the validation succeed, excluding the given columns.
func (c *Connection) ValidateAndSave(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	return nil, nil
}

// ValidateAndUpdate applies validation rules on the given entry, then update it
// if the validation succeed, excluding the given columns.
func (c *Connection) ValidateAndUpdate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	return nil, nil
}
