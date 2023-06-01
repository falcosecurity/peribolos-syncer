package sync

import (
	"github.com/spf13/cobra"

	"github.com/falcosecurity/peribolos-syncer/cmd/sync/github"
	"github.com/falcosecurity/peribolos-syncer/cmd/sync/local"
)

// New returns a new root command.
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   commandName,
		Short: commandShortDescription,
	}

	// Add sync subcommand.
	cmd.AddCommand(local.New())
	cmd.AddCommand(github.New())

	return cmd
}
