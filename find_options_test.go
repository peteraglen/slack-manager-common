package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindOptions(t *testing.T) {
	t.Run("NewFindOptions", func(t *testing.T) {
		f := NewFindOptions()
		assert.NotNil(t, f)
		assert.Len(t, f.FieldEquals(), 0)
		assert.Len(t, f.FieldNotEquals(), 0)
	})

	t.Run("WithFieldEquals", func(t *testing.T) {
		f := NewFindOptions()
		WithFieldEquals("test", "value")(f)
		WithFieldEquals("test2", "value2")(f)
		WithFieldEquals("test2", "value2")(f)
		assert.Len(t, f.FieldEquals(), 2)
		assert.Len(t, f.FieldNotEquals(), 0)
		assert.Equal(t, "value", f.FieldEquals()["test"])
		assert.Equal(t, "value2", f.FieldEquals()["test2"])
	})

	t.Run("WithFieldNotEquals", func(t *testing.T) {
		f := NewFindOptions()
		WithFieldNotEquals("test", "value")(f)
		WithFieldNotEquals("test2", "value2")(f)
		WithFieldNotEquals("test2", "value2")(f)
		assert.Len(t, f.FieldEquals(), 0)
		assert.Len(t, f.FieldNotEquals(), 2)
		assert.Equal(t, "value", f.FieldNotEquals()["test"])
		assert.Equal(t, "value2", f.FieldNotEquals()["test2"])
	})
}
