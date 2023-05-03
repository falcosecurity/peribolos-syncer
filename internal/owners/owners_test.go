package owners_test

import (
	"flag"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	prowflags "k8s.io/test-infra/prow/flagutil"
	gitv2 "k8s.io/test-infra/prow/git/v2"
	"k8s.io/test-infra/prow/github"
	"k8s.io/test-infra/prow/repoowners"

	. "github.com/maxgio92/peribolos-syncer/internal/owners"
)

const (
	repo       = "peribolos-syncer"
	ref        = "main"
	ownersPath = "OWNERS"
)

var _ = Describe("Creating new client", func() {
	var (
		ownersClient     *repoowners.Client
		githubOptions    *prowflags.GitHubOptions
		githubClient     github.Client
		flagset          *flag.FlagSet
		gitClientFactory gitv2.ClientFactory
	)

	BeforeEach(func() {
		githubClient = github.NewFakeClient()

		githubOptions = &prowflags.GitHubOptions{
			Host:                 "",
			TokenPath:            "",
			AllowAnonymous:       false,
			AllowDirectAccess:    false,
			AppID:                "",
			AppPrivateKeyPath:    "",
			ThrottleHourlyTokens: 0,
			ThrottleAllowBurst:   0,
			OrgThrottlers:        prowflags.Strings{},
		}

		flagset = flag.NewFlagSet("", flag.ExitOnError)
		githubOptions.AddFlags(flagset)

		s := ""
		gitClientFactory, _ = githubOptions.GitClientFactory("", &s, true)

		ownersClient = NewClient(githubClient, gitClientFactory)
	})

	It("should not be nil", func() {
		Expect(ownersClient).ToNot(BeNil())
	})
})
