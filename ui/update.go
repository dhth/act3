package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
			errorIndex := len(m.errors)
			var results []string
			for i := 0; i < 3; i++ {
				results = append(results, fmt.Sprintf("%s [%2d]", "error", errorIndex))
			}
			m.workFlowResults[msg.workflow.ID] = results
		} else {
			var results []string
			var count int
			for _, result := range msg.query.NodeResult.Workflow.Runs.Nodes {
				results = append(results, fmt.Sprintf("#%2d: %s", result.RunNumber, result.CheckSuite.Conclusion))
				if result.CheckSuite.Conclusion == "FAILURE" {
					var workflowRunKey string
					if msg.workflow.Key != nil {
						workflowRunKey = *msg.workflow.Key
					} else {
						workflowRunKey = fmt.Sprintf("%s:%s", msg.workflow.Repo, msg.workflow.Name)
					}
					m.failedWorkflowURLs[fmt.Sprintf("%s #%2d", workflowRunKey, result.RunNumber)] = result.Url
				}
				count += 1
			}
			for i := 0; i < 3-count; i++ {
				results = append(results, "")
			}
			m.workFlowResults[msg.workflow.ID] = results
		}
		m.numResults += 1
		if m.numResults >= len(m.workflows) {
			return m, quitProg()
		}
	case quitProgMsg:
		if !m.outputPrinted {
			switch m.outputFmt {
			case HTMLFmt:
				v := m.renderHTML()
				fmt.Print(v)
				m.outputPrinted = true
			}
		}
		return m, tea.Quit
	}
	return m, nil
}
