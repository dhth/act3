package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
	ghgql "github.com/cli/shurcooL-graphql"
	"github.com/dhth/act3/internal/gh"
)

func getWorkflowRuns(ghClient *ghapi.GraphQLClient, workflow gh.Workflow) tea.Cmd {
	return func() tea.Msg {
		variables := map[string]interface{}{
			"numWorkflowRuns": ghgql.Int(3),
			"workflowId":      ghgql.ID(workflow.ID),
		}
		var query gh.QueryResult
		err := ghClient.Query("GetWorkflows", &query, variables)

		return WorkflowRunsFetchedMsg{workflow, query, err}
	}
}

func quitProg() tea.Cmd {
	return func() tea.Msg {
		return quitProgMsg{}
	}
}
