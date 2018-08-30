package generate

import (
	"errors"

	"github.com/sedind/flow/app/dbe/migration"
	"github.com/spf13/cobra"
)

var configFile, migrationsPath string

// MigrationCmd generates sql migration files
var MigrationCmd = &cobra.Command{
	Use:   "migration [name]",
	Short: "Generates Up/Down migration files for your SQL database.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must supply a name for your migration")
		}
		return migration.Generate(migrationsPath, args[0], "sql", nil, nil)
	},
}

func init() {
	MigrationCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yml", "Configuration file path")
	MigrationCmd.PersistentFlags().StringVarP(&migrationsPath, "target", "t", "", "Target path where migration will be generated")
}
