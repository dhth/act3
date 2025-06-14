package cmd

import (
	_ "embed"
	"errors"
	"fmt"
)

//go:embed assets/sample-config.yml
var sampleConfig string

type ErrorFollowUp struct {
	IsUnexpected bool
	Message      string
}

// returns error follow up, and whether to follow up
func GetErrorFollowUp(err error) (ErrorFollowUp, bool) {
	var zero ErrorFollowUp

	if errors.Is(err, ErrConfigFileDoesntExit) {
		return expectedErr(fmt.Sprintf(`Here's a sample config:

---
%s---`, sampleConfig))
	}

	return zero, false
}

func expectedErr(message string) (ErrorFollowUp, bool) {
	return ErrorFollowUp{
		IsUnexpected: false,
		Message:      message,
	}, true
}
