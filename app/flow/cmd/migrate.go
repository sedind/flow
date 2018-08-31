package cmd

import (
	"github.com/sedind/flow/app/flow/cmd/migrate"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Aliases: []string{"m"},
	Short:   "Tools for working with your database migrations.",
}

func init() {
	migrate.Bind(migrateCmd)
	RootCmd.AddCommand(migrateCmd)
}
