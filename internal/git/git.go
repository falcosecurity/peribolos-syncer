package git

import (
	gitobject "github.com/go-git/go-git/v5/plumbing/object"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"

	"github.com/go-git/go-git/v5"
	gitplumbing "github.com/go-git/go-git/v5/plumbing"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func NewEphemeralGitBranch(repo *git.Repository, worktree *git.Worktree) (string, error) {
	id := uuid.New()
	refName := id.String()
	headRef, err := repo.Head()
	if err != nil {
		return "", errors.Wrap(err, "error getting repository HEAD reference")
	}

	ref := gitplumbing.NewHashReference(
		gitplumbing.NewBranchReferenceName(refName),
		headRef.Hash(),
	)
	if err = repo.Storer.SetReference(ref); err != nil {
		return "", errors.Wrap(err, "error setting head git reference")
	}

	if err = worktree.Checkout(&git.CheckoutOptions{
		Branch: gitplumbing.NewBranchReferenceName(refName),
	}); err != nil {
		return "", errors.Wrap(err, "error checking out just created branch")
	}

	return refName, nil
}

func StageAndCommit(repo *git.Repository, worktree *git.Worktree, author *gitobject.Signature,
	pgpEntity *openpgp.Entity, stagePath, commitMessage string) error {
	if worktree == nil {
		return errors.New("worktree cannot be empty")
	}

	// Stage the updated orgs config file to the index.
	if _, err := worktree.Add(stagePath); err != nil {
		return errors.Wrap(err, "error staging orgs config file")
	}

	if author == nil {
		return errors.New("git author cannot be empty")
	}

	commit, err := worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &gitobject.Signature{
			Name:  author.Name,
			Email: author.Email,
			When:  time.Now(),
		},
		SignKey: pgpEntity,
	})
	if err != nil {
		return errors.Wrap(err, "error creating orgs config update git commitAll")
	}

	// Create a commit.
	_, err = repo.CommitObject(commit)
	if err != nil {
		return errors.Wrap(err, "error retrieving the orgs config update commitAll hash")
	}

	return nil
}
