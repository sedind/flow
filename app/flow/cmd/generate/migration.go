package generate

import (
	"github.com/pkg/errors"
	"github.com/sedind/flow/app/config"

	"github.com/sedind/flow/app/dbe/migration"
	"github.com/spf13/cobra"
)

var configFile, migrationsPath string

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
			return migration.Generate(migrationsPath, args[0], "sql", nil, nil)
		}

		if configFile == "" {
			return errors.New("target not provided")
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

		return migration.Generate(path.Path, args[0], "sql", nil, nil)
	},
}

func init() {
	migrationCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yml", "Configuration file path")
	migrationCmd.PersistentFlags().StringVarP(&migrationsPath, "target", "t", "", "Target path where migration will be generated")
}

//assignee = currentUser()
