package gh

import (
	"fmt"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/shurcooL/githubv4"
)

// cs = check suite
// https://docs.github.com/en/graphql/reference/enums#checkconclusionstate
const (
	CSConclusionActionReq      = "ACTION_REQUIRED"
	CSConclusionCancelled      = "CANCELLED"
	CSConclusionFailure        = "FAILURE"
	CSConclusionNeutral        = "NEUTRAL"
	CSConclusionSkipped        = "SKIPPED"
	CSConclusionStartupFailure = "STARTUP_FAILURE"
	CSConclusionSuccess        = "SUCCESS"
	CSConclusionTimedOut       = "TIMED_OUT"
)

type Workflow struct {
	ID   string  `yaml:"id"`
	Repo string  `yaml:"repo"`
	Name string  `yaml:"name"`
	Key  *string `yaml:"key"`
	URL  *string `yaml:"url"`
}

type CheckSuite struct {
	Conclusion string
}

type WorkflowRunNodesResult struct {
	ID         string
	RunNumber  int
	URL        string
	CreatedAt  githubv4.DateTime
	CheckSuite CheckSuite
}

type WorkflowResult struct {
	Name string
	ID   string
	Runs struct {
		Nodes []WorkflowRunNodesResult
	} `graphql:"runs(first: $numWorkflowRuns)"`
}

type NodeResult struct {
	ID       string
	Workflow WorkflowResult `graphql:"... on Workflow"`
}

type QueryResult struct {
	NodeResult `graphql:"node(id: $workflowId)"`
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

func GetWorkflowDetails(ghClient *ghapi.RESTClient, repo string) (WorkflowDetails, error) {
	// https://docs.github.com/en/rest/actions/workflows?apiVersion=2022-11-28#list-repository-workflows
	var wd WorkflowDetails
	err := ghClient.Get(fmt.Sprintf("repos/%s/actions/workflows", repo), &wd)
	return wd, err
}

func (cs CheckSuite) IsAFailure() bool {
	switch cs.Conclusion {
	case CSConclusionActionReq, CSConclusionTimedOut, CSConclusionFailure, CSConclusionStartupFailure:
		return true
	default:
		return false
	}
}
