package migrate

import (
	"github.com/spf13/cobra"
)

var configFile string

// Bind package commands to parent command
func Bind(parentCmd *cobra.Command) {
	parentCmd.AddCommand(upCmd)
	parentCmd.AddCommand(downCmd)
	parentCmd.AddCommand(resetCmd)
	parentCmd.AddCommand(statusCmd)

	parentCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yml", "Configuration file path")
}
