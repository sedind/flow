package dbe

import (
	"github.com/pkg/errors"
	"github.com/sedind/flow/validate"
)

// BeforeValidator ensures BeforValidate call on model which implements Validator interface
type BeforeValidator interface {
	BeforeValidate(c *Connection) error
}

// Validator adds Validation support on Model
type Validator interface {
	Validate(c *Connection) (*validate.Errors, error)
}

// CreateValidator adds Validation support for Create method on model
type CreateValidator interface {
	ValidateCreate(c *Connection) (*validate.Errors, error)
}

// UpdateValidator adds Validation support for Update method on model
type UpdateValidator interface {
	ValidateUpdate(c *Connection) (*validate.Errors, error)
}

// DeleteValidator adds Validation support for Delete method on model
type DeleteValidator interface {
	ValidateDelete(c *Connection) (*validate.Errors, error)
}

func (m *Model) validate(c *Connection) (*validate.Errors, error) {
	if mv, ok := m.Value.(BeforeValidator); ok {
		if err := mv.BeforeValidate(c); err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}
	}

	if mv, ok := m.Value.(Validator); ok {
		return mv.Validate(c)
	}

	return validate.NewErrors(), nil
}

func (m *Model) validateCreate(c *Connection) (*validate.Errors, error) {
	verrs, err := m.validate(c)
	if err != nil {
		return verrs, errors.WithStack(err)
	}

	if mv, ok := m.Value.(CreateValidator); ok {
		v, err := mv.ValidateCreate(c)
		if v != nil {
			verrs.Append(v)
		}
		if err != nil {
			return verrs, errors.WithStack(err)
		}
	}
	return verrs, err
}

func (m *Model) validateUpdate(c *Connection) (*validate.Errors, error) {
	verrs, err := m.validate(c)
	if err != nil {
		return verrs, errors.WithStack(err)
	}

	if mv, ok := m.Value.(UpdateValidator); ok {
		v, err := mv.ValidateUpdate(c)
		if v != nil {
			verrs.Append(v)
		}
		if err != nil {
			return verrs, errors.WithStack(err)
		}
	}
	return verrs, err
}

func (m *Model) validateDelete(c *Connection) (*validate.Errors, error) {
	verrs, err := m.validate(c)
	if err != nil {
		return verrs, errors.WithStack(err)
	}

	if mv, ok := m.Value.(DeleteValidator); ok {
		v, err := mv.ValidateDelete(c)
		if v != nil {
			verrs.Append(v)
		}
		if err != nil {
			return verrs, errors.WithStack(err)
		}
	}
	return verrs, err
}
