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
