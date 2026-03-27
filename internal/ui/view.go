package ui

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/dhth/act3/internal/domain"
	humanize "github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
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

var (
	errCouldntRenderTable = errors.New("couldn't render table")
	ErrCouldntRenderHTML  = errors.New("couldn't render HTML")
)

//go:embed assets/template.html
var htmlTemplate string

func GetOutput(config domain.RunConfig, results []domain.ResultData) (string, error) {
	switch config.Fmt {
	case domain.TableFmt:
		return getTabularOutput(config, results)
	case domain.HTMLFmt:
		return getHTMLOutput(config, results)
	default:
		return getTerminalOutput(config, results), nil
	}
}

func getTabularOutput(config domain.RunConfig, results []domain.ResultData) (string, error) {
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

	headers := []string{"workflow", "last", "2nd-last", "3rd-last"}
	b := bytes.Buffer{}

	table := tablewriter.NewTable(&b,
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{
					Alignment:  tw.AlignCenter,
					AutoWrap:   tw.WrapNone,
					AutoFormat: tw.Off,
				},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{
					Alignment: tw.AlignLeft,
					AutoWrap:  tw.WrapNone,
				},
			},
		}),
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{Symbols: tw.NewSymbols(tw.StyleASCII)})),
		tablewriter.WithHeader(headers),
	)

	if err := table.Bulk(rows); err != nil {
		return "", fmt.Errorf("%w: %w", errCouldntRenderTable, err)
	}

	if err := table.Render(); err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntRenderTable, err.Error())
	}

	return b.String(), nil
}

func getTerminalOutput(config domain.RunConfig, results []domain.ResultData) string {
	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(" " + headerStyle.Render("act3"))

	if config.CurrentRepo != nil {
		s.WriteString(currentRepoStyle.Render(*config.CurrentRepo))
	}
	s.WriteString("\n\n")

	s.WriteString(workflowStyle.Render("workflow"))

	headers := []string{"last", "2nd last", "3rd last"}
	for _, header := range headers {
		fmt.Fprintf(&s, "%s    ", runNumberStyle.Render(header))
	}

	s.WriteString("\n\n")

	errorIndex := 0
	var errors []error
	nonSuccessfulRuns := make(map[string]string)
	unsuccessfulRuns := false

	for _, data := range results {
		if data.Workflow.Key != nil {
			s.WriteString(workflowStyle.Render(RightPadTrim(*data.Workflow.Key, runNumberWidth)))
		} else {
			var wf string
			if config.CurrentRepo != nil {
				wf = data.Workflow.Name
			} else {
				wf = fmt.Sprintf("%s:%s", data.Workflow.Repo, data.Workflow.Name)
			}
			s.WriteString(workflowStyle.Render(RightPadTrim(wf, runNumberWidth)))
		}
		if data.Err != nil {
			for range 3 {
				s.WriteString(runResultStyle.Render(fmt.Sprintf("%s %s %s",
					errorTextStyle.Render(RightPadTrim(fmt.Sprintf("#%d", errorIndex+1), runNumberPadding)),
					"😵",
					errorTextStyle.Render("(error)"),
				)))
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
				s.WriteString(runResultStyle.Render(fmt.Sprintf("%s %s %s",
					style.Render(RightPadTrim(fmt.Sprintf("#%d", rr.RunNumber), runNumberPadding)),
					indicator,
					faintStyle.Render(resultsDate),
				)))
			}
		}
		s.WriteString("\n")
	}

	if unsuccessfulRuns {
		s.WriteString("\n")
		s.WriteString(nonSuccessHeadingStyle.Render("Non successful runs"))
		s.WriteString("\n")
		for k, v := range nonSuccessfulRuns {
			s.WriteString(errorDetailStyle.Render(fmt.Sprintf("%s%s", RightPadTrim(k, 65), v)))
			s.WriteString("\n")
		}
	}

	if len(errors) > 0 {
		s.WriteString("\n")
		s.WriteString(errorHeadingStyle.Render("Errors"))
		s.WriteString("\n")
		for index, err := range errors {
			s.WriteString(errorDetailStyle.Render(fmt.Sprintf("[#%2d]: %s", index+1, err.Error())))
			s.WriteString("\n")
		}
	}

	return s.String()
}

func getHTMLOutput(config domain.RunConfig, results []domain.ResultData) (string, error) {
	columns := make([]string, 0, 4)
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
	hData.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05Z")

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
