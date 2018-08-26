package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Flow toolkit",
	Run: func(c *cobra.Command, args []string) {
		fmt.Printf("Flow toolkit version is: %s\n", Version)
	},
	// override the roon level pre-run func
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		return nil
	},
}
