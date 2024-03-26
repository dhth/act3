package ui

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"html/template"

	"github.com/charmbracelet/lipgloss"
)

const (
	RUN_NUMBER_WIDTH    = 25
	WORKFLOW_NAME_WIDTH = 30
)

const (
	ErrorFetchingVersion = "error"
	SystemNotFound       = "not found"
)

func (m model) renderHTML() string {

	var columns []string
	var rows []HTMLDataRow

	data := HTMLData{
		Title: "act3",
	}

	columns = append(columns, "workflow")
	for _, env := range []string{"last", "2nd last", "3rd last"} {
		columns = append(columns, env)
	}

	for _, workflow := range m.workflows {

		var workflowKey string
		if workflow.Key != nil {
			workflowKey = *workflow.Key
		} else {
			workflowKey = fmt.Sprintf("%s:%s", workflow.Repo, workflow.Name)
		}

		var data []HTMLWorkflowResult
		for _, workflowRunResult := range m.workFlowResults[workflow.ID] {
			data = append(data, HTMLWorkflowResult{
				Result:  workflowRunResult,
				Success: strings.Contains(workflowRunResult, "SUCCESS"),
			})
		}
		rows = append(rows, HTMLDataRow{
			Key:  workflowKey,
			Data: data,
		})
	}

	data.Columns = columns
	data.Rows = rows
	if len(m.errors) > 0 {
		data.Errors = &m.errors
	}
	if len(m.failedWorkflowURLs) > 0 {
		data.Failures = m.failedWorkflowURLs
	}
	data.Timestamp = time.Now().Format("2006-01-02 15:04:05 MST")

	var tmpl *template.Template
	var err error
	if m.htmlTemplate == "" {
		tmpl, err = template.New("act3").Parse(HTMLTemplText)
	} else {
		tmpl, err = template.New("act3").Parse(m.htmlTemplate)
	}
	if err != nil {
		return fmt.Sprintf(string(errorTemplate), err.Error())
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Sprintf(string(errorTemplate), err.Error())
	}

	return buf.String()
}

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
