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

package peribolos

import (
	"fmt"
	"github.com/spf13/pflag"
)

type PeribolosOptions struct {
	ConfigRepo    string
	ConfigPath    string
	ConfigBaseRef string
}

const (
	gitRef = "master"
)

func NewOptions() *PeribolosOptions {
	return &PeribolosOptions{}
}

// Validate validates peribolos options. It possibly returns an error.
func (o *PeribolosOptions) Validate() error {
	if o.ConfigRepo == "" {
		//nolint:goerr113
		return fmt.Errorf("organization config file's github repository name is empty")
	}

	return nil
}

// AddPFlags adds peribolos options' flags to a flag set.
func (o *PeribolosOptions) AddPFlags(pfs *pflag.FlagSet) {
	pfs.StringVar(&o.ConfigRepo, "peribolos-config-repository", "", "The name of the github repository that contains the peribolos organization config file")
	pfs.StringVarP(&o.ConfigPath, "peribolos-config-path", "c", "org.yaml", "The path to the peribolos organization config file from the root of the Git repository")
	pfs.StringVar(&o.ConfigBaseRef, "peribolos-config-git-ref", gitRef, "The base Git reference at which pull the peribolos config repository")
}
