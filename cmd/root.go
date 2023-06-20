package cmd

import (
	"github.com/spf13/cobra"

	"github.com/falcosecurity/peribolos-syncer/cmd/sync"
	"github.com/falcosecurity/peribolos-syncer/cmd/version"
	"github.com/falcosecurity/peribolos-syncer/internal/output"
)

// New returns a new root command.
func New() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = CommandName
	cmd.Long = CommandLongDescription
	cmd.DisableAutoGenTag = true

	// Add subcommands.
	cmd.AddCommand(sync.New())
	cmd.AddCommand(version.New())

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd := New()
	output.ExitOnErr(cmd.Execute())
}
