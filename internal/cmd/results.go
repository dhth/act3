package cmd

import (
	"sort"
	"sync"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/dhth/act3/internal/domain"
	"github.com/dhth/act3/internal/service"
)

func getResults(client *ghapi.GraphQLClient, workflows []domain.Workflow, forCurrentRepo bool) []domain.ResultData {
	semaphore := make(chan struct{}, maxConcurrentFetches)
	resultsMap := make(map[string]domain.ResultData)
	resultChannel := make(chan domain.ResultData)
	var wg sync.WaitGroup
	var results []domain.ResultData

	for _, wf := range workflows {
		wg.Add(1)
		go func(workflow domain.Workflow) {
			defer wg.Done()
			defer func() {
				<-semaphore
			}()
			semaphore <- struct{}{}
			resultChannel <- service.GetWorkflowRuns(client, workflow)
		}(wf)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	for r := range resultChannel {
		resultsMap[r.Workflow.ID] = r
	}

	if forCurrentRepo {
		resultsList := make([]domain.ResultData, 0, len(resultsMap))
		for _, r := range resultsMap {
			resultsList = append(resultsList, r)
		}
		// sort workflows alphabetically
		sort.Slice(resultsList, func(i, j int) bool {
			return resultsList[i].Workflow.Name < resultsList[j].Workflow.Name
		})
		results = resultsList
	} else {
		// sort results in the sequence of the workflows received
		resultsInConfigDefinedOrder := make([]domain.ResultData, len(workflows))
		for i, w := range workflows {
			resultsInConfigDefinedOrder[i] = resultsMap[w.ID]
		}
		results = resultsInConfigDefinedOrder
	}

	return results
}
