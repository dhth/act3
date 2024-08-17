package ui

import (
	"strings"

	"github.com/dhth/act3/internal/gh"
)

func RightPadTrim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

func Trim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s
}

func getCheckSuiteIndicator(conclusion string) string {
	switch conclusion {
	case gh.CSConclusionActionReq:
		return "🔄"
	case gh.CSConclusionTimedOut:
		return "⏰"
	case gh.CSConclusionCancelled:
		return "🚫"
	case gh.CSConclusionFailure:
		return "❌"
	case gh.CSConclusionSuccess:
		return "✅"
	case gh.CSConclusionNeutral:
		return "😐"
	case gh.CSConclusionSkipped:
		return "⏭️"
	case gh.CSConclusionStartupFailure:
		return "🛑"
	default:
		return "🟡"
	}
}
