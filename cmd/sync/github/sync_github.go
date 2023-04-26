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
	"os"
	"path"

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

	syncergit "github.com/maxgio92/peribolos-syncer/internal/git"
	"github.com/maxgio92/peribolos-syncer/internal/output"
	syncergithub "github.com/maxgio92/peribolos-syncer/pkg/github"
	"github.com/maxgio92/peribolos-syncer/pkg/owners"
	orgs "github.com/maxgio92/peribolos-syncer/pkg/peribolos"
	"github.com/maxgio92/peribolos-syncer/pkg/pgp"
	"github.com/maxgio92/peribolos-syncer/pkg/sync"
)

type options struct {
	*sync.Options

	author            gitobject.Signature
	privateGPGKeyPath string
	publicGPGKeyPath  string

	github syncergithub.GitHubOptions
	orgs   *orgs.PeribolosOptions
	owners *owners.OwnersOptions
	git.ListOptions
}

// New returns a new sync github command.
func New() *cobra.Command {
	o := &options{
		Options: &sync.Options{},
		author:  gitobject.Signature{},
		github:  syncergithub.GitHubOptions{},
		owners:  &owners.OwnersOptions{},
		orgs:    &orgs.PeribolosOptions{},
	}

	cmd := &cobra.Command{
		Use:   "github",
		Short: "Synchronize Peribolos org.yaml file from OWNERS file on remote github repositories via Pull Request",
		RunE:  o.Run,
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
		return fmt.Errorf("github organization name is empty")
	}

	if o.GitHubTeam == "" {
		return fmt.Errorf("github team name is empty")
	}

	if o.author.Name == "" {
		return fmt.Errorf("git author name is empty")
	}

	if o.author.Email == "" {
		return fmt.Errorf("git author email is empty")
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

func (o *options) Run(_ *cobra.Command, _ []string) error {
	if err := o.validate(); err != nil {
		return err
	}

	// Get the GitHub token from filesystem.
	token, err := os.ReadFile(o.github.TokenPath)
	if err != nil {
		return errors.Wrap(err, "error reading token file")
	}

	// Build GitHub client.
	githubClient, err := o.github.GitHubClientWithAccessToken(string(token))
	if err != nil {
		return errors.Wrap(err, "error generating GitHub client with specified access token")
	}

	// Load Owners hierarchy from specified repository.
	owners, err := o.loadOwnersFromGithub(githubClient, o.GitHubOrg, o.owners.OwnersRepo, o.owners.OwnersGitRef)
	if err != nil {
		return err
	}

	// Get the leaf approvers from the Owners hierarchy.
	approvers := maps.Keys(owners.LeafApprovers(o.owners.OwnersPath))

	// Clone the peribolos config repository.
	repo, worktree, local, err := o.github.ForkRepository(
		githubClient, o.GitHubOrg, o.orgs.ConfigRepo, string(token))
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
	if err := orgs.AddTeamMembers(config, o.GitHubOrg, o.GitHubTeam, approvers); err != nil {
		return errors.Wrap(err, "error updating maintainers github team from leaf approvers")
	}

	// Flush updated config to local working copy.
	if err = o.flushConfig(config, local); err != nil {
		return errors.Wrap(err, "error writing updated peribolos config")
	}

	// Store the change in a commitAll with a log.
	commitMsg := fmt.Sprintf(`chore(%s): update %s team members

The update reflects the content of the related repository root's OWNERS.
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
				Password: string(token),
			},
		}); err != nil {
			return errors.Wrap(err, "error pushing config update git branch")
		}

		// Create a Pull Request on GitHub.
		pr, err := githubClient.CreatePullRequest(
			o.GitHubOrg,
			o.orgs.ConfigRepo,
			fmt.Sprintf("Sync Github Team %s with %s owners", o.GitHubTeam, o.owners.OwnersRepo),
			fmt.Sprintf(`This PR synchronizes the Github Team %s with the leaf approvers declared in %s repository's [OWNERS](%s) file.

%s
`, o.GitHubTeam, o.owners.OwnersRepo, ownersDoc, syncerSignature),
			fmt.Sprintf("%s:%s", o.github.Username, ref),
			o.orgs.ConfigBaseRef,
			false,
		)
		if err != nil {
			return errors.Wrap(err, "error creating GitHub Pull Request")
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

func (o *options) loadOwnersFromGithub(githubClient github.Client, githubOrg, gitRepo, gitRef string) (repoowners.RepoOwner, error) {
	gitClientFactory, err := o.github.GetGitClientFactory()
	if err != nil {
		return nil, errors.Wrap(err, "error building git client gitClientFactory")
	}

	// Load Owners hierarchy from specified repository.
	owners, err := o.owners.LoadFromGitHub(
		githubClient, gitClientFactory,
		githubOrg, gitRepo, gitRef,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error loading owners from repository")
	}

	return owners, nil
}

func (o *options) flushConfig(config *peribolos.FullConfig, configPath string) error {
	b, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "error recompiling the peribolos config")
	}

	if err = os.WriteFile(path.Join(configPath, o.orgs.ConfigPath), b, 0644); err != nil {
		return errors.Wrap(err, "error writing the recompiled peribolos config")
	}

	return nil
}
