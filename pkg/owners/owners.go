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

package owners

import (
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/test-infra/prow/config"
	gitv2 "k8s.io/test-infra/prow/git/v2"
	"k8s.io/test-infra/prow/github"
	"k8s.io/test-infra/prow/plugins/ownersconfig"
	"k8s.io/test-infra/prow/repoowners"
)

const (
	baseRef  = "master"
	pathRoot = ""
)

// Handle represents a GitHub user handle.
type Handle string

// Owners represents a GitHub OWNERS flie.
type Owners struct {
	// Approvers is a list of users with ability to approve pull requests.
	Approvers []string `json:"approvers"`

	// EmeritusApprovers is a list of users which served for time as approvers.
	EmeritusApprovers []string `json:"emeritus_approvers"`

	// Reviewers is a list of users with ability to review pull requests.
	Reviewers []string `json:"reviewers"`
}

// OwnersOptions represents OWNERS loading options.
//
//nolint:revive
type OwnersOptions struct {
	OwnersRepo    string
	OwnersBaseRef string
	OwnersPath    string
}

func (o *OwnersOptions) Validate() error {
	if o.OwnersRepo == "" {
		//nolint:goerr113
		return fmt.Errorf("OWNERS file's github repository name is empty")
	}

	return nil
}

func (o *OwnersOptions) AddPFlags(pfs *pflag.FlagSet) {
	pfs.StringVar(&o.OwnersRepo, "owners-repository", "", "The name of the github repository from which parse OWNERS file")
	pfs.StringVarP(&o.OwnersBaseRef, "owners-base-ref", "r", baseRef, "The base Git reference at which parse the OWNERS hierarchy")
	pfs.StringVarP(&o.OwnersPath, "owners-file", "o", pathRoot, "The path to the OWNERS file from the root of the Git repository. Ignored with sync-github.")
}

func (o *OwnersOptions) BuildClient(githubClient github.Client, gitClientFactory gitv2.ClientFactory) (*repoowners.Client, error) {
	mdYAMLEnabled := func(org, repo string) bool {
		return false
	}

	skipCollaborators := func(org, repo string) bool {
		return true
	}

	ownersDirDenylist := func() *config.OwnersDirDenylist {
		return &config.OwnersDirDenylist{}
	}

	resolver := func(org, repo string) ownersconfig.Filenames {
		return ownersconfig.Filenames{
			Owners:        ownersconfig.DefaultOwnersFile,
			OwnersAliases: ownersconfig.DefaultOwnersAliasesFile,
		}
	}

	ownersClient := repoowners.NewClient(gitClientFactory, githubClient, mdYAMLEnabled, skipCollaborators, ownersDirDenylist, resolver)

	return ownersClient, nil
}
