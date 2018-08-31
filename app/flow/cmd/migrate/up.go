package migrate

import (
	"fmt"

	"github.com/spf13/cobra"
)

// upCmd generates sql migration files
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all of the 'up' migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Apply up migration")
		return nil
	},
}
