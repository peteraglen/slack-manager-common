package common

// WebhookButtonStyle represents a webhook button style.
type WebhookButtonStyle string

const (
	// WebhookButtonStylePrimary represents Slack button style 'primary'.
	WebhookButtonStylePrimary WebhookButtonStyle = "primary"

	// WebhookButtonStyleDanger represents Slack button style 'danger'.
	WebhookButtonStyleDanger WebhookButtonStyle = "danger"
)

// WebhookButtonStyleIsValid returns true if the provided WebhookButtonStyle is valid.
func WebhookButtonStyleIsValid(s WebhookButtonStyle) bool {
	switch s {
	case WebhookButtonStylePrimary, WebhookButtonStyleDanger:
		return true
	}
	return false
}

// ValidWebhookButtonStyles returns a slice of valid WebhookButtonStyle values.
func ValidWebhookButtonStyles() []string {
	return []string{
		string(WebhookButtonStylePrimary),
		string(WebhookButtonStyleDanger),
	}
}
