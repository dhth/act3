package gh

import (
	"fmt"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	ghgql "github.com/cli/shurcooL-graphql"
	"github.com/dhth/act3/internal/types"
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

	CSStateRequested  = "REQUESTED"
	CSStateQueued     = "QUEUED"
	CSStateInProgress = "IN_PROGRESS"
	CSStateCompleted  = "COMPLETED"
	CSStateWaiting    = "WAITING"
	CSStatePending    = "PENDING"
)

type CheckSuite struct {
	Status     string
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

type GetWorkflowResult struct {
	Repo    string
	Details WorkflowDetails
	Err     error
}

func GetWorkflowDetails(ghClient *ghapi.RESTClient, repo string) GetWorkflowResult {
	// https://docs.github.com/en/rest/actions/workflows?apiVersion=2022-11-28#list-repository-workflows
	var wd WorkflowDetails
	err := ghClient.Get(fmt.Sprintf("repos/%s/actions/workflows", repo), &wd)
	return GetWorkflowResult{Repo: repo, Details: wd, Err: err}
}

func (cs CheckSuite) IsAFailure() bool {
	switch cs.Conclusion {
	case CSConclusionActionReq, CSConclusionTimedOut, CSConclusionFailure, CSConclusionStartupFailure:
		return true
	default:
		return false
	}
}

func (cs CheckSuite) FinishedSuccessfully() bool {
	if cs.Status == CSStateCompleted && cs.Conclusion == CSConclusionSuccess {
		return true
	}
	return false
}

func (cs CheckSuite) ConclusionOrState() string {
	if cs.Status != CSStateCompleted {
		return cs.Status
	}
	return cs.Conclusion
}

type ResultData struct {
	Workflow types.Workflow
	Result   QueryResult
	Err      error
}

func GetWorkflowRuns(ghClient *ghapi.GraphQLClient, workflow types.Workflow) ResultData {
	variables := map[string]interface{}{
		"numWorkflowRuns": ghgql.Int(3),
		"workflowId":      ghgql.ID(workflow.ID),
	}
	var query QueryResult
	err := ghClient.Query("GetWorkflows", &query, variables)
	return ResultData{
		Workflow: workflow,
		Result:   query,
		Err:      err,
	}
}
