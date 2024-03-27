package ui

import "github.com/shurcooL/githubv4"

const (
	NUM_RUNS_TO_DISPLAY = 3
)

func InitialModel(ghClient *githubv4.Client, workflows []Workflow, outputFmt OutputFmt, htmlTemplate string) model {

	workflowResults := make(map[string]workflowRunResults)

	errors := make([]error, 0)
	failedWorkflowRunURLs := make(map[string]string)

	m := model{
		ghClient:           ghClient,
		workflows:          workflows,
		workFlowResults:    workflowResults,
		outputFmt:          outputFmt,
		htmlTemplate:       htmlTemplate,
		message:            "hello",
		errors:             errors,
		failedWorkflowURLs: failedWorkflowRunURLs,
	}
	return m
}
