package migrate

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ResetCmd generates sql migration files
var ResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "The equivalent of running `migrate down` and then `migrate up`",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("reset  migrations")
		return nil
	},
}
