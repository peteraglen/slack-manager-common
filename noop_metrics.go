package common

import "time"

// NoopMetrics is a no-op implementation of the Metrics interface.
// It does not record any metrics. Use if no metrics are needed.
type NoopMetrics struct{}

func (m *NoopMetrics) RegisterCounter(_, _ string, _ ...string) {
}

func (m *NoopMetrics) AddToCounter(_ string, _ float64, _ ...string) {
}

func (m *NoopMetrics) AddHTTPRequestMetric(_, _ string, _ int, _ time.Duration) {
}
