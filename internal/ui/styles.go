package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/act3/internal/gh"
)

const (
	bgColor          = "#282828"
	headerColor      = "#fe8019"
	runNumberColor   = "#83a598"
	workflowColor    = "#d3869b"
	csRunningColor   = "#fabd2f"
	csActionReqColor = "#83a598"
	csCancelledColor = "#fb4934"
	csFailureColor   = "#fb4934"
	csSuccessColor   = "#b8bb26"
	csDefaultColor   = "#928374"
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
			Width(runNumberWidth + 4)

	runResultStyle = nonFgStyle.
			Width(runNumberWidth + 4)

	getCheckRunColor = func(checkSuiteConclusion string) string {
		var color string
		switch checkSuiteConclusion {
		case gh.CSConclusionActionReq:
			color = csActionReqColor
		case gh.CSConclusionFailure, gh.CSConclusionStartupFailure:
			color = csFailureColor
		case gh.CSConclusionSuccess:
			color = csSuccessColor
		case "":
			color = errorColor
		default:
			color = csDefaultColor
		}

		return color
	}

	getResultStyle = func(checkSuiteConclusion string) lipgloss.Style {
		color := getCheckRunColor(checkSuiteConclusion)
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(color))
	}

	errorTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(errorColor))

	faintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(contextColor))

	nonSuccessHeadingStyle = nonFgStyle.
				Bold(true).
				Foreground(lipgloss.Color(csFailureColor))

	errorHeadingStyle = nonFgStyle.
				Bold(true).
				Foreground(lipgloss.Color(errorColor))

	errorDetailStyle = nonFgStyle.
				Foreground(lipgloss.Color(errorDetailColor))
)
