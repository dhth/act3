package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			m.message = msg.String()
		}
	case WorkflowRunsFetchedMsg:
		if msg.err != nil {
			m.errors = append(m.errors, msg.err)
			m.workFlowResults[msg.workflow.ID] = workflowRunResults{err: msg.err, errorIndex: len(m.errors) - 1}
		} else {
			m.workFlowResults[msg.workflow.ID] = workflowRunResults{results: msg.query.Workflow.Runs.Nodes, err: msg.err, errorIndex: len(m.errors)}
			for _, result := range msg.query.NodeResult.Workflow.Runs.Nodes {
				if result.CheckSuite.IsAFailure() {
					indicator := getCheckSuiteIndicator(result.CheckSuite.Conclusion)
					var failedWorkflowRunKey string
					if msg.workflow.Key != nil {
						failedWorkflowRunKey = fmt.Sprintf("%s %s", *msg.workflow.Key, indicator)
					} else {
						failedWorkflowRunKey = fmt.Sprintf("%s: %s", msg.workflow.Repo, msg.workflow.Name)
					}
					m.nonSuccessWorkflowURLs[fmt.Sprintf("%s #%2d", failedWorkflowRunKey, result.RunNumber)] = result.URL
				}
			}
			m.workFlowResults[msg.workflow.ID] = workflowRunResults{results: msg.query.Workflow.Runs.Nodes}
		}
		m.numResults++
		if m.numResults >= len(m.config.Workflows) {
			return m, quitProg()
		}
	case quitProgMsg:
		if !m.outputPrinted {
			switch m.config.Fmt {
			case HTMLFmt:
				v, err := m.renderHTML()
				// TODO: move this out to main
				if err != nil {
					fmt.Fprintf(os.Stderr, "Something went wrong generating HTML output.\nError: %s\n", err.Error())
					os.Exit(1)
				}
				fmt.Print(v)
				m.outputPrinted = true
			}
		}
		return m, tea.Quit
	}
	return m, nil
}
