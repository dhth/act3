package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	config                 Config
	workFlowResults        map[string]workflowRunResults
	numResults             int
	message                string
	errors                 []error
	nonSuccessWorkflowURLs map[string]string
	outputPrinted          bool
}

func (m Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.config.Workflows))
	for i, workflow := range m.config.Workflows {
		cmds[i] = getWorkflowRuns(m.config.GHClient, workflow)
	}
	return tea.Batch(cmds...)
}
