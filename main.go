package main

import (
	"fmt"
	"os"

	"github.com/dhth/act3/internal/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		errContext := cmd.GetErrorContext(err)
		if !errContext.FollowUp {
			os.Exit(1)
		}

		if errContext.Message != "" {
			fmt.Fprintf(os.Stderr, `
%s
`, errContext.Message)
		}

		if errContext.IsUnexpected {
			fmt.Fprintf(os.Stderr, `
------

This error is unexpected.
Let @dhth know about this via https://github.com/dhth/act3/issues.
`)
		}

		os.Exit(1)
	}
}
