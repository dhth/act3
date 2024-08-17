package ui

import (
	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/dhth/act3/internal/gh"
)

type Config struct {
	GHClient     *ghapi.GraphQLClient
	Workflows    []gh.Workflow
	CurrentRepo  *string
	Fmt          OutputFmt
	HTMLTemplate string
}
