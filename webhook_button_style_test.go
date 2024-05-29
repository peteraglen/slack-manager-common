package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookButtonStyle(t *testing.T) {
	assert.True(t, WebhookButtonStyleIsValid(WebhookButtonStylePrimary))
	assert.True(t, WebhookButtonStyleIsValid(WebhookButtonStyleDanger))
	assert.False(t, WebhookButtonStyleIsValid("invalid"))
}

func TestWebhookButtonStyleString(t *testing.T) {
	s := ValidWebhookButtonStyles()
	assert.Len(t, s, 2)
	assert.Contains(t, s, "primary")
	assert.Contains(t, s, "danger")
}
