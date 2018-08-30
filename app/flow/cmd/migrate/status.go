package migrate

import (
	"fmt"

	"github.com/spf13/cobra"
)

// StatusCmd generates sql migration files
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays the status of all migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("show migration status")
		return nil
	},
}
