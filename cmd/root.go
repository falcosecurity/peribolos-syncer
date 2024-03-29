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
