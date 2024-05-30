package common_test

import (
	"testing"

	common "github.com/peteraglen/slack-manager-common"
	"github.com/stretchr/testify/assert"
)

func TestFindOptions(t *testing.T) {
	t.Parallel()

	t.Run("NewFindOptions", func(t *testing.T) {
		t.Parallel()

		f := common.NewFindOptions()
		assert.NotNil(t, f)
		assert.Empty(t, f.FieldEquals())
		assert.Empty(t, f.FieldNotEquals())
	})

	t.Run("WithFieldEquals", func(t *testing.T) {
		t.Parallel()

		f := common.NewFindOptions()
		common.WithFieldEquals("test", "value")(f)
		common.WithFieldEquals("test2", "value2")(f)
		common.WithFieldEquals("test2", "value2")(f)
		assert.Len(t, f.FieldEquals(), 2)
		assert.Empty(t, f.FieldNotEquals())
		assert.Equal(t, "value", f.FieldEquals()["test"])
		assert.Equal(t, "value2", f.FieldEquals()["test2"])
	})

	t.Run("WithFieldNotEquals", func(t *testing.T) {
		t.Parallel()

		f := common.NewFindOptions()
		common.WithFieldNotEquals("test", "value")(f)
		common.WithFieldNotEquals("test2", "value2")(f)
		common.WithFieldNotEquals("test2", "value2")(f)
		assert.Empty(t, f.FieldEquals())
		assert.Len(t, f.FieldNotEquals(), 2)
		assert.Equal(t, "value", f.FieldNotEquals()["test"])
		assert.Equal(t, "value2", f.FieldNotEquals()["test2"])
	})
}
