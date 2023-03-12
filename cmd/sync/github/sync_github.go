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
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"io"
	"net/url"
	"os"
	"path"
	"time"

	gogit "github.com/go-git/go-git/v5"
	gogitplumbing "github.com/go-git/go-git/v5/plumbing"
	gogitobject "github.com/go-git/go-git/v5/plumbing/object"
	gogithttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	syncergithub "github.com/maxgio92/peribolos-owners-syncer/internal/github"
	"github.com/maxgio92/peribolos-owners-syncer/pkg/owners"
	"github.com/maxgio92/peribolos-owners-syncer/pkg/peribolos"
	"sigs.k8s.io/yaml"
)

type options struct {
	teamName string
	orgName  string

	gitAuthorEmail string
	gitAuthorName  string
	gitBaseBranch  string

	githubUsername string
	github         syncergithub.GitHubOptions

	peribolos *peribolos.PeribolosOptions
	owners    *owners.OwnersOptions
}

const (
	peribolosRepoBaseReference = "master"
	syncerSignature            = "Autogenerated with [peribolos-owners-syncer](https://github.com/maxgio92/peribolos-owners-syncer)."
)

// New returns a new sync github command.
// TODO: generate doc.
func New() *cobra.Command {
	o := &options{
		owners:    &owners.OwnersOptions{},
		peribolos: &peribolos.PeribolosOptions{},
	}

	cmd := &cobra.Command{
		Use:   "github",
		Short: "Synchronize Peribolos org.yaml file from OWNERS file on remote github repositories via Pull Request",
		RunE:  o.Run,
	}

	cmd.Flags().StringVar(&o.orgName, "org", "", "The name of the GitHub organization to update configuration for")
	cmd.Flags().StringVar(&o.teamName, "team", "", "The name of the GitHub team to update configuration for")

	cmd.Flags().StringVar(&o.gitAuthorName, "git-author-name", "", "The Git author name with which write commits for the update of the Peribolos config")
	cmd.Flags().StringVar(&o.gitAuthorEmail, "git-author-email", "", "The Git author email with which write commits for the update of the Peribolos config")
	cmd.Flags().StringVar(&o.gitBaseBranch, "git-base-branch", peribolosRepoBaseReference, "The Git base branch of the Peribolos config repository")
	cmd.Flags().StringVar(&o.githubUsername, "github-username", "", "The GitHub username with which authenticate")

	o.github.AddPFlags(cmd.Flags())
	o.owners.AddPFlags(cmd.Flags())
	o.peribolos.AddPFlags(cmd.Flags())

	return cmd
}

func (o *options) validate() error {
	if o.orgName == "" {
		return fmt.Errorf("github organization name is empty")
	}

	if o.teamName == "" {
		return fmt.Errorf("github team name is empty")
	}

	if o.gitAuthorName == "" {
		return fmt.Errorf("git author name is empty")
	}

	if o.gitAuthorEmail == "" {
		return fmt.Errorf("git author email is empty")
	}

	if o.githubUsername == "" {
		return fmt.Errorf("github username is empty")
	}

	if err := o.owners.Validate(); err != nil {
		return err
	}

	if err := o.peribolos.Validate(); err != nil {
		return err
	}

	return nil
}

func (o *options) Run(cmd *cobra.Command, agrs []string) error {
	if err := o.validate(); err != nil {
		return err
	}

	githubClient, err := o.github.BuildClient()
	if err != nil {
		return errors.Wrap(err, "error building github client")
	}

	gitClientFactory, err := o.github.GetGitClientFactory()
	if err != nil {
		return errors.Wrap(err, "error building git client factory")
	}

	ownersClient, err := o.owners.BuildClient(githubClient, gitClientFactory)
	if err != nil {
		return errors.Wrap(err, "error building owners client")
	}

	owners, err := ownersClient.LoadRepoOwners(o.orgName, o.owners.OwnersRepo, o.owners.OwnersRef)
	if err != nil {
		return errors.Wrap(err, "error loading owners")
	}

	token, err := os.ReadFile(o.github.TokenPath)
	if err != nil {
		return errors.Wrap(err, "error reading token file")
	}

	// TODO: fork the Peribolos config repository.

	// Clone Peribolos config repository.
	tmp, err := os.MkdirTemp("", "peribolos")
	if err != nil {
		return errors.Wrap(err, "error creating temporary direcotry for cloning git repo")
	}

	peribolosRepoURL, err := url.JoinPath("https://github.com", o.orgName, o.peribolos.ConfigRepo)
	if err != nil {
		return errors.Wrap(err, "error generating Peribolos config repository URL")
	}

	repo, err := gogit.PlainClone(tmp, false, &gogit.CloneOptions{
		Auth: &gogithttp.BasicAuth{
			Username: o.githubUsername, // yes, this can be anything except an empty string
			Password: string(token),
		},
		URL:      peribolosRepoURL,
		Progress: nil,
	})

	if err != nil {
		return errors.Wrap(err, "error cloning git repository")
	}

	// Get reposiotry HEAD symbolic reference
	headRef, err := repo.Head()

	id := uuid.New()
	branchRef := id.String()
	if err != nil {
		return errors.Wrap(err, "error getting repository HEAD reference")
	}

	ref := gogitplumbing.NewHashReference(
		gogitplumbing.NewBranchReferenceName(branchRef),
		headRef.Hash(),
	)
	err = repo.Storer.SetReference(ref)

	worktree, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "error getting repository worktree")
	}

	if err = worktree.Checkout(&gogit.CheckoutOptions{
		Branch: gogitplumbing.NewBranchReferenceName(branchRef),
	}); err != nil {
		return errors.Wrap(err, "error checking out just created branch")
	}

	approvers := maps.Keys(owners.LeafApprovers(o.owners.OwnersPath))

	orgs := peribolos.New()

	// Open Peribolos config file.
	file, err := worktree.Filesystem.Open(o.peribolos.ConfigPath)
	if err != nil {
		return errors.Wrap(err, "error opening Peribolos config file")
	}

	// Build Peribolos config.
	b, err := io.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "error reading Peribolos config file")
	}

	if err = yaml.Unmarshal(b, orgs); err != nil {
		return errors.Wrap(err, "error unmarshaling Peribolos config")
	}

	// Synchronize the Github Team config with Approvers.
	if err = peribolos.UpdateTeamMembers(orgs, o.orgName, o.teamName, approvers); err != nil {
		return errors.Wrap(err, "error updating Peribolos' maintainers from OWNERS's approvers")
	}

	// Write the update Peribolos config file.
	b, err = yaml.Marshal(orgs)
	if err != nil {
		return errors.Wrap(err, "error recompiling the Peribolos config")
	}

	if err = os.WriteFile(path.Join(tmp, o.peribolos.ConfigPath), b, 0644); err != nil {
		return errors.Wrap(err, "error writing the recompiled Peribolos config")
	}

	// Stage the updated Peribolos config file to the index.
	if _, err = worktree.Add(o.peribolos.ConfigPath); err != nil {
		return errors.Wrap(err, "error staging Peribolos config file")
	}

	// Store the change in a commit with a log.
	msg := fmt.Sprintf(`chore(org.yaml): update %s team members

The update reflects the content of the related repository root's OWNERS.
%s

Signed-off-by: %s <%s>
`, o.teamName, syncerSignature, o.gitAuthorName, o.gitAuthorEmail)
	commit, err := worktree.Commit(msg, &gogit.CommitOptions{
		Author: &gogitobject.Signature{
			Name:  o.gitAuthorName,
			Email: o.gitAuthorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return errors.Wrap(err, "error creating Peribolos config update git commit")
	}

	_, err = repo.CommitObject(commit)
	if err != nil {
		return errors.Wrap(err, "error retrieving the Peribolos config update commit hash")
	}

	if err = repo.Push(&gogit.PushOptions{
		Auth: &gogithttp.BasicAuth{
			Username: o.githubUsername, // yes, this can be anything except an empty string
			Password: string(token),
		},
	}); err != nil {
		return errors.Wrap(err, "error pushing Peribolos config update git branch")
	}

	// Create a Github Client with access token authentication.
	gh, err := o.github.GitHubClientWithAccessToken(string(token))
	if err != nil {
		return errors.Wrap(err, "error generating GitHub client with access token")
	}

	// Create a Pull Request with the change to the Peribolos config.
	pr, err := gh.CreatePullRequest(
		o.orgName,
		o.peribolos.ConfigRepo,
		fmt.Sprintf("Sync Github Team %s with %s owners", o.teamName, o.owners.OwnersRepo),
		fmt.Sprintf(`This PR synchronizes the Github Team %s with the root approvers in %s repository's OWNERS'.

%s
`, o.teamName, o.owners.OwnersRepo, syncerSignature),
		branchRef,
		o.gitBaseBranch,
		false,
	)
	if err != nil {
		return errors.Wrap(err, "error creating GitHub Pull Request")
	}

	fmt.Printf("Opened Pull request %d", pr)

	return nil
}
