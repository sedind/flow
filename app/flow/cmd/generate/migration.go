package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/sedind/flow/app/config"

	"github.com/spf13/cobra"
)

// migrationCmd generates sql migration files
var migrationCmd = &cobra.Command{
	Use:   "migration [name]",
	Short: "Generates Up/Down migration files for your SQL database.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must supply a name for your migration")
		}
		if migrationsPath != "" {
			// we will ignore project configuration and use migrations path to generate migration
			return generateMigrationFile(migrationsPath, args[0], "sql", nil, nil)
		}

		if configFile == "" {
			return errors.New("config file not provided")
		}

		var path struct {
			Path string `yaml:"migrations_path"`
		}

		err := config.LoadFromPath(configFile, &path)
		if err != nil {
			return errors.Wrapf(err, "Unable to load configuration %s", configFile)
		}

		if path.Path == "" {
			return errors.New("migrations_path can not be empty in configuration file")
		}

		return generateMigrationFile(path.Path, args[0], "sql", nil, nil)
	},
}

// generateMigrationFile writes contents for a given migration in normalized files
func generateMigrationFile(path, name, ext string, up, down []byte) error {
	n := time.Now().UTC()
	s := n.Format("20060102150405")

	err := os.MkdirAll(path, 0766)
	if err != nil {
		return errors.Wrapf(err, "couldn't create migrations path %s", path)
	}

	upf := filepath.Join(path, (fmt.Sprintf("%s_%s.up.%s", s, name, ext)))
	err = ioutil.WriteFile(upf, up, 0666)
	if err != nil {
		return errors.Wrapf(err, "couldn't write up migration %s", upf)
	}
	fmt.Printf("> %s\n", upf)

	downf := filepath.Join(path, (fmt.Sprintf("%s_%s.down.%s", s, name, ext)))
	err = ioutil.WriteFile(downf, down, 0666)
	if err != nil {
		return errors.Wrapf(err, "couldn't write up migration %s", downf)
	}

	fmt.Printf("> %s\n", downf)
	return nil
}
