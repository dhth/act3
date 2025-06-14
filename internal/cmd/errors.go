package cmd

import (
	_ "embed"
	"errors"
	"fmt"
)

//go:embed assets/sample-config.yml
var sampleConfig string

type ErrorContext struct {
	IsUnexpected bool
	Message      string
	FollowUp     bool
}

func GetErrorContext(err error) ErrorContext {
	var zero ErrorContext

	switch {
	case errors.Is(err, ErrConfigFileDoesntExit):
		return expectedErr(fmt.Sprintf(`Here's a sample config:

---
%s---`, sampleConfig))
	case errors.Is(err, ErrCouldntGetConfig):
		return expectedErr(fmt.Sprintf(`Make sure the config looks like this:

---
%s---`, sampleConfig))
	case errors.Is(err, ErrConfigNotValid):
		return expectedErr(fmt.Sprintf(`Here's a valid config:

---
%s---`, sampleConfig))
	case errors.Is(err, ErrCouldntMarshalConfigToYAML):
		return unexpectedErr()
	case errors.Is(err, ErrNotInAGitRepo):
		return expectedErr(`If you're looking for workflows for a git repository, run act3 from the repository root.
"act3 -g" can be run from anywhere.`)
	}

	return zero
}

func unexpectedErr() ErrorContext {
	return ErrorContext{
		IsUnexpected: true,
		FollowUp:     true,
	}
}

func expectedErr(message string) ErrorContext {
	return ErrorContext{
		IsUnexpected: false,
		Message:      message,
		FollowUp:     true,
	}
}
