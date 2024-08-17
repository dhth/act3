package ui

import "github.com/charmbracelet/lipgloss"

const (
	bgColor          = "#282828"
	headerColor      = "#fe8019"
	runNumberColor   = "#83a598"
	workflowColor    = "#d3869b"
	successColor     = "#b8bb26"
	nonSuccessColor  = "#fb4934"
	errorColor       = "#fabd2f"
	errorDetailColor = "#665c54"
	contextColor     = "#665c54"
	currentRepoColor = "#b8bb26"
)

var (
	fgStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color(bgColor))

	headerStyle = fgStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(headerColor))

	currentRepoStyle = fgStyle.
				PaddingLeft(1).
				Bold(true).
				Foreground(lipgloss.Color(currentRepoColor))

	runNumberStyle = fgStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(runNumberColor)).
			Width(runNumberWidth)

	nonFgStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	workflowStyle = nonFgStyle.
			Align(lipgloss.Left).
			Bold(true).
			Foreground(lipgloss.Color(workflowColor)).
			Width(workflowNameWidth + 4)

	runResultStyle = nonFgStyle.
			Width(runNumberWidth + 4)

	successTextStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(successColor))

	nonSuccessTextStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(nonSuccessColor))

	errorTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(errorColor))

	faintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(contextColor))

	nonSuccessHeadingStyle = nonFgStyle.
				Bold(true).
				Foreground(lipgloss.Color(nonSuccessColor))

	errorHeadingStyle = nonFgStyle.
				Bold(true).
				Foreground(lipgloss.Color(errorColor))

	errorDetailStyle = nonFgStyle.
				Foreground(lipgloss.Color(errorDetailColor))
)
