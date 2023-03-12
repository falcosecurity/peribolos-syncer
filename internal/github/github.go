package github

import (
	"flag"
	"fmt"
	"os"

	gogit "github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"k8s.io/test-infra/pkg/flagutil"
	prow "k8s.io/test-infra/prow/flagutil"
	gitv2 "k8s.io/test-infra/prow/git/v2"
	"k8s.io/test-infra/prow/github"
)

type GitHubOptions struct {
	prow.GitHubOptions

	DryRun bool
}

func (o *GitHubOptions) AddPFlags(pfs *pflag.FlagSet) {
	pfs.BoolVar(&o.DryRun, "dry-run", true, "Dry run for testing. Uses API tokens but does not mutate.")

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	for _, group := range []flagutil.OptionGroup{
		o,
	} {
		group.AddFlags(fs)
	}
	pfs.AddGoFlagSet(fs)
}

func (o *GitHubOptions) BuildClient() (github.Client, error) {
	githubClient, err := o.GitHubClient(o.DryRun)
	if err != nil {
		return nil, errors.Wrap(err, "error creating github client")
	}
	return githubClient, nil
}

func (o *GitHubOptions) GetGitClientFactory() (gitv2.ClientFactory, error) {
	s := ""
	factory, err := o.GitClientFactory("", &s, o.DryRun)
	if err != nil {
		return nil, errors.Wrap(err, "error creating git client")
	}

	return factory, nil
}

func (o *GitHubOptions) Clone(org, repo string) (*gogit.Repository, error) {
	return gogit.PlainClone("", false, &gogit.CloneOptions{
		URL:      fmt.Sprintf("https://github.com/%s/%s", org, repo),
		Progress: nil,
	})
}
