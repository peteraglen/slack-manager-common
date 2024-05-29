package common

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAlertConstructors(t *testing.T) {
	t.Run("panic alert", func(t *testing.T) {
		a := NewPanicAlert()
		assert.Equal(t, AlertPanic, a.Severity)
		assert.InDelta(t, time.Now().Unix(), a.Timestamp.Unix(), 1)
	})

	t.Run("error alert", func(t *testing.T) {
		a := NewErrorAlert()
		assert.Equal(t, AlertError, a.Severity)
		assert.InDelta(t, time.Now().Unix(), a.Timestamp.Unix(), 1)
	})

	t.Run("warning alert", func(t *testing.T) {
		a := NewWarningAlert()
		assert.Equal(t, AlertWarning, a.Severity)
		assert.InDelta(t, time.Now().Unix(), a.Timestamp.Unix(), 1)
	})

	t.Run("resolved alert", func(t *testing.T) {
		a := NewResolvedAlert()
		assert.Equal(t, AlertResolved, a.Severity)
		assert.InDelta(t, time.Now().Unix(), a.Timestamp.Unix(), 1)
	})

	t.Run("info alert", func(t *testing.T) {
		a := NewInfoAlert()
		assert.Equal(t, AlertInfo, a.Severity)
		assert.InDelta(t, time.Now().Unix(), a.Timestamp.Unix(), 1)
	})
}

func TestAlertDedupID(t *testing.T) {
	t.Run("dedup id", func(t *testing.T) {
		timestamp := time.Now()
		a := Alert{
			SlackChannelID: "C12345678",
			RouteKey:       "foo",
			CorrelationID:  "bar",
			Timestamp:      timestamp,
			Header:         "header",
			Text:           "text",
		}
		expected := hash("alert", "C12345678", "foo", "bar", timestamp.Format(time.RFC3339Nano), "header", "text")
		assert.Equal(t, expected, a.DedupID())
	})
}

func TestAlertClean(t *testing.T) {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	t.Run("timestamp newer than 7 days is kept", func(t *testing.T) {
		now := time.Now()
		a := Alert{
			Timestamp: now.Add(-7 * 24 * time.Hour).Add(10 * time.Second),
		}
		a.Clean()
		assert.Equal(t, now.Add(-7*24*time.Hour).Add(10*time.Second), a.Timestamp, "timestamp should not be updated when it's less than 7 days old")
	})

	t.Run("timestamp older than 7 days is ignored", func(t *testing.T) {
		now := time.Now()
		a := Alert{
			Timestamp: now.Add(-7 * 24 * time.Hour).Add(-1 * time.Second),
		}
		a.Clean()
		assert.InDelta(t, time.Now().Unix(), a.Timestamp.Unix(), 1, "timestamp should be updated to now when over 7 days old")
	})

	t.Run("type should be trimmed and lowercased", func(t *testing.T) {
		a := Alert{
			Type: "  FOO  ",
		}
		a.Clean()
		assert.Equal(t, "foo", a.Type)
	})

	t.Run("slackChannelID should be trimmed and uppercased", func(t *testing.T) {
		a := Alert{
			SlackChannelID: "  c12345678  ",
		}
		a.Clean()
		assert.Equal(t, "C12345678", a.SlackChannelID)
	})

	t.Run("routeKey should be trimmed and lowercased", func(t *testing.T) {
		a := Alert{
			RouteKey: "  FOO  ",
		}
		a.Clean()
		assert.Equal(t, "foo", a.RouteKey)
	})

	t.Run("header should be trimmed and newline replaced with space", func(t *testing.T) {
		a := Alert{
			Header:             "  Foo\nbar  ",
			HeaderWhenResolved: "  Hei\nresolved  ",
		}
		a.Clean()
		assert.Equal(t, "Foo bar", a.Header)
		assert.Equal(t, "Hei resolved", a.HeaderWhenResolved)
	})

	t.Run("text should be trimmed", func(t *testing.T) {
		a := Alert{
			Text:             "  Foo\nbar  ",
			TextWhenResolved: "  Hei\nresolved  ",
		}
		a.Clean()
		assert.Equal(t, "Foo\nbar", a.Text)
		assert.Equal(t, "Hei\nresolved", a.TextWhenResolved)
	})

	t.Run("fallbackText should be trimmed and simplified", func(t *testing.T) {
		a := Alert{
			FallbackText: "  Foo\nbar :status: ",
		}
		a.Clean()
		assert.Equal(t, "Foo bar", a.FallbackText)
	})

	t.Run("correlationID should be trimmed", func(t *testing.T) {
		a := Alert{
			CorrelationID: "  FOO  ",
		}
		a.Clean()
		assert.Equal(t, "FOO", a.CorrelationID)
	})

	t.Run("username should be trimmed", func(t *testing.T) {
		a := Alert{
			Username: "  FOO  ",
		}
		a.Clean()
		assert.Equal(t, "FOO", a.Username)
	})

	t.Run("author should be trimmed", func(t *testing.T) {
		a := Alert{
			Author: "  FOO  ",
		}
		a.Clean()
		assert.Equal(t, "FOO", a.Author)
	})

	t.Run("host should be trimmed", func(t *testing.T) {
		a := Alert{
			Host: "  FOO  ",
		}
		a.Clean()
		assert.Equal(t, "FOO", a.Host)
	})

	t.Run("footer should be trimmed", func(t *testing.T) {
		a := Alert{
			Footer: "  FOO  ",
		}
		a.Clean()
		assert.Equal(t, "FOO", a.Footer)
	})

	t.Run("iconEmoji should be trimmed and lowercased", func(t *testing.T) {
		a := Alert{
			IconEmoji: "  :Foo:  ",
		}
		a.Clean()
		assert.Equal(t, ":foo:", a.IconEmoji)
	})

	t.Run("severity should be trimmed and lowercased", func(t *testing.T) {
		a := Alert{
			Severity: "ERROR",
		}
		a.Clean()
		assert.Equal(t, AlertError, a.Severity)
	})

	t.Run("fallbackText should be truncated when too long", func(t *testing.T) {
		s := randString(151, randGen)
		a := Alert{
			FallbackText: s,
		}
		a.Clean()
		assert.Equal(t, s[:147]+"...", a.FallbackText)
	})

	t.Run("empty severity should default to error", func(t *testing.T) {
		a := Alert{
			Severity: "",
		}
		a.Clean()
		assert.Equal(t, AlertError, a.Severity)
	})

	t.Run("severity 'critical' should be converted to 'error'", func(t *testing.T) {
		a := Alert{
			Severity: "critical",
		}
		a.Clean()
		assert.Equal(t, AlertError, a.Severity)
	})

	t.Run("negative archivingDelaySeconds should be set to 0", func(t *testing.T) {
		a := Alert{
			ArchivingDelaySeconds: -1,
		}
		a.Clean()
		assert.Equal(t, 0, a.ArchivingDelaySeconds)
	})

	t.Run("negative notificationDelaySeconds should be set to 0", func(t *testing.T) {
		a := Alert{
			NotificationDelaySeconds: -1,
		}
		a.Clean()
		assert.Equal(t, 0, a.NotificationDelaySeconds)
	})

	t.Run("header should be truncated when too long", func(t *testing.T) {
		s := randString(131, randGen)
		s2 := randString(131, randGen)
		a := Alert{
			Header:             s,
			HeaderWhenResolved: s2,
		}
		a.Clean()
		assert.Equal(t, s[:127]+"...", a.Header)
		assert.Equal(t, s2[:127]+"...", a.HeaderWhenResolved)
	})

	t.Run("text should be truncated when too long", func(t *testing.T) {
		s := randString(10001, randGen)
		s2 := randString(10001, randGen) + "```" // Ends with code block
		a := Alert{
			Text:             s,
			TextWhenResolved: s2,
		}
		a.Clean()
		assert.Equal(t, s[:9997]+"...", a.Text)
		assert.Equal(t, s2[:9994]+"...```", a.TextWhenResolved)
	})

	t.Run("author should be truncated when too long", func(t *testing.T) {
		s := randString(101, randGen)
		a := Alert{
			Author: s,
		}
		a.Clean()
		assert.Equal(t, s[:97]+"...", a.Author)
	})

	t.Run("username should be truncated when too long", func(t *testing.T) {
		s := randString(101, randGen)
		a := Alert{
			Username: s,
		}
		a.Clean()
		assert.Equal(t, s[:97]+"...", a.Username)
	})

	t.Run("host should be truncated when too long", func(t *testing.T) {
		s := randString(101, randGen)
		a := Alert{
			Host: s,
		}
		a.Clean()
		assert.Equal(t, s[:97]+"...", a.Host)
	})

	t.Run("footer should be truncated when too long", func(t *testing.T) {
		s := randString(301, randGen)
		a := Alert{
			Footer: s,
		}
		a.Clean()
		assert.Equal(t, s[:297]+"...", a.Footer)
	})

	t.Run("field titles and values should be truncated when too long", func(t *testing.T) {
		title := randString(31, randGen)
		value := randString(201, randGen)
		a := Alert{
			Fields: []*Field{
				{Title: title, Value: value},
			},
		}
		a.Clean()
		assert.Equal(t, title[:27]+"...", a.Fields[0].Title)
		assert.Equal(t, value[:197]+"...", a.Fields[0].Value)
	})

	t.Run("webhook fields be trimmed", func(t *testing.T) {
		a := Alert{
			Webhooks: []*Webhook{
				{
					ID:               "	foo  ",
					URL:              "  http://foo.bar  ",
					ConfirmationText: "  some text  ",
					ButtonText:       "  press me  ",
					PlainTextInput: []*WebhookPlainTextInput{
						{
							ID:          "  foo  ",
							Description: "  bar  ",
						},
					},
				},
			},
		}
		a.Clean()
		assert.Equal(t, "foo", a.Webhooks[0].ID)
		assert.Equal(t, "http://foo.bar", a.Webhooks[0].URL)
		assert.Equal(t, "some text", a.Webhooks[0].ConfirmationText)
		assert.Equal(t, "press me", a.Webhooks[0].ButtonText)
		assert.Equal(t, "foo", a.Webhooks[0].PlainTextInput[0].ID)
		assert.Equal(t, "bar", a.Webhooks[0].PlainTextInput[0].Description)
	})

	t.Run("webhook button style 'default' should be replaced with empty string", func(t *testing.T) {
		a := Alert{
			Webhooks: []*Webhook{
				{
					ButtonStyle: "default",
				},
			},
		}
		a.Clean()
		assert.Equal(t, WebhookButtonStyle(""), a.Webhooks[0].ButtonStyle)
	})

	t.Run("alert escalation should be sorted by delay seconds", func(t *testing.T) {
		a := Alert{
			Escalation: []*Escalation{
				{DelaySeconds: 60},
				{DelaySeconds: 30},
			},
		}
		a.Clean()
		assert.Equal(t, 30, a.Escalation[0].DelaySeconds)
		assert.Equal(t, 60, a.Escalation[1].DelaySeconds)
	})

	t.Run("alert escalation moveToChannel should be trimmed and uppercased", func(t *testing.T) {
		a := Alert{
			Escalation: []*Escalation{
				{MoveToChannel: "  c12345678  "},
			},
		}
		a.Clean()
		assert.Equal(t, "C12345678", a.Escalation[0].MoveToChannel)
	})

	t.Run("alert escalation mentions should be trimmed", func(t *testing.T) {
		a := Alert{
			Escalation: []*Escalation{
				{SlackMentions: []string{"  <@foo>  "}},
			},
		}
		a.Clean()
		assert.Equal(t, "<@foo>", a.Escalation[0].SlackMentions[0])
	})
}

func TestAlertValidation(t *testing.T) {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	t.Run("valid minimum alert", func(t *testing.T) {
		var a *Alert
		assert.Error(t, a.Validate())
		a = &Alert{
			SlackChannelID: "C12345678",
			Header:         "foo",
		}
		a.Clean()
		assert.NoError(t, a.Validate())
	})

	t.Run("alert.slackChannelID is required when routeKey is empty", func(t *testing.T) {
		a := &Alert{
			Header: "foo",
		}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "slackChannelId")
	})

	t.Run("alert.slackChannelID and alert.routeKey should be on the correct format", func(t *testing.T) {
		a := &Alert{SlackChannelID: "abcdefghi"}
		a.Clean()
		assert.NoError(t, a.ValidateSlackChannelIDAndRouteKey())

		a = &Alert{SlackChannelID: "ABab129cf"}
		a.Clean()
		assert.NoError(t, a.ValidateSlackChannelIDAndRouteKey())

		a = &Alert{SlackChannelID: "abcdefghi9238yr"}
		a.Clean()
		assert.NoError(t, a.ValidateSlackChannelIDAndRouteKey())

		// Empty is not allowed
		a = &Alert{SlackChannelID: ""}
		a.Clean()
		assert.Error(t, a.ValidateSlackChannelIDAndRouteKey())

		// Channel names are allowed
		a = &Alert{SlackChannelID: "12345678"}
		a.Clean()
		assert.NoError(t, a.ValidateSlackChannelIDAndRouteKey())

		// Channel names are allowed
		a = &Alert{SlackChannelID: "foo-something"}
		a.Clean()
		assert.NoError(t, a.ValidateSlackChannelIDAndRouteKey())

		// Invalid characters
		a = &Alert{SlackChannelID: "sdkjsdf asdfasdf"}
		a.Clean()
		assert.Error(t, a.ValidateSlackChannelIDAndRouteKey())

		// Too long channelID
		a = &Alert{SlackChannelID: randString(MaxSlackChannelIDLength+1, randGen)}
		a.Clean()
		assert.Error(t, a.ValidateSlackChannelIDAndRouteKey())

		// routeKey is OK
		a = &Alert{RouteKey: "abcdefghi"}
		a.Clean()
		assert.NoError(t, a.ValidateSlackChannelIDAndRouteKey())

		// routeKey too long
		a = &Alert{RouteKey: randString(MaxRouteKeyLength+1, randGen)}
		a.Clean()
		assert.ErrorContains(t, a.ValidateSlackChannelIDAndRouteKey(), "routeKey")
	})

	t.Run("alert.header and alert.text cannot both be empty", func(t *testing.T) {
		a := &Alert{
			SlackChannelID: "C12345678",
		}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "header and text")
	})

	t.Run("alert.iconEmoji should be on the correct format", func(t *testing.T) {
		a := &Alert{Header: "a", RouteKey: "b", IconEmoji: ":foo:"}
		a.Clean()
		assert.NoError(t, a.Validate())

		// Invalid format
		a = &Alert{Header: "a", RouteKey: "b", IconEmoji: ":foo:bar:"}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "iconEmoji")

		// Invalid format
		a = &Alert{Header: "a", RouteKey: "b", IconEmoji: "foo"}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "iconEmoji")

		// Invalid format
		a = &Alert{Header: "a", RouteKey: "b", IconEmoji: "foo:"}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "iconEmoji")

		// Too long
		a = &Alert{Header: "a", RouteKey: "b", IconEmoji: ":" + randString(MaxIconEmojiLength+1, randGen) + ":"}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "iconEmoji")
	})

	t.Run("alert.link should be on the correct format", func(t *testing.T) {
		// Empty is OK
		a := &Alert{Header: "a", RouteKey: "b"}
		a.Clean()
		assert.NoError(t, a.Validate())

		// Invalid format
		a = &Alert{Header: "a", RouteKey: "b", Link: "foo"}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "link is not a valid absolute URL")

		// Valid format
		a = &Alert{Header: "a", RouteKey: "b", Link: "http://foo.bar?foo=bar#sfd"}
		a.Clean()
		assert.NoError(t, a.Validate())

		// Relative url is not allowed
		a = &Alert{Header: "a", RouteKey: "b", Link: "/foo"}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "link is not a valid absolute URL")
	})

	t.Run("alert.severity should be on the correct format", func(t *testing.T) {
		a := &Alert{Header: "a", RouteKey: "b", Severity: AlertError}
		a.Clean()
		assert.NoError(t, a.Validate())

		a = &Alert{Header: "a", RouteKey: "b", Severity: "foo"} // Invalid severity
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "severity")
	})

	t.Run("alert.correlationID should be on the correct format", func(t *testing.T) {
		// Empty is OK
		a := &Alert{Header: "a", RouteKey: "b"}
		a.Clean()
		assert.NoError(t, a.Validate())

		// Valid format
		a = &Alert{Header: "a", RouteKey: "b", CorrelationID: "foo"}
		a.Clean()
		assert.NoError(t, a.Validate())

		// Too long
		a = &Alert{Header: "a", RouteKey: "b", CorrelationID: randString(MaxCorrelationIDLength+1, randGen)}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "correlationId")
	})

	t.Run("alert.autoResolveSeconds should be on the correct format", func(t *testing.T) {
		// Min value is OK
		a := &Alert{Header: "a", RouteKey: "b", IssueFollowUpEnabled: true, AutoResolveSeconds: MinAutoResolveSeconds}
		a.Clean()
		assert.NoError(t, a.Validate())

		// Too small
		a = &Alert{Header: "a", RouteKey: "b", IssueFollowUpEnabled: true, AutoResolveSeconds: MinAutoResolveSeconds - 1}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "autoResolveSeconds")

		// Negative value
		a = &Alert{Header: "a", RouteKey: "b", IssueFollowUpEnabled: true, AutoResolveSeconds: -1}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "autoResolveSeconds")

		// Ignore invalid value when issueFollowUpEnabled is false
		a = &Alert{Header: "a", RouteKey: "b", IssueFollowUpEnabled: false, AutoResolveSeconds: -1}
		a.Clean()
		assert.NoError(t, a.Validate())

		// Too long
		a = &Alert{Header: "a", RouteKey: "b", IssueFollowUpEnabled: true, AutoResolveSeconds: MaxAutoResolveSeconds + 1}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "autoResolveSeconds")

		// Maximum value is OK
		a = &Alert{Header: "a", RouteKey: "b", IssueFollowUpEnabled: true, AutoResolveSeconds: MaxAutoResolveSeconds}
		a.Clean()
		assert.NoError(t, a.Validate())
	})

	t.Run("alert.ignoreIfTextContains should be on the correct format", func(t *testing.T) {
		// Empty is OK
		a := &Alert{Header: "a", RouteKey: "b", IgnoreIfTextContains: []string{}}
		a.Clean()
		assert.NoError(t, a.Validate())

		// All good
		a = &Alert{Header: "a", RouteKey: "b", IgnoreIfTextContains: []string{"foo", "bar"}}
		a.Clean()
		assert.NoError(t, a.Validate())

		// Max length is OK
		a = &Alert{Header: "a", RouteKey: "b", IgnoreIfTextContains: []string{randString(MaxIgnoreIfTextContainsLength, randGen)}}
		a.Clean()
		assert.NoError(t, a.Validate())

		// Too long
		a = &Alert{Header: "a", RouteKey: "b", IgnoreIfTextContains: []string{"foo", randString(MaxIgnoreIfTextContainsLength+1, randGen)}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "ignoreIfTextContains")
	})

	t.Run("alert.Fields should not have too many items", func(t *testing.T) {
		a := &Alert{Header: "a", RouteKey: "b"}
		for i := 1; i <= MaxFieldCount+1; i++ {
			a.Fields = append(a.Fields, &Field{Title: "foo", Value: "bar"})
		}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "too many fields")
	})

	t.Run("alert.webhooks should be on the correct format", func(t *testing.T) {
		// Empty is OK
		a := &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{}}
		a.Clean()
		assert.NoError(t, a.Validate())

		// Max webhooks is 5
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{}}
		for i := 1; i <= 6; i++ {
			a.Webhooks = append(a.Webhooks, &Webhook{ID: "foo", URL: "http://foo.bar", ButtonText: "press me"})
		}
		a.Clean()
		assert.Error(t, a.Validate(), "too many webhooks")

		// ID is required
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "", URL: "http://foo.bar", ButtonText: "press me"}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].id is required")

		// ID must be unique
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{
			{ID: "foo", URL: "http://foo.bar", ButtonText: "press me"},
			{ID: "foo", URL: "http://foo.bar", ButtonText: "press me"},
		}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[1].id must be unique")

		// Url is required
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "", ButtonText: "press me"}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].url is required")

		// Url max length is 1000
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar/" + randString(986, randGen), ButtonText: "press me"}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].url is too long")

		// Url must be valid
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "foo", ButtonText: "press me"}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].url is not a valid absolute URL")

		// Url must be absolute
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "/foo", ButtonText: "press me"}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].url is not a valid absolute URL")

		// Button text is required
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: ""}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].buttonText is required")

		// Button text max length is 25
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: randString(26, randGen)}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].buttonText is too long")

		// Confirmation text max length is 1000
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", ConfirmationText: randString(1001, randGen)}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].confirmationText is too long")

		// Button style must be valid
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", ButtonStyle: "foo"}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].buttonStyle 'foo' is not valid")

		// Access level must be valid
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", AccessLevel: "foo"}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].accessLevel 'foo' is not valid")

		// Display mode must be valid
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", DisplayMode: "foo"}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].displayMode 'foo' is not valid")

		// Max payload size is 50
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", Payload: map[string]interface{}{}}}}
		for i := 1; i <= 51; i++ {
			a.Webhooks[0].Payload[randString(10, randGen)] = randString(10, randGen)
		}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].payload item count is too large")

		// Max plain text input size is 10
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me"}}}
		for i := 1; i <= 11; i++ {
			a.Webhooks[0].PlainTextInput = append(a.Webhooks[0].PlainTextInput, &WebhookPlainTextInput{ID: randString(5, randGen), Description: randString(5, randGen)})
		}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].plainTextInput item count is too large")

		// Max checkbox size is 10
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me"}}}
		for i := 1; i <= 11; i++ {
			a.Webhooks[0].CheckboxInput = append(a.Webhooks[0].CheckboxInput, &WebhookCheckboxInput{ID: randString(5, randGen), Label: randString(5, randGen)})
		}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].checkboxInput item count is too large")

		// Input ID is required
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", PlainTextInput: []*WebhookPlainTextInput{{ID: "", Description: "foo"}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].plainTextInput[0].id is required")

		// Input ID must be unique
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", PlainTextInput: []*WebhookPlainTextInput{{ID: "foo", Description: "foo"}, {ID: "foo", Description: "foo"}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].plainTextInput[1].id must be unique")

		// Input ID max length is 200
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", PlainTextInput: []*WebhookPlainTextInput{{ID: randString(201, randGen), Description: "foo"}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].plainTextInput[0].id is too long")

		// Input description max length is 200
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", PlainTextInput: []*WebhookPlainTextInput{{ID: "foo", Description: randString(201, randGen)}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].plainTextInput[0].description is too long")

		// Input minLength cannot be negative
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", PlainTextInput: []*WebhookPlainTextInput{{ID: "foo", Description: "foo", MinLength: -1}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].plainTextInput[0].minLength must be >=0")

		// Input minstLength cannt be larger than 3000
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", PlainTextInput: []*WebhookPlainTextInput{{ID: "foo", Description: "foo", MinLength: 3001}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].plainTextInput[0].minLength must be <=3000")

		// Input maxLength cannot be negative
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", PlainTextInput: []*WebhookPlainTextInput{{ID: "foo", Description: "foo", MaxLength: -1}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].plainTextInput[0].maxLength must be >=0")

		// Input maxLength cannot be larger than 3000
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", PlainTextInput: []*WebhookPlainTextInput{{ID: "foo", Description: "foo", MaxLength: 3001}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].plainTextInput[0].maxLength must be <=3000")

		// Input minLength cannot be larger than maxLength
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", PlainTextInput: []*WebhookPlainTextInput{{ID: "foo", Description: "foo", MinLength: 10, MaxLength: 5}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].plainTextInput[0].maxLength cannot be smaller than minLength")

		// Input initialValue cannot be longer than maxLength
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", PlainTextInput: []*WebhookPlainTextInput{{ID: "foo", Description: "foo", MaxLength: 5, InitialValue: "123456"}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].plainTextInput[0].initialValue cannot be longer than maxLength")

		// Checkbox ID is required
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", CheckboxInput: []*WebhookCheckboxInput{{ID: "", Label: "foo"}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].checkboxInput[0].id is required")

		// Checkbox ID must be unique
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", CheckboxInput: []*WebhookCheckboxInput{{ID: "foo", Label: "foo"}, {ID: "foo", Label: "foo"}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].checkboxInput[1].id must be unique")

		// Checkbox ID max length is 200
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", CheckboxInput: []*WebhookCheckboxInput{{ID: randString(201, randGen), Label: "foo"}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].checkboxInput[0].id is too long")

		// Checkbox label max length is 200
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", CheckboxInput: []*WebhookCheckboxInput{{ID: "foo", Label: randString(201, randGen)}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].checkboxInput[0].label is too long")

		// Checkbox options length cannot be larger than 5
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", CheckboxInput: []*WebhookCheckboxInput{{ID: "foo", Label: "foo"}}}}}
		for i := 1; i <= 6; i++ {
			a.Webhooks[0].CheckboxInput[0].Options = append(a.Webhooks[0].CheckboxInput[0].Options, &WebhookCheckboxOption{Value: "foo"})
		}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].checkboxInput[0].options item count is too large")

		// Checkbox option value is required
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", CheckboxInput: []*WebhookCheckboxInput{{ID: "foo", Label: "foo", Options: []*WebhookCheckboxOption{{Value: ""}}}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].checkboxInput[0].options[0].value is required")

		// Checkbox option value must be unique
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", CheckboxInput: []*WebhookCheckboxInput{{ID: "foo", Label: "foo", Options: []*WebhookCheckboxOption{{Value: "foo"}, {Value: "foo"}}}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].checkboxInput[0].options[1].value must be unique")

		// Checkbox option text max length is 50
		a = &Alert{Header: "a", RouteKey: "b", Webhooks: []*Webhook{{ID: "foo", URL: "http://foo.bar", ButtonText: "press me", CheckboxInput: []*WebhookCheckboxInput{{ID: "foo", Label: "foo", Options: []*WebhookCheckboxOption{{Value: "foo", Text: randString(51, randGen)}}}}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "webhook[0].checkboxInput[0].options[0].text is too long")
	})

	t.Run("alert.escalation should be on the correct format", func(t *testing.T) {
		// Empty is OK
		a := &Alert{Header: "a", RouteKey: "b", Escalation: []*Escalation{}}
		a.Clean()
		assert.NoError(t, a.Validate())

		// Escalation delay must be at least 30 seconds
		a = &Alert{Header: "a", RouteKey: "b", Escalation: []*Escalation{{DelaySeconds: 29, Severity: AlertError}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "escalation[0].delaySeconds '29' is too low")

		// Escalation delay must be at least 30 seconds larger than the previous escalation
		a = &Alert{Header: "a", RouteKey: "b", Escalation: []*Escalation{{DelaySeconds: 30, Severity: AlertError}, {DelaySeconds: 59, Severity: AlertPanic}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "escalation[1].delaySeconds '59' is too small compared to previous escalation")

		// Escalation severity must be valid
		a = &Alert{Header: "a", RouteKey: "b", Escalation: []*Escalation{{DelaySeconds: 30, Severity: "foo"}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "escalation[0].severity 'foo' is not valid")
		a = &Alert{Header: "a", RouteKey: "b", Escalation: []*Escalation{{DelaySeconds: 30, Severity: AlertInfo}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "escalation[0].severity 'info' is not valid")
		a = &Alert{Header: "a", RouteKey: "b", Escalation: []*Escalation{{DelaySeconds: 30, Severity: AlertResolved}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "escalation[0].severity 'resolved' is not valid")

		// Escalation mentions count must be at most 10
		a = &Alert{Header: "a", RouteKey: "b", Escalation: []*Escalation{{DelaySeconds: 30, Severity: AlertError, SlackMentions: []string{}}}}
		for i := 1; i <= 11; i++ {
			a.Escalation[0].SlackMentions = append(a.Escalation[0].SlackMentions, "<@foo>")
		}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "escalation[0].slackMentions item count is too large")

		// Escalation mentions must be valid
		a = &Alert{Header: "a", RouteKey: "b", Escalation: []*Escalation{{DelaySeconds: 30, Severity: AlertError, SlackMentions: []string{"foo"}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "escalation[0].slackMentions[0] is not valid")
		a = &Alert{Header: "a", RouteKey: "b", Escalation: []*Escalation{{DelaySeconds: 30, Severity: AlertError, SlackMentions: []string{"<@" + randString(MaxMentionLength+1, randGen) + ">"}}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "escalation[0].slackMentions[0] is not valid")

		// Escalation moveToChannel must be a valid channel ID or channel name
		a = &Alert{Header: "a", RouteKey: "b", Escalation: []*Escalation{{DelaySeconds: 30, Severity: AlertError, MoveToChannel: "foo bar"}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "escalation[0].moveToChannel is not valid")
		a = &Alert{Header: "a", RouteKey: "b", Escalation: []*Escalation{{DelaySeconds: 30, Severity: AlertError, MoveToChannel: randString(81, randGen)}}}
		a.Clean()
		assert.ErrorContains(t, a.Validate(), "escalation[0].moveToChannel is not valid")
	})
}

var testLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int, randGen *rand.Rand) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = testLetters[randGen.Intn(len(testLetters))]
	}
	return string(b)
}
