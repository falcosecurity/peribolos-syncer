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

package sync

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/maxgio92/peribolos-owners-syncer/pkg/owners"
	"github.com/maxgio92/peribolos-owners-syncer/pkg/peribolos"
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
		Use:   "sync",
		Short: "Synchronize Peribolos org,yaml file from OWNERS file",
	}

	cmd.RunE = o.Run

	cmd.Flags().StringVarP(&o.ownersFilepath, flagOwnersFilePath, "o", defaultOwnersFilepath, "The path to the OWNERS file")
	cmd.Flags().StringVarP(&o.peribolosConfigFilepath, flagPeribolosConfigFilepath, "c", defaultPeribolosConfigFilepath, "The path to the Peribolos org.yaml file")
	cmd.Flags().StringVar(&o.orgName, "org", "", "The name of the Github organization to update")
	cmd.Flags().StringVar(&o.teamName, "team", "", "The name of the Github organization to update")

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
	if err := o.validate(); err != nil {
		return errors.Wrap(err, "error validating parameters")
	}

	f, err := os.ReadFile(o.ownersFilepath)
	if err != nil {
		return errors.Wrap(err, "error reading OWNERS file")
	}

	ow := owners.New()

	if err = yaml.Unmarshal(f, ow); err != nil {
		return errors.Wrap(err, "error unmarshaling owners")
	}

	f, err = os.ReadFile(o.peribolosConfigFilepath)
	if err != nil {
		return errors.Wrap(err, "error reading Peribolos config file")
	}

	orgs := peribolos.New()

	if err = yaml.Unmarshal(f, orgs); err != nil {
		return errors.Wrap(err, "error unmarshaling Peribolos config")
	}

	if err = peribolos.UpdateTeamMembers(orgs, o.orgName, o.teamName, ow.Approvers); err != nil {
		return errors.Wrap(err, "error updating Peribolos' maintainers from OWNERS's approvers")
	}

	compiled, err := yaml.Marshal(orgs)
	if err != nil {
		return errors.Wrap(err, "error recompiling the Peribolos config")
	}

	if err = os.WriteFile(o.peribolosConfigFilepath, compiled, FilePerm); err != nil {
		return errors.Wrap(err, "error writing the recompiled Peribolos config")
	}

	fmt.Println("The Peribolos configuration has been updated.")

	return nil
}
