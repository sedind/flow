package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is the current version of the flow binary
const Version = "v0.0.1"

// ProjectFile is name of file which holds flow configuration
const ProjectFile = "flow.yml"

var globalCommands = []string{"init", "version", "help"}

// RootCmd is the hook for all of the other commands in the flow binary
var RootCmd = &cobra.Command{
	SilenceErrors: true,
	Use:           "flow",
	Short:         "Flow simplifies your development flow",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		var isGlobalCmd = false
		for _, globalCmd := range globalCommands {
			if globalCmd == cmd.Name() {
				isGlobalCmd = true
			}
		}
		if isGlobalCmd {
			return nil
		}

		if !isInsideProject() {
			return errors.New("you need to be inside Flow project path to run this command")
		}

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func isInsideProject() bool {
	if _, err := os.Stat(ProjectFile); err != nil {
		return false
	}
	return true
}
