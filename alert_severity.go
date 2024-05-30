package common

// AlertSeverity represents the severity for a given alert
type AlertSeverity string

const (
	// AlertPanic is used for panic situations (panic icon in Slack).
	AlertPanic AlertSeverity = "panic"

	// AlertError is used for error (critical) situations (red error icon in Slack).
	// Alerts with severity 'critical' are automatically converted to 'error'.
	AlertError AlertSeverity = "error"

	// AlertWarning is used for warning situations (yellow warning icon in Slack).
	AlertWarning AlertSeverity = "warning"

	// AlertResolved is used when a previous panic/error/warning situation has been resolved (and IssueFollowUpEnabled is true).
	// The previous icon is replaced with a green OK icon.
	// Not to be confused with an info alert!
	AlertResolved AlertSeverity = "resolved"

	// AlertInfo is used for pure info situations (blue info icon in Slack).
	// Typically used for fire-and-forget status messages, where IssueFollowUpEnabled is false.
	// Not to be confused with a resolved alert!
	AlertInfo AlertSeverity = "info"
)

func SeverityIsValid(s AlertSeverity) bool {
	switch s {
	case AlertPanic, AlertError, AlertWarning, AlertResolved, AlertInfo:
		return true
	}
	return false
}

func SeverityPriority(s AlertSeverity) int {
	switch s {
	case AlertPanic:
		return 3
	case AlertError:
		return 2
	case AlertWarning:
		return 1
	case AlertResolved, AlertInfo:
		return 0
	default:
		return -1
	}
}

func ValidSeverities() []string {
	return []string{
		string(AlertPanic),
		string(AlertError),
		string(AlertWarning),
		string(AlertResolved),
		string(AlertInfo),
	}
}
