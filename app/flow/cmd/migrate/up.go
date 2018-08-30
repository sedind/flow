package migrate

import (
	"fmt"

	"github.com/spf13/cobra"
)

// UpCmd generates sql migration files
var UpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all of the 'up' migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Apply up migration")
		return nil
	},
}
