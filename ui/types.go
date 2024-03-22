package ui

type Workflow struct {
	ID   string  `yaml:"id"`
	Repo string  `yaml:"repo"`
	Name string  `yaml:"name"`
	Key  *string `yaml:"key"`
}

type CommitResult struct {
	Oid     string
	Message string
}

type WorkflowRunNodesResult struct {
	Id         string
	RunNumber  int
	Url        string
	CreatedAt  string
	CheckSuite struct {
		Commit     CommitResult
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
