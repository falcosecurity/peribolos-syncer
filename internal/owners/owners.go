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
	baseRef = "master"
)

// Handle represents a GitHub user handle.
type Handle string

// OwnersLoadingOptions represents OWNERS loading options.
//
//nolint:revive
type OwnersLoadingOptions struct {

	// RepositoryName represents the git repository name at which load the Owners config.
	RepositoryName string

	// GitRef represents the git reference at which load the Owners config.
	GitRef string

	// ConfigPath represents the path to the Owners config file in the repository.
	ConfigPath string

	// ApproversOnly represents the option to load only the approvers.
	ApproversOnly bool

	// ReviewersOnly represents the option to load only the reviewers.
	ReviewersOnly bool
}

func (o *OwnersLoadingOptions) Validate() error {
	if o.RepositoryName == "" {
		//nolint:goerr113
		return fmt.Errorf("owners file's github repository is empty")
	}

	return nil
}

func (o *OwnersLoadingOptions) AddPFlags(pfs *pflag.FlagSet) {
	pfs.StringVar(&o.RepositoryName, "owners-repository", "", "The name of the github repository from which parse OWNERS file")
	pfs.StringVarP(&o.GitRef, "owners-git-ref", "r", baseRef, "The base Git reference at which parse the OWNERS hierarchy")
	pfs.StringVar(&o.ConfigPath, "owners-config-path", "", "The path to the Owners config file from the root of the Git repository. When specified, they are considered people for which the roles are applied from the root until the specified path.")
	pfs.BoolVar(&o.ApproversOnly, "approvers-only", false, "Whether to load only the approvers from the Owners config")
	pfs.BoolVar(&o.ReviewersOnly, "reviewers-only", false, "Whether to load only the reviewers from the Owners config")
}

// NewClient returns a new repoowners.Client from a prow/github.Client and prow/git/v2.ClientFactory.
// It wraps around repoowners.NewClient facilitating the client build configuration with default settings.
// It possibly returns an error.
func NewClient(githubClient github.Client,
	gitClientFactory gitv2.ClientFactory) *repoowners.Client {
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

	ownersClient := repoowners.NewClient(gitClientFactory, githubClient, mdYAMLEnabled,
		skipCollaborators, ownersDirDenylist, resolver)

	return ownersClient
}
