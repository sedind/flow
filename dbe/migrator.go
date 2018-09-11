package dbe

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/pkg/errors"
)

var migrationRegEx = regexp.MustCompile(`(\d+)_([^\.]+)(\.[a-z]+)?\.(up|down)\.(sql)`)

// default migartion schema name
const defaultMigrationSchema string = "schema_migration"

// NewMigrator returns a new "blank" migrator.
// A "blank" Migrator should only be used as
//the basis for a new type of migration system.
// It is recommended to use something like FileMigrator.
func NewMigrator(conn *Connection) Migrator {
	return Migrator{
		Conn: conn,
		Migrations: map[string]Migrations{
			"up":   Migrations{},
			"down": Migrations{},
		},
	}
}

// Migrator forms the basis of all migrations systems.
// It does the actual heavy lifting of running migrations.
// When building a new migration system, you should embed this
// type into your migrator.
type Migrator struct {
	Conn       *Connection
	SchemaPath string
	Migrations map[string]Migrations
}

// Up runs pending "up" migrations and applies them to the database.
func (m Migrator) Up() error {
	return m.exec(func() error {

		migrations := m.Migrations["up"]
		sort.Sort(migrations)
		for _, migration := range migrations {
			exists, err := migration.Exists(m.Conn, m.migrationSchema())
			if err != nil {
				return errors.Wrapf(err, "problem checking for migration version %s", migration.Version)
			}

			if exists {
				continue //migration is executed skip to next
			}

			err = migration.Run(m.Conn)
			if err != nil {
				return errors.WithStack(err)
			}

			_, err = m.Conn.Store.Exec(fmt.Sprintf("insert into %s (version,name) values (?,?)", m.migrationSchema()), migration.Version, migration.Name)
			if err != nil {
				return errors.WithStack(err)
			}

			fmt.Printf("> %s\n", migration.Name)
		}
		return nil
	})
}

// Down runs pending "down" migrations and rolls back the
// database by the specified number of steps.
func (m Migrator) Down(step int) error {
	return m.exec(func() error {
		count, err := m.getExecutedMigrationsCount()
		if err != nil {
			return errors.Wrap(err, "migration down: unable to count migrations")
		}

		migrations := m.Migrations["down"]
		// sorting magic :)
		sort.Sort(sort.Reverse(migrations))

		//skip all executed migrations
		if len(migrations) > count {
			migrations = migrations[len(migrations)-count:]
		}

		// and run only required steps
		if step > 0 && len(migrations) >= step {
			migrations = migrations[:step]
		}
		for _, migration := range migrations {
			exists, err := migration.Exists(m.Conn, m.migrationSchema())
			if err != nil || !exists {
				return errors.Wrapf(err, "problem checking for migration version %s", migration.Version)
			}

			err = migration.Run(m.Conn)
			if err != nil {
				return errors.WithStack(err)
			}

			_, err = m.Conn.Store.Exec(fmt.Sprintf("delete from %s where version = ? ", m.migrationSchema()), migration.Version)
			if err != nil {
				return errors.WithStack(err)
			}
			fmt.Printf("< %s\n", migration.Name)

		}
		return nil
	})
}

// Status prints out the status of applied/pending migrations.
func (m Migrator) Status() error {
	return m.exec(func() error {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "Version\t\tName\t\tStatus")
		for _, migration := range m.Migrations["up"] {
			exists, err := migration.Exists(m.Conn, m.migrationSchema())
			if err != nil {
				return errors.Wrapf(err, "problem checking for migration version %s", migration.Version)
			}
			status := "Pending"
			if exists {
				status = "Applied"
			}
			fmt.Fprintf(w, "%s\t\t%s\t\t%s\t\t\n", migration.Version, migration.Name, status)
		}
		return w.Flush()
	})
}

// Reset the database by runing the down migrations followed by the up migrations.
func (m Migrator) Reset() error {
	err := m.Down(-1)
	if err != nil {
		return errors.WithStack(err)
	}
	return m.Up()
}

// CreateSchemaMigrations sets up a table to track migrations.
func (m Migrator) CreateSchemaMigrations() error {
	//check if migrations table exists
	if m.hasMigrationsSchema() {
		return nil // nothing to to here, go back
	}
	dialect := m.Conn.Details.Dialect
	sql := m.getMigrationsSchema(dialect)

	if sql == "" {
		return errors.Errorf("Version Schema missing for dialect %s", dialect)
	}

	_, err := m.Conn.Store.Exec(sql)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// exec internal helper execution function
func (m Migrator) exec(fn func() error) error {
	now := time.Now()
	defer printTimer(now)

	err := m.CreateSchemaMigrations()
	if err != nil {
		return errors.Wrap(err, "Migrator: problem creating schema migrations")
	}

	return fn()
}

func (m Migrator) hasMigrationsSchema() bool {
	var currentDatabase string
	var count int

	m.Conn.Store.QueryRow("SELECT DATABASE()").Scan(&currentDatabase)

	m.Conn.Store.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ? AND table_name = ?", currentDatabase, m.migrationSchema()).Scan(&count)
	return count > 0
}

func (m Migrator) getMigrationsSchema(dialect string) string {
	switch dialect {
	case "mysql":
		return fmt.Sprintf(mySQLMigrationTblTpl, m.migrationSchema())
	}
	return ""
}

func (m Migrator) getExecutedMigrationsCount() (int, error) {
	var count int
	err := m.Conn.Store.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", m.migrationSchema())).Scan(&count)

	return count, err
}

func (m Migrator) migrationSchema() string {
	migrationTable := defaultMigrationSchema
	if v, ok := m.Conn.Details.Options["migration_schema"]; ok {
		migrationTable = v
	}
	return migrationTable
}

// printTimer prints time difference between startTime and time of execution
func printTimer(startTime time.Time) {
	diff := time.Now().Sub(startTime).Seconds()
	if diff > 60 {
		fmt.Printf("\n%.4f minutes \n", diff/60)
	} else {
		fmt.Printf("\n%.4f seconds \n", diff)
	}
}

var mySQLMigrationTblTpl = `
	CREATE TABLE %s ( 
	version NVARCHAR(14) NOT NULL, 
	name NVARCHAR(255) NULL, 
	UNIQUE INDEX  schema_version_idx (version ASC));
`
