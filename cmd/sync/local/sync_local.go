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

package local

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/test-infra/prow/repoowners"
	"sigs.k8s.io/yaml"

	"github.com/maxgio92/peribolos-owners-syncer/internal/output"
	"github.com/maxgio92/peribolos-owners-syncer/pkg/orgs"
	"github.com/maxgio92/peribolos-owners-syncer/pkg/sync"
)

type options struct {
	peribolosConfigFilepath string
	ownersFilepath          string

	*sync.Options
}

// New returns a new sync local command.
func New() *cobra.Command {
	o := &options{
		Options: &sync.Options{},
	}

	cmd := &cobra.Command{
		Use:   "local",
		Short: "Synchronize Peribolos org.yaml file from OWNERS file in the local filesystem",
	}

	cmd.RunE = o.Run

	cmd.Flags().StringVarP(&o.ownersFilepath, flagOwnersFilePath, "o", defaultOwnersFilepath, "The path to the OWNERS file")
	cmd.Flags().StringVarP(&o.peribolosConfigFilepath, flagPeribolosConfigFilepath, "c", defaultPeribolosConfigFilepath, "The path to the Peribolos org.yaml file")
	cmd.Flags().StringVar(&o.GitHubOrg, "org", "", "The name of the GitHub organization to update")
	cmd.Flags().StringVar(&o.GitHubTeam, "team", "", "The name of the GitHub organization to update")

	return cmd
}

func (o *options) validate() error {
	if o.GitHubOrg == "" {
		//nolint:goerr113
		return fmt.Errorf("org name is empty")
	}

	if o.GitHubTeam == "" {
		//nolint:goerr113
		return fmt.Errorf("team name is empty")
	}

	return nil
}

func (o *options) Run(cmd *cobra.Command, agrs []string) error {
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "error validating parameters")
	}

	b, err := os.ReadFile(o.ownersFilepath)
	if err != nil {
		return errors.Wrap(err, "error reading OWNERS file")
	}

	owners, err := repoowners.LoadSimpleConfig(b)
	if err != nil {
		return errors.Wrap(err, "error unmarshaling owners")
	}

	b, err = os.ReadFile(o.peribolosConfigFilepath)
	if err != nil {
		return errors.Wrap(err, "error reading Peribolos config file")
	}

	orgsConfig := orgs.NewConfig()

	if err = yaml.Unmarshal(b, orgsConfig); err != nil {
		return errors.Wrap(err, "error unmarshaling Peribolos config")
	}

	if err = orgs.UpdateTeamMembers(orgsConfig, o.GitHubOrg, o.GitHubTeam, owners.Approvers); err != nil {
		return errors.Wrap(err, "error updating Peribolos' maintainers from OWNERS's approvers")
	}

	compiled, err := yaml.Marshal(orgsConfig)
	if err != nil {
		return errors.Wrap(err, "error recompiling the Peribolos config")
	}

	if err = os.WriteFile(o.peribolosConfigFilepath, compiled, FilePerm); err != nil {
		return errors.Wrap(err, "error writing the recompiled Peribolos config")
	}

	output.Print("The Peribolos configuration has been updated.")

	return nil
}
