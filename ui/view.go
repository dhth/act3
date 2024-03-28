package ui

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"html/template"

	"github.com/charmbracelet/lipgloss"
	humanize "github.com/dustin/go-humanize"
)

const (
	RUN_NUMBER_WIDTH    = 30
	WORKFLOW_NAME_WIDTH = 30
)

const (
	ErrorFetchingVersion = "error"
	SystemNotFound       = "not found"
)

func (m model) renderHTML() string {

	var columns []string
	var rows []htmlDataRow

	data := htmlData{
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

		var data []htmlWorkflowResult
		workflowResults := m.workFlowResults[workflow.ID]
		if workflowResults.err != nil {
			for i := 0; i < 3; i++ {
				data = append(data, htmlWorkflowResult{
					Details: htmlRunDetails{
						NumberFormatted: fmt.Sprintf("#%2d", workflowResults.errorIndex),
						Indicator:       "😵",
						Context:         "(error)",
					},
					Success: false,
					Error:   true,
				})
			}

		} else {
			for _, rr := range workflowResults.results {
				var resultSignifier string
				var success bool
				if rr.CheckSuite.Conclusion == "SUCCESS" {
					resultSignifier = "✅"
					success = true
				} else {
					resultSignifier = "❌"
					success = false
				}
				var resultsDate = "(" + rr.CreatedAt.Time.Format("Jan 2") + ")"

				var url string
				if workflow.Url != nil {
					url = strings.Replace(*workflow.Url, "{{runNumber}}", fmt.Sprintf("%d", rr.RunNumber), -1)
				} else {
					url = rr.Url
				}
				data = append(data, htmlWorkflowResult{
					Details: htmlRunDetails{
						NumberFormatted: fmt.Sprintf("#%2d", rr.RunNumber),
						RunNumber:       fmt.Sprintf("%d", rr.RunNumber),
						Indicator:       resultSignifier,
						Context:         resultsDate,
					},
					Success: success,
					Url:     url,
				},
				)

			}
		}
		rows = append(rows, htmlDataRow{
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
			s += fmt.Sprintf("%s", workflowStyle.Render(RightPadTrim(*workflow.Key, WORKFLOW_NAME_WIDTH)))
		} else {
			s += fmt.Sprintf("%s", workflowStyle.Render(RightPadTrim(fmt.Sprintf("%s:%s", workflow.Repo, workflow.Name), WORKFLOW_NAME_WIDTH)))
		}
		workflowResults := m.workFlowResults[workflow.ID]
		if workflowResults.err != nil {
			for i := 0; i < 3; i++ {
				s += runResultStyle.Render(fmt.Sprintf("%s %s %s",
					errorTextStyle.Render(fmt.Sprintf("#%2d", workflowResults.errorIndex)),
					"😵",
					errorTextStyle.Render("(error)"),
				))
			}
		} else {

			for _, rr := range workflowResults.results {
				var resultSignifier string
				if rr.CheckSuite.Conclusion == "SUCCESS" {
					resultSignifier = "✅"
					style = successTextStyle
				} else {
					resultSignifier = "❌"
					style = failureTextStyle
				}
				var resultsDate = "(" + humanize.Time(rr.CreatedAt.Time) + ")"
				s += runResultStyle.Render(fmt.Sprintf("%s %s %s",
					style.Render(fmt.Sprintf("#%2d", rr.RunNumber)),
					resultSignifier,
					faintStyle.Render(resultsDate),
				))
			}
		}
		s += "\n"
	}

	if len(m.failedWorkflowURLs) > 0 {
		s += "\n"
		s += failureHeadingStyle.Render("Failed runs")
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
			s += errorDetailStyle.Render(fmt.Sprintf("[#%2d]: %s", index, err.Error()))
			s += "\n"
		}
	}
	return s
}
