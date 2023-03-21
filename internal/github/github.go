package github

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"k8s.io/test-infra/pkg/flagutil"
	prow "k8s.io/test-infra/prow/flagutil"
	gitv2 "k8s.io/test-infra/prow/git/v2"
)

type GitHubOptions struct {
	Username string

	DryRun bool

	prow.GitHubOptions
}

func (o *GitHubOptions) AddPFlags(pfs *pflag.FlagSet) {
	pfs.BoolVar(&o.DryRun, "dry-run", true, "Dry run for testing. Uses API tokens but does not mutate.")
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
