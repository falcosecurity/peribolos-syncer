/*
Copyright © 2023 maxgio92 me@maxgio.me

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package github

import (
	"fmt"
	"github.com/spf13/cobra"
)

type options struct {
	peribolosConfigFilepath string
	ownersFilepath          string
	orgName                 string
	teamName                string
}

// New returns a new root command.
func New() *cobra.Command {
	o := &options{}

	cmd := &cobra.Command{
		Use:   "github",
		Short: "Synchronize Peribolos org.yaml file from OWNERS file on remote Github repositories via Pull Request",
	}

	cmd.RunE = o.Run

	return cmd
}

func (o *options) validate() error {
	if o.orgName == "" {
		return fmt.Errorf("org name is empty")
	}

	if o.teamName == "" {
		return fmt.Errorf("team name is empty")
	}

	return nil
}

func (o *options) Run(cmd *cobra.Command, agrs []string) error {
	return nil
}
