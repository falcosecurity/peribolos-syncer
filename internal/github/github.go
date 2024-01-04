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

package github

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"k8s.io/test-infra/pkg/flagutil"
	prowflags "k8s.io/test-infra/prow/flagutil"
	gitv2 "k8s.io/test-infra/prow/git/v2"
	prowgithub "k8s.io/test-infra/prow/github"
)

// GitHubOptions represents options to interact with GitHub.
//
//nolint:revive
type GitHubOptions struct {
	Username string

	DryRun bool

	prowflags.GitHubOptions
}

func (o *GitHubOptions) AddPFlags(pfs *pflag.FlagSet) {
	pfs.BoolVar(&o.DryRun, "dry-run", false, "Dry run for testing. Uses API tokens but does not mutate.")
	pfs.StringVar(&o.Username, "github-username", "", "The GitHub username")

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	for _, group := range []flagutil.OptionGroup{
		o,
	} {
		group.AddFlags(fs)
	}

	pfs.AddGoFlagSet(fs)
}

func (o *GitHubOptions) ValidateAll() error {
	if o.Username == "" {
		//nolint:goerr113
		return fmt.Errorf("github Username is empty")
	}

	return nil
}

func (o *GitHubOptions) GetGitClientFactory() (gitv2.ClientFactory, error) {
	s := ""

	factory, err := o.GitClientFactory("", &s, o.DryRun)
	if err != nil {
		return nil, errors.Wrap(err, "error creating git client")
	}

	return factory, nil
}

func (o *GitHubOptions) ForkRepository(githubClient prowgithub.Client, githubOrg, githubRepo, token string) (*git.Repository, *git.Worktree, string, error) {
	githubClient.Used()

	path, err := os.MkdirTemp("", "orgs")
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "error creating temporary directory for cloning git repository")
	}

	fork, err := githubClient.EnsureFork(o.Username, githubOrg, githubRepo)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "error creating a fork of the orgs config repository")
	}

	configRepoURL, err := url.JoinPath(fmt.Sprintf("https://%s", o.Host), o.Username, fork)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "error generating orgs config repository URL")
	}

	repository, err := git.PlainClone(path, false, &git.CloneOptions{
		Auth: &githttp.BasicAuth{
			Username: o.Username,
			Password: token,
		},
		URL:      configRepoURL,
		Progress: nil,
	})
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "error cloning git repository")
	}

	worktree, err := repository.Worktree()
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "error getting repository worktree")
	}

	return repository, worktree, path, nil
}
