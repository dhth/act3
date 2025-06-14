package main

import (
	"fmt"
	"os"

	"github.com/dhth/act3/internal/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		followUp, toFollowUp := cmd.GetErrorFollowUp(err)
		if !toFollowUp {
			os.Exit(1)
		}

		if followUp.Message != "" {
			fmt.Fprintf(os.Stderr, `
%s
`, followUp.Message)
		}

		if followUp.IsUnexpected {
			fmt.Fprintf(os.Stderr, `
------

This error is unexpected.
Let @dhth know about this via https://github.com/dhth/act3/issues.
`)
		}

		os.Exit(1)
	}
}
