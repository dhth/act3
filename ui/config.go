package ui

import ghapi "github.com/cli/go-gh/v2/pkg/api"

type Config struct {
	GHClient     *ghapi.GraphQLClient
	Workflows    []Workflow
	CurrentRepo  *string
	Fmt          OutputFmt
	HTMLTemplate string
}
