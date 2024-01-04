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

package peribolos_test

import (
	"fmt"

	"github.com/go-git/go-billy/v5/memfs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	peribolos "k8s.io/test-infra/prow/config/org"
	"sigs.k8s.io/yaml"

	. "github.com/falcosecurity/peribolos-syncer/pkg/peribolos"
)

const (
	filename = "orgs.yaml"
	org      = "acme"
	team     = "admins"
	admin    = "alice"
	member   = "bob"
)

var _ = Describe("Creating new Peribolos config", func() {
	var config *peribolos.FullConfig

	BeforeEach(func() {
		config = NewConfig()
	})

	It("should not be nil", func() {
		Expect(config).ToNot(BeNil())
		Expect(config.Orgs).ToNot(BeNil())
	})
})

var _ = Describe("Loading Peribolos config from filesystem", func() {
	var (
		err        error
		fs         = memfs.New()
		configYaml string
		config     = &peribolos.FullConfig{Orgs: map[string]peribolos.Config{}}
	)

	Context("config is a valid peribolos fullconfig", func() {
		BeforeEach(func() {
			configYaml = fmt.Sprintf(`
orgs:
  %s:
    admins:
    - %s
    members:
    - %s
    - %s
    teams:
      %s:
        maintainers:
        - %s
        members:
        - %s
        - %s
`, org, admin, admin, member, team, admin, admin, member)

			err = yaml.Unmarshal([]byte(configYaml), config)
			Expect(err).To(Succeed())

			file, _ := fs.Create(filename)
			Expect(file).ToNot(BeNil())
			file.Write([]byte(configYaml))

			config, err = LoadConfigFromFilesystem(fs, filename)
		})

		It("should not error", func() {
			Expect(config).ToNot(BeNil())
			Expect(err).To(Succeed())
		})

		It("should build correct teams", func() {
			Expect(config.Orgs[org].Teams).To(HaveLen(1))
		})

		It("should build correct teams maintainers", func() {
			Expect(config.Orgs[org].Teams[team].Maintainers).To(HaveLen(1))
		})

		It("should build correct teams members", func() {
			Expect(config.Orgs[org].Teams[team].Members).To(HaveLen(2))
		})
	})

	Context("config file is not valid", func() {
		BeforeEach(func() {
			configYaml = `wrong`

			file, _ := fs.Create(filename)
			Expect(file).ToNot(BeNil())
			file.Write([]byte(configYaml))

			config, err = LoadConfigFromFilesystem(fs, filename)
		})

		It("should error", func() {
			Expect(err).To(HaveOccurred())
			Expect(config).To(BeNil())
		})
	})
})

var _ = Describe("Updating Team's members", func() {
	var (
		err    error
		config = &peribolos.FullConfig{Orgs: map[string]peribolos.Config{}}
	)

	BeforeEach(func() {
		config.Orgs = map[string]peribolos.Config{
			org: {
				Teams: map[string]peribolos.Team{
					team: {
						Members:     []string{"alice", "bob"},
						Maintainers: []string{"alice"},
					},
				},
				Members: nil,
				Admins:  nil,
			},
		}
	})

	Context("the org exists", func() {
		Context("the team exists", func() {
			BeforeEach(func() {
				err = AddTeamMembers(config, org, team, []string{"charlie"})
			})
			It("should not error", func() {
				Expect(err).To(Succeed())
			})
			It("should add the member to the specified team", func() {
				Expect(config.Orgs[org].Teams[team].Members).To(HaveLen(3))
			})
			It("should not remove existing members", func() {
				Expect(config.Orgs[org].Teams[team].Members).To(Equal([]string{admin, member, "charlie"}))
			})
		})

		Context("the team does not exist", func() {
			BeforeEach(func() {
				err = AddTeamMembers(config, org, "nonexistent", []string{"charlie"})
			})

			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should not change the config", func() {
				Expect(config.Orgs[org].Teams).To(HaveLen(1))
				Expect(config.Orgs[org].Teams[team].Members).To(Equal([]string{"alice", "bob"}))
				Expect(config.Orgs[org].Teams[team].Maintainers).To(Equal([]string{"alice"}))
			})
		})
	})

	Context("the org does not exist", func() {
		BeforeEach(func() {
			err = AddTeamMembers(config, "nonexistent", team, []string{"charlie"})
		})

		It("should error", func() {
			Expect(err).To(HaveOccurred())
		})

		It("should not change the config", func() {
			Expect(config.Orgs[org].Teams).To(HaveLen(1))
			Expect(config.Orgs[org].Teams[team].Members).To(Equal([]string{"alice", "bob"}))
			Expect(config.Orgs[org].Teams[team].Maintainers).To(Equal([]string{"alice"}))
		})
	})
})

var _ = Describe("Updating Team's maintainers", func() {
	var (
		err    error
		config = &peribolos.FullConfig{Orgs: map[string]peribolos.Config{}}
	)

	BeforeEach(func() {
		config.Orgs = map[string]peribolos.Config{
			org: {
				Teams: map[string]peribolos.Team{
					team: {
						Members:     []string{"alice", "bob"},
						Maintainers: []string{"alice"},
					},
				},
				Members: nil,
				Admins:  nil,
			},
		}
	})

	Context("the org exists", func() {
		Context("the team exists", func() {
			BeforeEach(func() {
				err = AddTeamMaintainers(config, org, team, []string{"charlie"})
			})
			It("should not error", func() {
				Expect(err).To(Succeed())
			})
			It("should add the member to the specified team", func() {
				Expect(config.Orgs[org].Teams[team].Maintainers).To(HaveLen(2))
			})
			It("should not remove existing members", func() {
				Expect(config.Orgs[org].Teams[team].Maintainers).To(Equal([]string{admin, "charlie"}))
			})
		})

		Context("the team does not exist", func() {
			BeforeEach(func() {
				err = AddTeamMaintainers(config, org, "nonexistent", []string{"charlie"})
			})

			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should not change the config", func() {
				Expect(config.Orgs[org].Teams).To(HaveLen(1))
				Expect(config.Orgs[org].Teams[team].Members).To(Equal([]string{"alice", "bob"}))
				Expect(config.Orgs[org].Teams[team].Maintainers).To(Equal([]string{"alice"}))
			})
		})
	})

	Context("the org does not exist", func() {
		BeforeEach(func() {
			err = AddTeamMaintainers(config, "nonexistent", team, []string{"charlie"})
		})

		It("should error", func() {
			Expect(err).To(HaveOccurred())
		})

		It("should not change the config", func() {
			Expect(config.Orgs[org].Teams).To(HaveLen(1))
			Expect(config.Orgs[org].Teams[team].Members).To(Equal([]string{"alice", "bob"}))
			Expect(config.Orgs[org].Teams[team].Maintainers).To(Equal([]string{"alice"}))
		})
	})
})
