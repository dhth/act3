package ui

func InitialModel(config Config) Model {
	workflowResults := make(map[string]workflowRunResults)

	errors := make([]error, 0)
	failedWorkflowRunURLs := make(map[string]string)

	m := Model{
		config:                 config,
		workFlowResults:        workflowResults,
		message:                "hello",
		errors:                 errors,
		nonSuccessWorkflowURLs: failedWorkflowRunURLs,
	}
	return m
}
