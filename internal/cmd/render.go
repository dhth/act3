package cmd

import (
	"fmt"

	"github.com/dhth/act3/internal/domain"
	"github.com/dhth/act3/internal/gh"
	"github.com/dhth/act3/internal/ui"
)

func render(results []gh.ResultData, config domain.RunConfig) error {
	output, err := ui.GetOutput(config, results)
	if err != nil {
		return err
	}

	fmt.Print(output)
	return nil
}
