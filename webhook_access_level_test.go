package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookAccessLevelValidation(t *testing.T) {
	assert.True(t, WebhookAccessLevelIsValid(WebhookAccessLevelGlobalAdmins))
	assert.True(t, WebhookAccessLevelIsValid(WebhookAccessLevelChannelAdmins))
	assert.True(t, WebhookAccessLevelIsValid(WebhookAccessLevelChannelMembers))
	assert.False(t, WebhookAccessLevelIsValid("invalid"))
}

func TestWebhookAccessLevelString(t *testing.T) {
	s := ValidWebhookAccessLevels()
	assert.Len(t, s, 3)
	assert.Contains(t, s, "global_admins")
	assert.Contains(t, s, "channel_admins")
	assert.Contains(t, s, "channel_members")
}
