package generate

import (
	"github.com/spf13/cobra"
)

var configFile, migrationsPath string

// Bind package commands to parent command
func Bind(parentCmd *cobra.Command) {
	parentCmd.AddCommand(migrationCmd)

	parentCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yml", "Configuration file path")
	parentCmd.PersistentFlags().StringVarP(&migrationsPath, "target", "t", "", "Target path where migration will be generated")
}
