package cmd

import (
	"fmt"
	"os"

	"github.com/dhth/act3/internal/gh"
	"github.com/dhth/act3/internal/utils"
)

func openFailedWorkflows(results []gh.ResultData) {
	var urls []string
	for _, r := range results {
		if r.Err != nil {
			continue
		}

		for _, rr := range r.Result.Workflow.Runs.Nodes {
			if rr.CheckSuite.IsAFailure() {
				urls = append(urls, rr.URL)
			}
		}
	}

	if len(urls) == 0 {
		return
	}

	err := utils.OpenURLsInBrowser(urls)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening URLs: %s", err.Error())
	}
}
