package domain

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
