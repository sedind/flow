package migrate

import (
	"github.com/spf13/cobra"
)

// Bind package commands to parent command
func Bind(parentCmd *cobra.Command) {
	parentCmd.AddCommand(upCmd)
	parentCmd.AddCommand(downCmd)
	parentCmd.AddCommand(resetCmd)
	parentCmd.AddCommand(statusCmd)
}
