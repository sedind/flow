package cmd

import (
	"github.com/sedind/flow/flow/cmd/generate"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Generates things that we need ",
}

func init() {
	generate.Bind(generateCmd)
	RootCmd.AddCommand(generateCmd)
}
