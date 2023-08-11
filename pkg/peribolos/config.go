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

package peribolos

import (
	"fmt"
	"io"

	"bitbucket.org/creachadair/stringset"
	"github.com/go-git/go-billy/v5"
	"github.com/pkg/errors"
	peribolos "k8s.io/test-infra/prow/config/org"
	"sigs.k8s.io/yaml"
)

// NewConfig returns a new orgs.FullConfig structure.
func NewConfig() *peribolos.FullConfig {
	return &peribolos.FullConfig{Orgs: map[string]peribolos.Config{}}
}

// LoadConfigFromFilesystem loads the peribolos config from the filesystem.
// It possibly returns an error.
func LoadConfigFromFilesystem(fs billy.Filesystem, configPath string) (*peribolos.FullConfig, error) {
	r, err := fs.Open(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "error opening peribolos config file")
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "error reading peribolos config file")
	}

	config := NewConfig()

	if err = yaml.Unmarshal(b, config); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling peribolos config")
	}

	return config, nil
}

// AddTeamMaintainers updates the maintainers of the specified Team in the specified Organization, adding the maintainers
// list specified as agument.
func AddTeamMaintainers(config *peribolos.FullConfig, org, team string, maintainers []string) error {
	orgConfig, ok := config.Orgs[org]
	if !ok {
		//nolint:goerr113
		return fmt.Errorf("organization not found in peribolos config")
	}

	teamConfig, ok := orgConfig.Teams[team]
	if !ok {
		//nolint:goerr113
		return fmt.Errorf("team not fonud in organization %s peribolos config", org)
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

// AddTeamMembers updates the members of the specified Team in the specified Organization, adding the members list
// specified as argument.
func AddTeamMembers(config *peribolos.FullConfig, org, team string, members []string) error {
	orgConfig, ok := config.Orgs[org]
	if !ok {
		//nolint:goerr113
		return errors.New("organization not found in peribolos config")
	}

	teamConfig, ok := orgConfig.Teams[team]
	if !ok {
		//nolint:goerr113
		return fmt.Errorf("team not found in organization %s peribolos config", org)
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
