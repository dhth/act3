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
	outputFmt          OutputFmt
	message            string
	htmlTemplate       string
	errors             []error
	failedWorkflowURLs map[string]string
	outputPrinted      bool
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, workflow := range m.workflows {
		cmds = append(cmds, getWorkflowRuns(m.ghClient, workflow))
	}
	return tea.Batch(cmds...)
}
