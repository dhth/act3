package ui

import "github.com/dhth/act3/internal/gh"

type OutputFmt uint

const (
	UnspecifiedFmt OutputFmt = iota
	HTMLFmt
)

type workflowRunResults struct {
	results    []gh.WorkflowRunNodesResult
	err        error
	errorIndex int
}

type htmlRunDetails struct {
	NumberFormatted string
	RunNumber       string
	Indicator       string
	Context         string
}

type htmlWorkflowResult struct {
	Details    htmlRunDetails
	URL        string
	Success    bool
	Color      string
	Conclusion string
	Error      bool
}

type htmlDataRow struct {
	Key  string
	Data []htmlWorkflowResult
}

type htmlData struct {
	Title       string
	CurrentRepo *string
	Columns     []string
	Rows        []htmlDataRow
	Failures    map[string]string
	Errors      *[]error
	Timestamp   string
}
