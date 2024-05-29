package common

import (
	"testing"
	"time"
)

func TestNoopMetrics(t *testing.T) {
	m := &NoopMetrics{}
	m.RegisterCounter("", "")
	m.AddToCounter("", 0)
	m.AddHTTPRequestMetric("", "", 0, time.Second)
}
