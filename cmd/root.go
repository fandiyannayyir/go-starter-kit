// Package cmd wires the cobra command tree. Subcommands are registered
// explicitly in Execute (no init() magic), so the command set is traceable.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-starter-kit",
	Short: "go-starter-kit service",
}

// Execute registers the subcommands and runs the root command.
func Execute() {
	rootCmd.AddCommand(serverCmd, migrateCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
