package ui

import "github.com/shurcooL/githubv4"

func InitialModel(ghClient *githubv4.Client, workflows []Workflow, outputFmt OutputFmt, htmlTemplate string) model {

	workflowResults := make(map[string][]string)

	for _, w := range workflows {
		var results []string
		for i := 0; i < 3; i++ {
			results = append(results, "...")
		}
		workflowResults[w.ID] = results
	}

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
