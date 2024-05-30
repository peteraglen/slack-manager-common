package common

// WebhookAccessLevel is the access level required to invoke a webhook.
type WebhookAccessLevel string

const (
	// WebhookAccessLevelGlobalAdmins indicates that a webhook is available only to Slack Manager global admins.
	WebhookAccessLevelGlobalAdmins WebhookAccessLevel = "global_admins"

	// WebhookAccessLevelChannelAdmins indicates that a webhook is available to the Slack Manager channel admins (and global admins).
	WebhookAccessLevelChannelAdmins WebhookAccessLevel = "channel_admins"

	// WebhookAccessLevelChannelMembers indicates that a webhook is available to all members in a channel.
	WebhookAccessLevelChannelMembers WebhookAccessLevel = "channel_members"
)

// WebhookAccessLevelIsValid returns true if the provided WebhookAccessLevel is valid.
func WebhookAccessLevelIsValid(s WebhookAccessLevel) bool {
	switch s {
	case WebhookAccessLevelGlobalAdmins, WebhookAccessLevelChannelAdmins, WebhookAccessLevelChannelMembers:
		return true
	}
	return false
}

// ValidWebhookAccessLevels returns a slice of valid WebhookAccessLevel values.
func ValidWebhookAccessLevels() []string {
	return []string{
		string(WebhookAccessLevelGlobalAdmins),
		string(WebhookAccessLevelChannelAdmins),
		string(WebhookAccessLevelChannelMembers),
	}
}
