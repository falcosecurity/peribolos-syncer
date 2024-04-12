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
	"fmt"
	"os"
	"path"
	"strings"
	"unicode"

	"github.com/go-git/go-git/v5"
	gitobject "github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	peribolos "k8s.io/test-infra/prow/config/org"
	"k8s.io/test-infra/prow/github"
	"k8s.io/test-infra/prow/repoowners"
	"sigs.k8s.io/yaml"

	syncergit "github.com/falcosecurity/peribolos-syncer/internal/git"
	syncergithub "github.com/falcosecurity/peribolos-syncer/internal/github"
	"github.com/falcosecurity/peribolos-syncer/internal/output"
	"github.com/falcosecurity/peribolos-syncer/internal/owners"
	"github.com/falcosecurity/peribolos-syncer/internal/sync"
	orgs "github.com/falcosecurity/peribolos-syncer/pkg/peribolos"
	"github.com/falcosecurity/peribolos-syncer/pkg/pgp"
)

type options struct {
	*sync.CommonOptions

	author            gitobject.Signature
	privateGPGKeyPath string
	publicGPGKeyPath  string

	github syncergithub.GitHubOptions
	orgs   *orgs.Options
	owners *owners.OwnersLoadingOptions

	git.ListOptions
}

// New returns a new sync github command.
func New() *cobra.Command {
	o := &options{
		CommonOptions: &sync.CommonOptions{},
		author:        gitobject.Signature{},
		github:        syncergithub.GitHubOptions{},
		owners:        &owners.OwnersLoadingOptions{},
		orgs:          &orgs.Options{},
	}

	cmd := &cobra.Command{
		Use:     commandName,
		Short:   commandShortDescription,
		Example: commandExample,
		RunE:    o.Run,
	}

	// Organization sync options.
	cmd.Flags().StringVar(&o.GitHubOrg, "org", "", "The name of the GitHub organization to update configuration for")
	cmd.Flags().StringVar(&o.GitHubTeam, "team", "", "The name of the GitHub team to update configuration for")

	// Git author options.
	cmd.Flags().StringVar(&o.author.Name, "git-author-name", "", "The Git author name with which write commits for the update of the Peribolos config")
	cmd.Flags().StringVar(&o.author.Email, "git-author-email", "", "The Git author email with which write commits for the update of the Peribolos config")
	cmd.Flags().StringVar(&o.publicGPGKeyPath, "gpg-public-key", "", "The path to the public GPG key for signing git commits")
	cmd.Flags().StringVar(&o.privateGPGKeyPath, "gpg-private-key", "", "The path to the private GPG key for signing git commits")

	// GitHub options.
	o.github.AddPFlags(cmd.Flags())

	// Owners options.
	o.owners.AddPFlags(cmd.Flags())

	// Orgs config options.
	o.orgs.AddPFlags(cmd.Flags())

	return cmd
}

func (o *options) validate() error {
	if o.GitHubOrg == "" {
		return errors.New("github organization name is empty")
	}

	if o.GitHubTeam == "" {
		return errors.New("github team name is empty")
	}

	if o.author.Name == "" {
		return errors.New("git author name is empty")
	}

	if o.author.Email == "" {
		return errors.New("git author email is empty")
	}

	if o.publicGPGKeyPath == "" {
		return errors.New("git author public pgp key path cannot be empty")
	}

	if o.privateGPGKeyPath == "" {
		return errors.New("git author private pgp key path cannot be empty")
	}

	if err := o.owners.Validate(); err != nil {
		return err
	}

	if err := o.orgs.Validate(); err != nil {
		return err
	}

	if err := o.github.ValidateAll(); err != nil {
		return err
	}

	return nil
}

//nolint:funlen
func (o *options) Run(_ *cobra.Command, _ []string) error {
	if err := o.validate(); err != nil {
		return err
	}

	token, err := getTokenFromFile(o.github.TokenPath)
	if err != nil {
		return errors.Wrap(err, "error reading token from file")
	}

	// Build GitHub client.
	githubClient, err := o.github.GitHubClientWithAccessToken(token)
	if err != nil {
		return errors.Wrap(err, "error generating github client with specified access token")
	}

	// Load Owners hierarchy from specified repository.
	owners, err := o.loadOwnersFromGithub(githubClient)
	if err != nil {
		return err
	}

	// Load specified people from the Owners structure.
	people := o.loadPeopleFromOwners(owners)

	// Clone the peribolos config repository.
	repo, worktree, local, err := o.github.ForkRepository(
		githubClient, o.GitHubOrg, o.orgs.ConfigRepo, token)
	if err != nil {
		return errors.Wrap(err, "error forking the config repository")
	}

	// Load GitHub orgs config from the git working tree's filesystem.
	config, err := orgs.LoadConfigFromFilesystem(worktree.Filesystem, o.orgs.ConfigPath)
	if err != nil {
		return errors.Wrap(err, "error loading the config")
	}

	// Create an ephemeral branch for the changes.
	ref, err := syncergit.NewEphemeralGitBranch(repo, worktree)
	if err != nil {
		return errors.Wrap(err, "error creating new branch for changes on the config")
	}

	// Synchronize the Github Team config with Approvers.
	if err = orgs.AddTeamMembers(config, o.GitHubOrg, o.GitHubTeam, people); err != nil {
		return errors.Wrap(err, "error updating maintainers github team from leaf approvers")
	}

	// Flush updated config to local working copy.
	if err = o.flushConfig(config, local); err != nil {
		return errors.Wrap(err, "error writing updated peribolos config")
	}

	// Store the change in a commitAll with a log.
	commitMsg := fmt.Sprintf(`chore(%s): update %s team members

The update reflects the content of the related repository's OWNERS tree.
%s

Signed-off-by: %s <%s>
`, peribolosConfigFile, o.GitHubTeam, syncerSignature, o.author.Name, o.author.Email)

	// Generate a PGP entity to sign the git commits.
	pgpEntity, err := pgp.NewPGPEntity(o.author.Name, o.author.Email, o.publicGPGKeyPath, o.privateGPGKeyPath)
	if err != nil {
		return errors.Wrap(err, "error generating the pgp entity")
	}

	// Stage the change to the config and create a commit for it.
	if err = syncergit.StageAndCommit(repo, worktree, &o.author, pgpEntity, o.orgs.ConfigPath, commitMsg); err != nil {
		return errors.Wrap(err, "error committing the changes on config")
	}

	// Skip push to remote and pull request creation when dry run.
	if !o.github.DryRun {
		// Push the new branch to the remote.
		if err := repo.Push(&git.PushOptions{
			Auth: &githttp.BasicAuth{
				Username: o.github.Username,
				Password: token,
			},
		}); err != nil {
			return errors.Wrap(err, "error pushing config update git branch")
		}

		// Create a Pull Request on GitHub.
		pr, err := githubClient.CreatePullRequest(
			o.GitHubOrg,
			o.orgs.ConfigRepo,
			fmt.Sprintf("Sync Github Team %s with %s owners", o.GitHubTeam, o.owners.RepositoryName),
			fmt.Sprintf(`This PR synchronizes the Github Team %s with the leaf approvers declared in %s repository's [OWNERS](%s) file.

%s
`, o.GitHubTeam, o.owners.RepositoryName, ownersDoc, syncerSignature),
			fmt.Sprintf("%s:%s", o.github.Username, ref),
			o.orgs.ConfigBaseRef,
			false,
		)
		if err != nil {
			return errors.Wrap(err, "error creating github pull request")
		}

		output.Print(
			fmt.Sprintf("A Pull Request has been opened: https://%s/%s/%s/pull/%d",
				o.github.Host, o.GitHubOrg, o.orgs.ConfigRepo, pr),
		)

		return nil
	}

	output.Print("Skipping pull request.")

	return nil
}

func (o *options) loadOwnersFromGithub(githubClient github.Client) (repoowners.RepoOwner, error) {
	gitClientFactory, err := o.github.GetGitClientFactory()
	if err != nil {
		return nil, errors.Wrap(err, "error building git client gitclientfactory")
	}

	ownersClient := owners.NewClient(githubClient, gitClientFactory)

	// Load Owners hierarchy from specified repository.
	owners, err := ownersClient.LoadRepoOwners(o.GitHubOrg, o.owners.RepositoryName, o.owners.GitRef)
	if err != nil {
		return nil, errors.Wrap(err, "error loading owners from repository")
	}

	return owners, nil
}

func (o *options) loadPeopleFromOwners(owners repoowners.RepoOwner) []string {
	var people []string

	switch {
	// Limiting the scope of the roles.
	case o.owners.ConfigPath != "":
		//nolint:gocritic
		if o.owners.ApproversOnly {
			// Approvers of the subpart of the repository.
			people = maps.Keys(owners.Approvers(o.owners.ConfigPath).Set())
		} else if o.owners.ReviewersOnly {
			// Reviewers of the subpart of the repository.
			people = maps.Keys(owners.Reviewers(o.owners.ConfigPath).Set())
		} else {
			// Both approvers and reviewers of the subpart of the repository.
			people = maps.Keys(owners.Approvers(o.owners.ConfigPath).
				Union(owners.Reviewers(o.owners.ConfigPath)).Set())
		}
	// Approvers of the whole repository.
	case o.owners.ApproversOnly:
		people = maps.Keys(owners.AllApprovers())
	// Reviewers of the whole repository.
	case o.owners.ReviewersOnly:
		people = maps.Keys(owners.AllReviewers())
	// Both approvers and reviewers of the whole repository.
	default:
		people = maps.Keys(owners.AllOwners())
	}

	return people
}

func (o *options) flushConfig(config *peribolos.FullConfig, configPath string) error {
	b, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "error recompiling the peribolos config")
	}

	if err = os.WriteFile(path.Join(configPath, o.orgs.ConfigPath), b, modeConfigFile); err != nil {
		return errors.Wrap(err, "error writing the recompiled peribolos config")
	}

	return nil
}

func getTokenFromFile(path string) (string, error) {
	token, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "error reading token file")
	}

	return removeNonPrintableChars(string(token)), nil
}

func removeNonPrintableChars(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case unicode.IsPrint(r):
			return r
		default:
			return -1
		}
	}, s)
}
