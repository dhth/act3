package ui

type WorkflowRunsFetchedMsg struct {
	workflow Workflow
	query    QueryResult
	err      error
}

type quitProgMsg struct{}
