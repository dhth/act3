package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shurcooL/githubv4"
)

func getWorkflowRuns(ghClient *githubv4.Client, workflow Workflow) tea.Cmd {
	return func() tea.Msg {
		variables := map[string]interface{}{
			"numWorkflowRuns": githubv4.Int(3),
			"workflowId":      githubv4.ID(workflow.ID),
		}
		var query QueryResult
		err := ghClient.Query(context.Background(), &query, variables)

		return WorkflowRunsFetchedMsg{workflow, query, err}
	}
}

func quitProg() tea.Cmd {
	return func() tea.Msg {
		return quitProgMsg{}
	}
}
