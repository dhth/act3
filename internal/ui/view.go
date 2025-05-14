package ui

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/dhth/act3/internal/gh"
	"github.com/dhth/act3/internal/types"
	humanize "github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
)

const (
	runNumberWidth   = 32
	runNumberPadding = 8
	dateFormat       = "Jan 02"
)

const (
	ErrorFetchingVersion = "error"
	SystemNotFound       = "not found"
)

//go:embed assets/template.html
var htmlTemplate string

func GetOutput(config types.Config, results []gh.ResultData) (string, error) {
	switch config.Fmt {
	case types.TableFmt:
		return getTabularOutput(config, results), nil
	case types.HTMLFmt:
		return getHTMLOutput(config, results)
	default:
		return getTerminalOutput(config, results), nil
	}
}

func getTabularOutput(config types.Config, results []gh.ResultData) string {
	rows := make([][]string, len(results))

	for i, data := range results {
		var row []string
		if data.Workflow.Key != nil {
			row = append(row, *data.Workflow.Key)
		} else {
			var wf string
			if config.CurrentRepo != nil {
				wf = data.Workflow.Name
			} else {
				wf = fmt.Sprintf("%s:%s", data.Workflow.Repo, data.Workflow.Name)
			}
			row = append(row, wf)
		}
		if data.Err != nil {
			for range 3 {
				row = append(row, "ERROR")
			}
		} else {
			for _, rr := range data.Result.Workflow.Runs.Nodes {

				resultsDate := "(" + rr.CreatedAt.Format(dateFormat) + ")"
				var conclusion string
				if !rr.CheckSuite.FinishedSuccessfully() {
					conclusion = fmt.Sprintf(" %s", rr.CheckSuite.ConclusionOrState())
				}
				row = append(row, fmt.Sprintf("%s%s%s",
					RightPadTrim(fmt.Sprintf("#%d", rr.RunNumber), runNumberPadding),
					resultsDate,
					conclusion,
				))
			}
		}
		rows[i] = row
	}

	b := bytes.Buffer{}
	table := tablewriter.NewWriter(&b)

	headers := []string{"workflow", "last", "2nd-last", "3rd-last"}
	table.SetHeader(headers)

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.AppendBulk(rows)

	table.Render()

	return b.String()
}

func getTerminalOutput(config types.Config, results []gh.ResultData) string {
	var s string
	s += "\n"
	s += " " + headerStyle.Render("act3")

	if config.CurrentRepo != nil {
		s += currentRepoStyle.Render(*config.CurrentRepo)
	}
	s += "\n\n"

	s += workflowStyle.Render("workflow")

	headers := []string{"last", "2nd last", "3rd last"}
	for _, header := range headers {
		s += fmt.Sprintf("%s    ", runNumberStyle.Render(header))
	}

	s += "\n\n"

	errorIndex := 0
	var errors []error
	nonSuccessfulRuns := make(map[string]string)
	unsuccessfulRuns := false

	for _, data := range results {
		if data.Workflow.Key != nil {
			s += workflowStyle.Render(RightPadTrim(*data.Workflow.Key, runNumberWidth))
		} else {
			var wf string
			if config.CurrentRepo != nil {
				wf = data.Workflow.Name
			} else {
				wf = fmt.Sprintf("%s:%s", data.Workflow.Repo, data.Workflow.Name)
			}
			s += workflowStyle.Render(RightPadTrim(wf, runNumberWidth))
		}
		if data.Err != nil {
			for range 3 {
				s += runResultStyle.Render(fmt.Sprintf("%s %s %s",
					errorTextStyle.Render(RightPadTrim(fmt.Sprintf("#%d", errorIndex+1), runNumberPadding)),
					"😵",
					errorTextStyle.Render("(error)"),
				))
			}
			errors = append(errors, data.Err)
			errorIndex++
		} else {
			for _, rr := range data.Result.Workflow.Runs.Nodes {
				indicator := getCheckSuiteIndicator(rr.CheckSuite)
				if !rr.CheckSuite.FinishedSuccessfully() {
					var failedWorkflowRunKey string
					if data.Workflow.Key != nil {
						failedWorkflowRunKey = *data.Workflow.Key
					} else {
						if config.CurrentRepo != nil {
							failedWorkflowRunKey = data.Workflow.Name
						} else {
							failedWorkflowRunKey = fmt.Sprintf("%s: %s", data.Workflow.Repo, data.Workflow.Name)
						}
					}
					nonSuccessfulRuns[fmt.Sprintf("%s #%d (%s) ", failedWorkflowRunKey, rr.RunNumber, rr.CheckSuite.ConclusionOrState())] = rr.URL
					unsuccessfulRuns = true
				}

				style := getResultStyle(rr.CheckSuite.Conclusion)

				resultsDate := "(" + humanize.Time(rr.CreatedAt.Time) + ")"
				s += runResultStyle.Render(fmt.Sprintf("%s %s %s",
					style.Render(RightPadTrim(fmt.Sprintf("#%d", rr.RunNumber), runNumberPadding)),
					indicator,
					faintStyle.Render(resultsDate),
				))
			}
		}
		s += "\n"
	}

	if unsuccessfulRuns {
		s += "\n"
		s += nonSuccessHeadingStyle.Render("Non successful runs")
		s += "\n"
		for k, v := range nonSuccessfulRuns {
			s += errorDetailStyle.Render(fmt.Sprintf("%s%s", RightPadTrim(k, 65), v))
			s += "\n"
		}
	}

	if len(errors) > 0 {
		s += "\n"
		s += errorHeadingStyle.Render("Errors")
		s += "\n"
		for index, err := range errors {
			s += errorDetailStyle.Render(fmt.Sprintf("[#%2d]: %s", index+1, err.Error()))
			s += "\n"
		}
	}

	return s
}

func getHTMLOutput(config types.Config, results []gh.ResultData) (string, error) {
	var columns []string
	rows := make([]htmlDataRow, len(results))

	hData := htmlData{
		Title:       config.HTMLTitle,
		CurrentRepo: config.CurrentRepo,
	}

	columns = append(columns, "workflow")
	columns = append(columns, []string{"last", "2nd last", "3rd last"}...)
	errorIndex := 0
	var errors []error
	nonSuccessfulRuns := make(map[string]string)
	unsuccessfulRuns := false

	for i, data := range results {

		var workflowKey string
		if data.Workflow.Key != nil {
			workflowKey = *data.Workflow.Key
		} else {
			if config.CurrentRepo != nil {
				workflowKey = data.Workflow.Name
			} else {
				workflowKey = fmt.Sprintf("%s:%s", data.Workflow.Repo, data.Workflow.Name)
			}
		}

		var resultData []htmlWorkflowResult
		if data.Err != nil {
			for range 3 {
				resultData = append(resultData, htmlWorkflowResult{
					Details: htmlRunDetails{
						NumberFormatted: fmt.Sprintf("#%2d", errorIndex),
						Indicator:       "😵",
						Context:         "(error)",
					},
					Success: false,
					Error:   true,
					Color:   errorColor,
				})
			}
			errors = append(errors, data.Err)
			errorIndex++
		} else {
			for _, rr := range data.Result.Workflow.Runs.Nodes {
				if !rr.CheckSuite.FinishedSuccessfully() {
					var failedWorkflowRunKey string
					if data.Workflow.Key != nil {
						failedWorkflowRunKey = *data.Workflow.Key
					} else {
						if config.CurrentRepo != nil {
							failedWorkflowRunKey = data.Workflow.Name
						} else {
							failedWorkflowRunKey = fmt.Sprintf("%s: %s", data.Workflow.Repo, data.Workflow.Name)
						}
					}
					nonSuccessfulRuns[fmt.Sprintf("%s #%d (%s) ", failedWorkflowRunKey, rr.RunNumber, rr.CheckSuite.ConclusionOrState())] = rr.URL
					unsuccessfulRuns = true
				}

				success := !rr.CheckSuite.IsAFailure()
				indicator := getCheckSuiteIndicator(rr.CheckSuite)
				resultsDate := "(" + rr.CreatedAt.Format(dateFormat) + ")"

				var url string
				if data.Workflow.URL != nil {
					url = strings.ReplaceAll(*data.Workflow.URL, "{{runNumber}}", fmt.Sprintf("%d", rr.RunNumber))
				} else {
					url = rr.URL
				}
				resultData = append(resultData, htmlWorkflowResult{
					Details: htmlRunDetails{
						NumberFormatted: fmt.Sprintf("#%2d", rr.RunNumber),
						RunNumber:       fmt.Sprintf("%d", rr.RunNumber),
						Indicator:       indicator,
						Context:         resultsDate,
					},
					Success:    success,
					URL:        url,
					Conclusion: rr.CheckSuite.Conclusion,
					Color:      getCheckRunColor(rr.CheckSuite.Conclusion),
				},
				)

			}
		}
		rows[i] = htmlDataRow{
			Key:  workflowKey,
			Data: resultData,
		}
	}

	hData.Columns = columns
	hData.Rows = rows
	if len(errors) > 0 {
		hData.Errors = &errors
	}
	if unsuccessfulRuns {
		hData.Failures = nonSuccessfulRuns
	}
	hData.Timestamp = time.Now().Format("2006-01-02 15:04:05 MST")

	var tmpl *template.Template
	var err error
	if config.HTMLTemplate == "" {
		tmpl, err = template.New("act3").Parse(htmlTemplate)
	} else {
		tmpl, err = template.New("act3").Parse(config.HTMLTemplate)
	}
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, hData)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
