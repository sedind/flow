package dbe

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sedind/flow/validate"
)

// Create add a new given entry to the database, excluding the given columns.
// It updates `created_at` and `updated_at` columns automatically.
func (c *Connection) Create(model interface{}, excludeColumns ...string) error {

	m := &Model{Value: model}

	if err := m.beforeCreate(c); err != nil {
		return err
	}

	m.TouchCreatedAt()
	m.TouchUpdatedAt()

	cols := m.Columns()
	cols.Remove(excludeColumns...)

	stmt, err := c.Dialect.CreateStmt(m.TableName(), cols.String(), cols.ParamString())
	if err != nil {
		return errors.WithStack(err)
	}

	Logger.Info(stmt)

	res, err := c.Store.NamedExec(stmt, m.Value)
	if err != nil {
		return errors.WithStack(err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return errors.WithStack(err)
	}
	m.setID(id)

	return m.afterCreate(c)
}

// Delete given model from database
func (c *Connection) Delete(model interface{}) error {
	m := &Model{Value: model}

	if err := m.beforeDelete(c); err != nil {
		return err
	}

	stmt, err := c.Dialect.DeleteStmt(m.TableName(), m.WhereID())

	if err != nil {
		return errors.WithStack(err)
	}

	Logger.Info(stmt)

	_, err = c.Store.Exec(stmt)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := m.afterDelete(c); err != nil {
		return err
	}
	return nil
}

// Save wraps the Create and Update methods. It executes a Create if no ID is provided with the entry;
// or issues an Update otherwise.
func (c *Connection) Save(model interface{}, excludeColumns ...string) error {
	m := &Model{Value: model}
	id := m.ID()
	if fmt.Sprint(id) == "0" {
		return c.Create(m.Value, excludeColumns...)
	}
	return c.Update(m.Value, excludeColumns...)

}

// Update given model to database
func (c *Connection) Update(model interface{}, excludeColumns ...string) error {
	m := &Model{Value: model}
	if err := m.beforeUpdate(c); err != nil {
		return err
	}

	cols := m.Columns()
	cols.Remove("id", "created_at")
	cols.Remove(excludeColumns...)

	m.TouchUpdatedAt()

	stmt, err := c.Dialect.UpdateStmt(m.TableName(), cols.UpdateString(), m.WhereID())

	if err != nil {
		return errors.WithStack(err)
	}

	Logger.Info(stmt)

	_, err = c.Store.NamedExec(stmt, m.Value)

	if err != nil {
		return errors.WithStack(err)
	}

	return m.afterUpdate(c)
}

// ValidateAndCreate applies validation rules on the given entry, then creates it
// if the validation succeed, excluding the given columns.
func (c *Connection) ValidateAndCreate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	m := &Model{Value: model}

	verrs, err := m.validateCreate(c)
	if err != nil {
		return verrs, err
	}

	if verrs.HasAny() {
		return verrs, nil
	}

	return verrs, c.Create(m.Value, excludeColumns...)
}

// ValidateAndSave applies validation rules on the given entry, then save it
// if the validation succeed, excluding the given columns.
func (c *Connection) ValidateAndSave(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	m := &Model{Value: model}
	id := m.ID()
	if fmt.Sprint(id) == "0" {
		return c.ValidateAndCreate(m.Value, excludeColumns...)
	}
	return c.ValidateAndUpdate(m.Value, excludeColumns...)
}

// ValidateAndUpdate applies validation rules on the given entry, then update it
// if the validation succeed, excluding the given columns.
func (c *Connection) ValidateAndUpdate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	m := &Model{Value: model}

	verrs, err := m.validateUpdate(c)
	if err != nil {
		return verrs, err
	}

	if verrs.HasAny() {
		return verrs, nil
	}

	return verrs, c.Update(m.Value, excludeColumns...)
}

// ValidateAndDelete applies validation rules on the given entry,
// then delte it if the validation succeed.
func (c *Connection) ValidateAndDelete(model interface{}) (*validate.Errors, error) {
	m := &Model{Value: model}
	verrs, err := m.validateDelete(c)
	if err != nil {
		return verrs, err
	}

	if verrs.HasAny() {
		return verrs, nil
	}

	return verrs, c.Delete(m.Value)
}
