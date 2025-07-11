package cmd

import (
	"regexp"
	"sort"
	"sync"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/dhth/act3/internal/gh"
	"github.com/dhth/act3/internal/types"
)

type WorkflowError struct {
	Repo string
	Err  error
}

func getWorkflowsForRepos(ghClient *ghapi.RESTClient, repos []string, filter *regexp.Regexp) ([]types.Workflow, []WorkflowError) {
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
	var errors []WorkflowError
	for r := range resultChan {
		if r.Err != nil {
			errors = append(errors, WorkflowError{
				Repo: r.Repo,
				Err:  r.Err,
			})
			continue
		}

		for _, w := range r.Details.Workflows {
			if filter != nil && !filter.MatchString(w.Name) {
				continue
			}

			workflows = append(workflows, types.Workflow{
				ID:   w.NodeID,
				Repo: r.Repo,
				Name: w.Name,
			})
		}
	}
	sort.Slice(workflows, func(i, j int) bool {
		if workflows[i].Repo == workflows[j].Repo {
			return workflows[i].Name < workflows[j].Name
		}
		return workflows[i].Repo < workflows[j].Repo
	})

	return workflows, errors
}
