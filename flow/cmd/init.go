package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Flow project",
	Run: func(c *cobra.Command, args []string) {
		fmt.Println("init")
	},
	// override the roon level pre-run func
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		return nil
	},
}
