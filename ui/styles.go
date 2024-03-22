package ui

import "github.com/charmbracelet/lipgloss"

var (
	fgStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color("#282828"))

	fgStylePlain = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	headerStyle = fgStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#fe8019"))

	headerStylePlain = fgStylePlain.Copy().
				Align(lipgloss.Center)

	runNumberStyle = fgStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#83a598")).
			Width(RUN_NUMBER_WIDTH)

	nonFgStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	workflowStyle = nonFgStyle.Copy().
			Align(lipgloss.Left).
			Bold(true).
			Foreground(lipgloss.Color("#83a598")).
			Width(WORKFLOW_NAME_WIDTH)

	successResultStyle = nonFgStyle.Copy().
				Align(lipgloss.Center).
				Bold(true).
				Foreground(lipgloss.Color("#b8bb26")).
				Width(RUN_NUMBER_WIDTH)

	failureResultStyle = nonFgStyle.Copy().
				Align(lipgloss.Center).
				Bold(true).
				Foreground(lipgloss.Color("#fb4934")).
				Width(RUN_NUMBER_WIDTH).
				Underline(true)

	errorHeadingStyle = nonFgStyle.Copy().
				Bold(true).
				Foreground(lipgloss.Color("#fb4934"))

	errorDetailStyle = nonFgStyle.Copy().
				Foreground(lipgloss.Color("#665c54"))
)
