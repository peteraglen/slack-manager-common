package types_test

import (
	"testing"

	"github.com/slackmgr/types"
	"github.com/slackmgr/types/dbtests"
)

func TestInMemoryDB(t *testing.T) {
	t.Parallel()

	dbtests.RunAllTests(t, types.NewInMemoryDB())
}
