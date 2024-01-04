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

package version_test

import (
	"fmt"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmd "github.com/falcosecurity/peribolos-syncer/cmd/version"
)

var _ = Describe("Version", func() {
	version := &cmd.Version{
		SemVersion: "0.0.0",
		GitCommit:  "ef30d58",
		BuildDate:  "2024-03-25_17:55:07",
		GoVersion:  "go1.20",
		Compiler:   "gc",
		Platform:   "linux/test",
	}

	Context("testing Version Print function", func() {
		It("should not error", func() {
			Expect(version.Print()).Error().ShouldNot(HaveOccurred())
		})
	})

	Context("testing NewVersion function", func() {
		It("should return the expected Version struct", func() {
			v := cmd.NewVersion()
			Expect(v.Compiler).Should(Equal(runtime.Compiler))
			Expect(v.SemVersion).Should(Equal(cmd.SemVersion))
			Expect(v.GitCommit).Should(Equal(cmd.GitCommit))
			Expect(v.BuildDate).Should(Equal(cmd.BuildDate))
			Expect(v.GoVersion).Should(Equal(runtime.Version()))
			Expect(v.Platform).Should(Equal(fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)))
		})
	})
})

func TestVersion(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Version Suite")
}
