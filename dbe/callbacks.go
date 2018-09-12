package dbe

// BeforeCreater callback will be called before record is created
type BeforeCreater interface {
	BeforeCreate(*Connection) error
}

func (m *Model) beforeCreate(c *Connection) error {
	if cm, ok := m.Value.(BeforeCreater); ok {
		return cm.BeforeCreate(c)
	}
	return nil
}

// AfterCreater callback will be called after record is created
type AfterCreater interface {
	AfterCreate(*Connection) error
}

func (m *Model) afterCreate(c *Connection) error {
	if cm, ok := m.Value.(AfterCreater); ok {
		return cm.AfterCreate(c)
	}
	return nil
}

// BeforeDeleter callback will be called before record is deleted
type BeforeDeleter interface {
	AfterDelete(*Connection) error
}

func (m *Model) beforeDelete(c *Connection) error {
	if cm, ok := m.Value.(BeforeDeleter); ok {
		return cm.AfterDelete(c)
	}
	return nil
}

// AfterDeleter callback will be called after record is deleted
type AfterDeleter interface {
	AfterDelete(*Connection) error
}

func (m *Model) afterDelete(c *Connection) error {
	if cm, ok := m.Value.(AfterDeleter); ok {
		return cm.AfterDelete(c)
	}
	return nil
}

// BeforeUpdater callback will be called before record is updated
type BeforeUpdater interface {
	BeforeUpdate(*Connection) error
}

func (m *Model) beforeUpdate(c *Connection) error {
	if cm, ok := m.Value.(BeforeUpdater); ok {
		return cm.BeforeUpdate(c)
	}
	return nil
}

// AfterUpdater callback will be called after record is updated
type AfterUpdater interface {
	AfterUpdate(*Connection) error
}

func (m *Model) afterUpdate(c *Connection) error {
	if cm, ok := m.Value.(AfterUpdater); ok {
		return cm.AfterUpdate(c)
	}
	return nil
}
