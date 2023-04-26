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

package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/maxgio92/peribolos-syncer/internal/output"
)

var (
	// Semantic version that refers to ghe version (git tag) of syncer that is released.
	// NOTE: The $Format strings are replaced during 'git archive' thanks to the
	// companion .gitattributes file containing 'export-subst' in this same
	// directory.  See also https://git-scm.com/docs/gitattributes.
	semVersion = "v0.0.0-master+$Format:%H$"

	// sha1 from git, output of $(git rev-parse HEAD).
	gitCommit = "$Format:%H$"

	// build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ').
	buildDate = "1970-01-01T00:00:00Z"
)

type version struct {
	SemVersion string `json:"sem_version"`
	GitCommit  string `json:"git_commit"`
	BuildDate  string `json:"build_date"`
	GoVersion  string `json:"go_version"`
	Compiler   string `json:"compiler"`
	Platform   string `json:"platform"`
}

// New returns a new version command.
func New() *cobra.Command {
	v := newVersion()

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Return the syncer version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return v.Print()
		},
	}

	return cmd
}

func (o *version) Print() error {
	output.Print(o.SemVersion)
	return nil
}

func newVersion() version {
	// These variables usually come from -ldflags settings and in their
	// absence fallback to the ones defined in the var section.
	return version{
		SemVersion: semVersion,
		GitCommit:  gitCommit,
		BuildDate:  buildDate,
		GoVersion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
