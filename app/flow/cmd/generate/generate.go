package generate

import (
	"github.com/spf13/cobra"
)

// Bind package commands to parent command
func Bind(parentCmd *cobra.Command) {
	parentCmd.AddCommand(migrationCmd)
}
