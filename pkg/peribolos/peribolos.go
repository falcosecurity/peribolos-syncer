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

package peribolos

import (
	"fmt"

	"bitbucket.org/creachadair/stringset"
	peribolos "k8s.io/test-infra/prow/config/org"
)

// New returns a new peribolos.FullConfig structure.
func New() *peribolos.FullConfig {
	return &peribolos.FullConfig{}
}

// UpdateTeamMaintainers updates the maintainers of the specified Team in the specified Organization, adding the maintainers
// list specified as agument.
func UpdateTeamMaintainers(config *peribolos.FullConfig, org, team string, maintainers []string) error {
	orgConfig, ok := config.Orgs[org]
	if !ok {
		return fmt.Errorf("organization not found in Peribolos config")
	}

	teamConfig, ok := orgConfig.Teams[team]
	if !ok {
		return fmt.Errorf("team not fonud in organization %s Peribolos config", org)
	}

	m := teamConfig.Maintainers
	for _, v := range maintainers {
		if !stringset.Contains(m, v) {
			m = append(m, v)
		}
	}

	teamConfig.Maintainers = m
	orgConfig.Teams[team] = teamConfig
	config.Orgs[org] = orgConfig

	return nil
}

// UpdateTeamMembers updates the members of the specified Team in the specified Organization, adding the members list
// specified as argument.
func UpdateTeamMembers(config *peribolos.FullConfig, org, team string, members []string) error {
	orgConfig, ok := config.Orgs[org]
	if !ok {
		return fmt.Errorf("organization not found in Peribolos config")
	}

	teamConfig, ok := orgConfig.Teams[team]
	if !ok {
		return fmt.Errorf("team not fonud in organization %s Peribolos config", org)
	}

	m := teamConfig.Members
	for _, v := range members {
		if !stringset.Contains(m, v) {
			m = append(m, v)
		}
	}

	teamConfig.Members = m
	orgConfig.Teams[team] = teamConfig
	config.Orgs[org] = orgConfig

	return nil
}
