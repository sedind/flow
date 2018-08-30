package generate

import (
	"errors"

	"github.com/spf13/cobra"
)

// MigrationCmd generates sql migration files
var MigrationCmd = &cobra.Command{
	Use:   "migration [name]",
	Short: "Generates Up/Down migrations for your database using SQL.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must supply a name for your migration")
		}
		return nil
	},
}
