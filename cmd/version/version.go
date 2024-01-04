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

package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/falcosecurity/peribolos-syncer/internal/output"
)

const (
	// Semantic Version that refers to ghe Version (git tag) of syncer that is released.
	// NOTE: The $Format strings are replaced during 'git archive' thanks to the
	// companion .gitattributes file containing 'export-subst' in this same
	// directory.  See also https://git-scm.com/docs/gitattributes.
	SemVersion = "v0.0.0-master+$Format:%H$"

	// sha1 from git, output of $(git rev-parse HEAD).
	GitCommit = "$Format:%H$"

	// build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ').
	BuildDate = "1970-01-01T00:00:00Z"
)

type Version struct {
	SemVersion string `json:"sem_version"`
	GitCommit  string `json:"git_commit"`
	BuildDate  string `json:"build_date"`
	GoVersion  string `json:"go_version"`
	Compiler   string `json:"compiler"`
	Platform   string `json:"platform"`
}

// New returns a new Version command.
func New() *cobra.Command {
	v := NewVersion()

	cmd := &cobra.Command{
		Use:   "Version",
		Short: "Return the syncer Version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return v.Print()
		},
	}

	return cmd
}

func (o *Version) Print() error {
	output.Print(o.SemVersion)

	return nil
}

func NewVersion() Version {
	// These variables usually come from -ldflags settings and in their
	// absence fallback to the ones defined in the var section.
	return Version{
		SemVersion: SemVersion,
		GitCommit:  GitCommit,
		BuildDate:  BuildDate,
		GoVersion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
