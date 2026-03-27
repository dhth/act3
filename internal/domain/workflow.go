package domain

import (
	"strings"

	"github.com/dhth/act3/internal/utils"
)

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

type WorkflowDetailsResult struct {
	NodeID string `json:"node_id"`
	Name   string
	State  string
}

type WorkflowDetails struct {
	TotalCount int
	Workflows  []WorkflowDetailsResult
}

type GetWorkflowResult struct {
	Repo    string
	Details WorkflowDetails
	Err     error
}
