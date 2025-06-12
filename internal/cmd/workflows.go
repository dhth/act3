package cmd

import (
	"sync"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/dhth/act3/internal/gh"
	"github.com/dhth/act3/internal/types"
)

func getWorkflowsForRepos(ghClient *ghapi.RESTClient, repos []string) ([]types.Workflow, []error) {
	semaphore := make(chan struct{}, maxConcurrentFetches)
	resultChan := make(chan gh.GetWorkflowResult)
	var wg sync.WaitGroup

	for _, repo := range repos {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			defer func() {
				<-semaphore
			}()
			semaphore <- struct{}{}
			resultChan <- gh.GetWorkflowDetails(ghClient, repo)
		}(repo)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var workflows []types.Workflow
	var errors []error
	for r := range resultChan {
		if r.Err != nil {
			errors = append(errors, r.Err)
		} else {
			for _, w := range r.Details.Workflows {
				workflows = append(workflows, types.Workflow{
					ID:   w.NodeID,
					Repo: r.Repo,
					Name: w.Name,
				})
			}
		}
	}

	return workflows, errors
}
