package gh

import (
	"fmt"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

type WorkflowDetailsResult struct {
	NodeID string `json:"node_id"`
	Name   string
	State  string
}

type WorkflowDetails struct {
	TotalCount int
	Workflows  []WorkflowDetailsResult
}

func GetWorkflowDetails(ghClient *ghapi.RESTClient, repo string) (WorkflowDetails, error) {
	// https://docs.github.com/en/rest/actions/workflows?apiVersion=2022-11-28#list-repository-workflows
	var wd WorkflowDetails
	err := ghClient.Get(fmt.Sprintf("repos/%s/actions/workflows", repo), &wd)
	return wd, err
}
