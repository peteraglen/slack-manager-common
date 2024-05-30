package common_test

import (
	"testing"
	"time"

	common "github.com/peteraglen/slack-manager-common"
)

func TestNoopMetrics(t *testing.T) {
	t.Parallel()

	m := &common.NoopMetrics{}
	m.RegisterCounter("", "")
	m.AddToCounter("", 0)
	m.AddHTTPRequestMetric("", "", 0, time.Second)
}
