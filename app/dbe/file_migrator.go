package dbe

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// FileMigrator is a migrator for SQL
// files on disk at a specified path.
type FileMigrator struct {
	Migrator
	Path string
}

// NewFileMigrator for a path and a Connection
func NewFileMigrator(path string, conn *Connection) (FileMigrator, error) {
	fm := FileMigrator{
		Migrator: NewMigrator(conn),
		Path:     path,
	}

	fm.SchemaPath = path

	err := fm.loadMigrations()
	if err != nil {
		return fm, errors.WithStack(err)
	}

	return fm, nil
}

func (fm *FileMigrator) loadMigrations() error {
	dir := fm.Path
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		return nil
	}

	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			matches := migrationRegEx.FindAllStringSubmatch(info.Name(), -1)
			if matches == nil || len(matches) == 0 {
				return nil
			}

			match := matches[0]
			dir := match[4]

			migration := Migration{
				Path:      p,
				Version:   match[1],
				Name:      match[2],
				Direction: dir,
				Runner: func(migration Migration, conn *Connection) error {
					f, err := os.Open(p)
					if err != nil {
						return errors.WithStack(err)
					}

					content, err := migrateContent(migration, f)
					if err != nil {
						return errors.Wrapf(err, "error processing %s", migration.Path)
					}

					if content == "" {
						return nil
					}

					tx, err := conn.DB.Begin()
					if err != nil {
						return err
					}

					defer func() {
						if err != nil {
							tx.Rollback()
							fmt.Printf("Migration %s Failed. Rolling back...", migration.Name)
							return
						}
						tx.Commit()
					}()

					_, err = tx.Exec(content)
					if err != nil {
						return errors.Wrapf(err, "error executing %s, sql: %s", migration.Path, content)
					}

					return nil
				}, // Runner end
			} // Migration end

			fm.Migrations[dir] = append(fm.Migrations[dir], migration)

		}
		return nil
	})
	return nil
}

func migrateContent(m Migration, r io.Reader) (string, error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return "", errors.WithStack(err)
	}

	content := string(raw)

	temp := template.Must(template.New("sql").Parse(content))

	var buff bytes.Buffer

	err = temp.Execute(&buff, nil)
	if err != nil {
		return "", errors.Wrapf(err, "could not execute migration template %s", m.Path)
	}

	content = buff.String()

	return content, nil

}
