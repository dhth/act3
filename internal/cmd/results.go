package cmd

import (
	"sort"
	"sync"

	"github.com/dhth/act3/internal/gh"
	"github.com/dhth/act3/internal/types"
)

func getResults(workflows []types.Workflow, config types.RunConfig) []gh.ResultData {
	semaphore := make(chan struct{}, maxConcurrentFetches)
	resultsMap := make(map[string]gh.ResultData)
	resultChannel := make(chan gh.ResultData)
	var wg sync.WaitGroup
	var results []gh.ResultData

	for _, wf := range workflows {
		wg.Add(1)
		go func(workflow types.Workflow) {
			defer wg.Done()
			defer func() {
				<-semaphore
			}()
			semaphore <- struct{}{}
			resultChannel <- gh.GetWorkflowRuns(config.GHClient, workflow)
		}(wf)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	for r := range resultChannel {
		resultsMap[r.Workflow.ID] = r
	}

	if config.CurrentRepo != nil {
		resultsList := make([]gh.ResultData, 0, len(resultsMap))
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
		resultsInConfigDefinedOrder := make([]gh.ResultData, len(workflows))
		for i, w := range workflows {
			resultsInConfigDefinedOrder[i] = resultsMap[w.ID]
		}
		results = resultsInConfigDefinedOrder
	}

	return results
}
