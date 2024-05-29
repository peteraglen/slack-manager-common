package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookDisplayMode(t *testing.T) {
	assert.True(t, WebhookDisplayModeIsValid(WebhookDisplayModeAlways))
	assert.True(t, WebhookDisplayModeIsValid(WebhookDisplayModeOpenIssue))
	assert.True(t, WebhookDisplayModeIsValid(WebhookDisplayModeResolvedIssue))
	assert.False(t, WebhookDisplayModeIsValid("invalid"))
}

func TestWebhookDisplayModeString(t *testing.T) {
	s := ValidWebhookDisplayModes()
	assert.Len(t, s, 3)
	assert.Contains(t, s, "always")
	assert.Contains(t, s, "open_issue")
	assert.Contains(t, s, "resolved_issue")
}
