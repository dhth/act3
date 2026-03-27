package ui

import (
	"strings"

	"github.com/dhth/act3/internal/domain"
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

func getCheckSuiteIndicator(checkSuite domain.CheckSuite) string {
	if checkSuite.Status != domain.CSStateCompleted {
		return getCheckSuiteStateIndicator(checkSuite.Status)
	}
	return getCheckSuiteConclusionIndicator(checkSuite.Conclusion)
}

func getCheckSuiteStateIndicator(state string) string {
	switch state {
	case domain.CSStateRequested:
		return "🙏"
	case domain.CSStateQueued:
		return "⏯"
	case domain.CSStateInProgress:
		return "⏳"
	case domain.CSStateWaiting:
		return "🕓"
	case domain.CSStatePending:
		return "🟡"
	default:
		return ""
	}
}

func getCheckSuiteConclusionIndicator(conclusion string) string {
	switch conclusion {
	case domain.CSConclusionActionReq:
		return "🔄"
	case domain.CSConclusionTimedOut:
		return "⏰"
	case domain.CSConclusionCancelled:
		return "🚫"
	case domain.CSConclusionFailure:
		return "❌"
	case domain.CSConclusionSuccess:
		return "✅"
	case domain.CSConclusionNeutral:
		return "😐"
	case domain.CSConclusionSkipped:
		return "⏭️"
	case domain.CSConclusionStartupFailure:
		return "🛑"
	default:
		return "🟡"
	}
}
