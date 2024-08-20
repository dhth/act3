package cmd

import (
	"fmt"
	"sort"

	"github.com/dhth/act3/internal/gh"
	"github.com/dhth/act3/internal/types"
	"github.com/dhth/act3/internal/ui"
)

func render(workflows []types.Workflow, config types.Config) error {
	resultsMap := make(map[string]gh.ResultData)
	resultChannel := make(chan gh.ResultData)
	var results []gh.ResultData

	for _, wf := range workflows {
		go func(workflow types.Workflow) {
			resultChannel <- gh.GetWorkflowRuns(config.GHClient, workflow)
		}(wf)
	}

	for range workflows {
		r := <-resultChannel
		resultsMap[r.Workflow.ID] = r
	}

	if config.CurrentRepo != nil {
		var resultsList []gh.ResultData
		for _, r := range resultsMap {
			resultsList = append(resultsList, r)
		}
		// sort workflows alphabetically
		sort.Slice(resultsList, func(i, j int) bool {
			return resultsList[i].Workflow.Name < resultsList[j].Workflow.Name
		})
		results = resultsList
	} else {
		// sort workflows in the sequence of the config file
		resultsInConfigDefinedOrder := make([]gh.ResultData, len(workflows))
		for i, w := range workflows {
			resultsInConfigDefinedOrder[i] = resultsMap[w.ID]
		}
		results = resultsInConfigDefinedOrder
	}

	output, err := ui.GetOutput(config, results)
	if err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}
