package types

import (
	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

type OutputFmt uint

const (
	DefaultFmt OutputFmt = iota
	TableFmt
	HTMLFmt
)

type RunConfig struct {
	GHClient     *ghapi.GraphQLClient
	CurrentRepo  *string
	Fmt          OutputFmt
	HTMLTitle    string
	HTMLTemplate string
}

type Workflow struct {
	ID   string  `yaml:"id"`
	Repo string  `yaml:"repo"`
	Name string  `yaml:"name"`
	Key  *string `yaml:"key,omitempty"`
	URL  *string `yaml:"url,omitempty"`
}
