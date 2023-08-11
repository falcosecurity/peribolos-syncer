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

package owners_test

import (
	"flag"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	prowflags "k8s.io/test-infra/prow/flagutil"
	gitv2 "k8s.io/test-infra/prow/git/v2"
	"k8s.io/test-infra/prow/github"
	"k8s.io/test-infra/prow/repoowners"

	. "github.com/falcosecurity/peribolos-syncer/internal/owners"
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
