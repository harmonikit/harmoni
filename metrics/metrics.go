// Package metrics defines interfaces for application metrics.
//
// It provides Counter, Gauge, and Histogram abstractions that can be backed
// by Prometheus, OpenTelemetry, expvar, or any other metrics backend.
//
// Example:
//
//	var requestCount metrics.Counter
//	requestCount.With("method", "GET").Add(1)
package metrics

// Counter is a monotonically increasing metric.
type Counter interface {
	// With returns a new Counter with the given label values added.
	// Implementations should not mutate the receiver.
	With(labelValues ...string) Counter

	// Add increments the counter by delta. Delta must be >= 0.
	Add(delta float64)
}

// Gauge is a metric that can go up and down.
type Gauge interface {
	// With returns a new Gauge with the given label values added.
	// Implementations should not mutate the receiver.
	With(labelValues ...string) Gauge

	// Set sets the gauge to an absolute value.
	Set(value float64)

	// Add increments or decrements the gauge by delta.
	Add(delta float64)
}

// Histogram records the distribution of observed values.
type Histogram interface {
	// With returns a new Histogram with the given label values added.
	// Implementations should not mutate the receiver.
	With(labelValues ...string) Histogram

	// Observe records a value in the histogram.
	Observe(value float64)
}
