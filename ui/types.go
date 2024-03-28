package ui

import "github.com/shurcooL/githubv4"

type Workflow struct {
	ID   string  `yaml:"id"`
	Repo string  `yaml:"repo"`
	Name string  `yaml:"name"`
	Key  *string `yaml:"key"`
	Url  *string `yaml:"url"`
}

type WorkflowRunNodesResult struct {
	Id         string
	RunNumber  int
	Url        string
	CreatedAt  githubv4.DateTime
	CheckSuite struct {
		Conclusion string
	}
}

type WorkflowResult struct {
	Name string
	Id   string
	Runs struct {
		Nodes []WorkflowRunNodesResult
	} `graphql:"runs(first: $numWorkflowRuns)"`
}

type NodeResult struct {
	Id       string
	Workflow WorkflowResult `graphql:"... on Workflow"`
}

type QueryResult struct {
	NodeResult `graphql:"node(id: $workflowId)"`
}

type OutputFmt uint

const (
	UnspecifiedFmt OutputFmt = iota
	HTMLFmt
)

type workflowRunResults struct {
	results    []WorkflowRunNodesResult
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
	Details htmlRunDetails
	Url     string
	Success bool
	Error   bool
}

type htmlDataRow struct {
	Key  string
	Data []htmlWorkflowResult
}

type htmlData struct {
	Title     string
	Columns   []string
	Rows      []htmlDataRow
	Failures  map[string]string
	Errors    *[]error
	Timestamp string
}
