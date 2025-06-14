package types

import (
	"strings"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/dhth/act3/internal/utils"
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

func (w Workflow) Validate() []string {
	var errors []string

	if strings.TrimSpace(w.ID) == "" {
		errors = append(errors, "workflow ID is empty")
	}

	if !utils.IsRepoNameValid(w.Repo) {
		errors = append(errors, "repo name is invalid")
	}

	if strings.TrimSpace(w.Name) == "" {
		errors = append(errors, "workflow name is empty")
	}

	if w.Key != nil && strings.TrimSpace(*w.Key) == "" {
		errors = append(errors, "workflow key is empty")
	}

	if w.URL != nil && !strings.HasPrefix(*w.URL, "https://") {
		errors = append(errors, "URL is invalid")
	}

	return errors
}
