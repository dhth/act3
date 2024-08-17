package ui

import "github.com/dhth/act3/internal/gh"

type WorkflowRunsFetchedMsg struct {
	workflow gh.Workflow
	query    gh.QueryResult
	err      error
}

type quitProgMsg struct{}
