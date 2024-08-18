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

type Config struct {
	GHClient     *ghapi.GraphQLClient
	CurrentRepo  *string
	Fmt          OutputFmt
	HTMLTemplate string
}

type Workflow struct {
	ID   string  `yaml:"id"`
	Repo string  `yaml:"repo"`
	Name string  `yaml:"name"`
	Key  *string `yaml:"key"`
	URL  *string `yaml:"url"`
}
