package ui

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	humanize "github.com/dustin/go-humanize"
)

const (
	runNumberWidth    = 30
	workflowNameWidth = 30
	runNumberPadding  = 6
)

const (
	ErrorFetchingVersion = "error"
	SystemNotFound       = "not found"
)

var (
	//go:embed assets/template.html
	htmlTemplate string

	//go:embed assets/error.html
	htmlErrorTemplate string
)

func (m Model) renderHTML() string {
	var columns []string
	rows := make([]htmlDataRow, len(m.config.Workflows))

	data := htmlData{
		Title:       "act3",
		CurrentRepo: m.config.CurrentRepo,
	}

	columns = append(columns, "workflow")
	columns = append(columns, []string{"last", "2nd last", "3rd last"}...)

	for i, workflow := range m.config.Workflows {

		var workflowKey string
		if workflow.Key != nil {
			workflowKey = *workflow.Key
		} else {
			if m.config.CurrentRepo != nil {
				workflowKey = workflow.Name
			} else {
				workflowKey = fmt.Sprintf("%s:%s", workflow.Repo, workflow.Name)
			}
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
				resultsDate := "(" + rr.CreatedAt.Time.Format("Jan 2") + ")"

				var url string
				if workflow.URL != nil {
					url = strings.Replace(*workflow.URL, "{{runNumber}}", fmt.Sprintf("%d", rr.RunNumber), -1)
				} else {
					url = rr.URL
				}
				data = append(data, htmlWorkflowResult{
					Details: htmlRunDetails{
						NumberFormatted: fmt.Sprintf("#%2d", rr.RunNumber),
						RunNumber:       fmt.Sprintf("%d", rr.RunNumber),
						Indicator:       resultSignifier,
						Context:         resultsDate,
					},
					Success: success,
					URL:     url,
				},
				)

			}
		}
		rows[i] = htmlDataRow{
			Key:  workflowKey,
			Data: data,
		}
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
	if m.config.HTMLTemplate == "" {
		tmpl, err = template.New("act3").Parse(htmlTemplate)
	} else {
		tmpl, err = template.New("act3").Parse(m.config.HTMLTemplate)
	}
	if err != nil {
		return fmt.Sprintf(htmlErrorTemplate, err.Error())
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Sprintf(htmlErrorTemplate, err.Error())
	}

	return buf.String()
}

func (m Model) View() string {
	var s string

	s += "\n"
	s += " " + headerStyle.Render("act3")
	if m.config.CurrentRepo != nil {
		s += currentRepoStyle.Render(*m.config.CurrentRepo)
	}
	s += "\n\n"

	s += workflowStyle.Render("workflow")

	headers := []string{"last", "2nd last", "3rd last"}
	for _, header := range headers {
		s += fmt.Sprintf("%s    ", runNumberStyle.Render(header))
	}

	s += "\n\n"
	var style lipgloss.Style
	for _, workflow := range m.config.Workflows {
		if workflow.Key != nil {
			s += workflowStyle.Render(RightPadTrim(*workflow.Key, workflowNameWidth))
		} else {
			var wf string
			if m.config.CurrentRepo != nil {
				wf = workflow.Name
			} else {
				wf = fmt.Sprintf("%s:%s", workflow.Repo, workflow.Name)
			}
			s += workflowStyle.Render(RightPadTrim(wf, workflowNameWidth))
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
				resultsDate := "(" + humanize.Time(rr.CreatedAt.Time) + ")"
				s += runResultStyle.Render(fmt.Sprintf("%s %s %s",
					style.Render(RightPadTrim(fmt.Sprintf("#%d", rr.RunNumber), runNumberPadding)),
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
			s += errorDetailStyle.Render(fmt.Sprintf("%s%s", RightPadTrim(k, 65), v))
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
