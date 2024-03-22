package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	RUN_NUMBER_WIDTH    = 25
	WORKFLOW_NAME_WIDTH = 30
)

func (m model) View() string {
	var s string

	s += "\n"
	s += " " + headerStyle.Render("act3")
	s += "\n\n"

	s += fmt.Sprintf("%s", workflowStyle.Render("workflow"))

	headers := []string{"last", "2nd last", "3rd last"}
	for _, header := range headers {
		s += fmt.Sprintf("%s    ", runNumberStyle.Render(header))
	}

	s += "\n\n"
	var style lipgloss.Style
	for _, workflow := range m.workflows {
		if workflow.Key != nil {
			s += fmt.Sprintf("%s", workflowStyle.Render(RightPadTrim(*workflow.Key, 28)))
		} else {
			s += fmt.Sprintf("%s", workflowStyle.Render(RightPadTrim(fmt.Sprintf("%s:%s", workflow.Repo, workflow.Name), 28)))
		}
		for _, workflowRunResult := range m.workFlowResults[workflow.ID] {
			if strings.Contains(workflowRunResult, "SUCCESS") {
				style = successResultStyle
			} else {
				style = failureResultStyle
			}
			s += fmt.Sprintf("%s    ", style.Render(workflowRunResult))
		}
		s += "\n"
	}

	if len(m.failedWorkflowURLs) > 0 {
		s += "\n"
		s += errorHeadingStyle.Render("Failed runs")
		s += "\n"
		for k, v := range m.failedWorkflowURLs {
			s += errorDetailStyle.Render(fmt.Sprintf("%s:\t%s", k, v))
			s += "\n"
		}
	}

	if len(m.errors) > 0 {
		s += "\n"
		s += errorHeadingStyle.Render("Errors")
		s += "\n"
		for index, err := range m.errors {
			s += errorDetailStyle.Render(fmt.Sprintf("[%2d]: %s", index+1, err.Error()))
			s += "\n"
		}
	}
	return s
}
