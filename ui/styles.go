package ui

import "github.com/charmbracelet/lipgloss"

const (
	BACKGROUND_COLOR   = "#282828"
	HEADER_COLOR       = "#fe8019"
	RUN_NUMBER_COLOR   = "#83a598"
	WORKFLOW_COLOR     = "#d3869b"
	SUCCESS_COLOR      = "#b8bb26"
	FAILURE_COLOR      = "#fb4934"
	ERROR_COLOR        = "#fabd2f"
	ERROR_DETAIL_COLOR = "#665c54"
	CONTEXT_COLOR      = "#665c54"
)

var (
	fgStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color(BACKGROUND_COLOR))

	headerStyle = fgStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(HEADER_COLOR))

	runNumberStyle = fgStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(RUN_NUMBER_COLOR)).
			Width(RUN_NUMBER_WIDTH)

	nonFgStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	workflowStyle = nonFgStyle.
			Align(lipgloss.Left).
			Bold(true).
			Foreground(lipgloss.Color(WORKFLOW_COLOR)).
			Width(WORKFLOW_NAME_WIDTH)

	runResultStyle = nonFgStyle.
			PaddingLeft((RUN_NUMBER_WIDTH - 20) / 2). // TODO: This is a clumsy hack; make it better
			Width(RUN_NUMBER_WIDTH + 4)

	successTextStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(SUCCESS_COLOR))

	failureTextStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(FAILURE_COLOR))

	errorTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ERROR_COLOR))

	faintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(CONTEXT_COLOR))

	failureHeadingStyle = nonFgStyle.
				Bold(true).
				Foreground(lipgloss.Color(FAILURE_COLOR))

	errorHeadingStyle = nonFgStyle.
				Bold(true).
				Foreground(lipgloss.Color(ERROR_COLOR))

	errorDetailStyle = nonFgStyle.
				Foreground(lipgloss.Color(ERROR_DETAIL_COLOR))
)
