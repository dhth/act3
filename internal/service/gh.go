package service

import (
	"fmt"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	ghgql "github.com/cli/shurcooL-graphql"
	"github.com/dhth/act3/internal/domain"
)

func GetWorkflowDetails(ghClient *ghapi.RESTClient, repo string) domain.GetWorkflowResult {
	// https://docs.github.com/en/rest/actions/workflows?apiVersion=2022-11-28#list-repository-workflows
	var wd domain.WorkflowDetails
	err := ghClient.Get(fmt.Sprintf("repos/%s/actions/workflows", repo), &wd)
	return domain.GetWorkflowResult{Repo: repo, Details: wd, Err: err}
}

func GetWorkflowRuns(ghClient *ghapi.GraphQLClient, workflow domain.Workflow) domain.ResultData {
	variables := map[string]any{
		"numWorkflowRuns": ghgql.Int(3),
		"workflowId":      ghgql.ID(workflow.ID),
	}
	var query domain.QueryResult
	err := ghClient.Query("GetWorkflows", &query, variables)
	return domain.ResultData{
		Workflow: workflow,
		Result:   query,
		Err:      err,
	}
}
