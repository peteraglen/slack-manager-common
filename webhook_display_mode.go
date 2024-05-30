package common

// WebhookDisplayMode represents a display mode for a webhook button.
type WebhookDisplayMode string

const (
	// WebhookDisplayModeAlways means that the webhook button is always displayed, regardless of issue state.
	WebhookDisplayModeAlways WebhookDisplayMode = "always"

	// WebhookDisplayModeOpenIssue means that the webhook button is displayed for open issues only (error, warning etc).
	WebhookDisplayModeOpenIssue WebhookDisplayMode = "open_issue"

	// WebhookDisplayModeResolvedIssue means that the webhook button is displayed for resolved issues only.
	WebhookDisplayModeResolvedIssue WebhookDisplayMode = "resolved_issue"
)

// WebhookDisplayModeIsValid returns true if the provided WebhookDisplayMode is valid.
func WebhookDisplayModeIsValid(s WebhookDisplayMode) bool {
	switch s {
	case WebhookDisplayModeAlways, WebhookDisplayModeOpenIssue, WebhookDisplayModeResolvedIssue:
		return true
	}
	return false
}

// ValidWebhookDisplayModes returns a slice of valid WebhookDisplayMode values.
func ValidWebhookDisplayModes() []string {
	return []string{
		string(WebhookDisplayModeAlways),
		string(WebhookDisplayModeOpenIssue),
		string(WebhookDisplayModeResolvedIssue),
	}
}
