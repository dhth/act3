package domain

import (
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
	Workflow Workflow
	Result   QueryResult
	Err      error
}
