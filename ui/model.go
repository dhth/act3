package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shurcooL/githubv4"
)

type model struct {
	workflows          []Workflow
	ghClient           *githubv4.Client
	workFlowResults    map[string][]string
	numResults         int
	message            string
	errors             []error
	failedWorkflowURLs map[string]string
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, workflow := range m.workflows {
		cmds = append(cmds, getWorkflowRuns(m.ghClient, workflow))
	}
	return tea.Batch(cmds...)
}
