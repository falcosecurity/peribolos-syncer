// Copyright 2023 The Falco Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
