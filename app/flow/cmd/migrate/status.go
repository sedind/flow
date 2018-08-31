package migrate

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd generates sql migration files
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays the status of all migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(configFile)
		fmt.Println(migrationsPath)
		return nil
	},
}
