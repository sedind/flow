package dbe

import (
	"fmt"

	"github.com/pkg/errors"
)

// Migration handles the data for a given database migration
type Migration struct {
	// Path to the migration (./migrations/20180515170142_initial.up.sql)
	Path string
	// Version of the migration (20180515170142)
	Version string
	// Name of the migration (initial)
	Name string
	// Direction of the migration (up)
	Direction string
	// Runner function to run/execute the migration
	Runner func(Migration, *Connection) error
}

// Run the migration. Returns an error if there is
// no mf.Runner defined.
func (m Migration) Run(conn *Connection) error {
	if m.Runner == nil {
		return errors.Errorf("no runner defined for $s", m.Path)
	}

	return m.Runner(m, conn)
}

// Exists checks if migration exists in DB
func (m Migration) Exists(conn *Connection, migrationTable string) (bool, error) {
	var count int
	err := conn.DB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE version = ?", migrationTable), m.Version).Scan(&count)

	if err != nil {
		return false, errors.WithStack(err)
	}

	return count > 0, err
}

// Migrations is a collection of Migration
type Migrations []Migration

// Len returns number of migrations
func (m Migrations) Len() int {
	return len(m)
}

// Less checks if migration version at index i is lesss than migration at index j
func (m Migrations) Less(i, j int) bool {
	return m[i].Version < m[j].Version
}

// Swap migrations at index i and j
func (m Migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
