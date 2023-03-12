package sync

import (
	"github.com/maxgio92/peribolos-owners-syncer/cmd/sync/github"
	"github.com/spf13/cobra"

	"github.com/maxgio92/peribolos-owners-syncer/cmd/sync/local"
)

// New returns a new root command.
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize Peribolos org,yaml file from OWNERS file",
	}

	// Add sync subcommand.
	cmd.AddCommand(local.New())
	cmd.AddCommand(github.New())

	return cmd
}
