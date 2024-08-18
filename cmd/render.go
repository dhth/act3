package cmd

import (
	"fmt"
	"sort"

	"github.com/dhth/act3/internal/gh"
	"github.com/dhth/act3/internal/types"
	"github.com/dhth/act3/internal/ui"
)

func render(workflows []types.Workflow, config types.Config) error {
	results := make([]gh.ResultData, len(workflows))
	resultChannel := make(chan gh.ResultData)

	for _, wf := range workflows {
		go func(workflow types.Workflow) {
			resultChannel <- gh.GetWorkflowRuns(config.GHClient, workflow)
		}(wf)
	}

	for i := range workflows {
		r := <-resultChannel
		results[i] = r
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Workflow.Name < results[j].Workflow.Name
	})

	output, err := ui.GetOutput(config, results)
	if err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}
