package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlertSeverityValidation(t *testing.T) {
	assert.True(t, SeverityIsValid(AlertPanic))
	assert.True(t, SeverityIsValid(AlertError))
	assert.True(t, SeverityIsValid(AlertWarning))
	assert.True(t, SeverityIsValid(AlertResolved))
	assert.True(t, SeverityIsValid(AlertInfo))
	assert.False(t, SeverityIsValid("invalid"))
}

func TestAlertPriority(t *testing.T) {
	assert.Equal(t, 3, SeverityPriority(AlertPanic))
	assert.Equal(t, 2, SeverityPriority(AlertError))
	assert.Equal(t, 1, SeverityPriority(AlertWarning))
	assert.Equal(t, 0, SeverityPriority(AlertResolved))
	assert.Equal(t, 0, SeverityPriority(AlertInfo))
	assert.Equal(t, -1, SeverityPriority("invalid"))
}

func TestValidSeverities(t *testing.T) {
	s := ValidSeverities()
	assert.Len(t, s, 5)
	assert.Contains(t, s, "panic")
	assert.Contains(t, s, "error")
	assert.Contains(t, s, "warning")
	assert.Contains(t, s, "resolved")
	assert.Contains(t, s, "info")
}
