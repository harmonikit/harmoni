package metrics_test

import (
	"fmt"

	"github.com/harmonikit/harmoni/metrics"
)

func ExampleCounter() {
	var requestCount metrics.Counter = metrics.NewNopCounter()
	requestCount.With("method", "GET", "status", "200").Add(1)
	fmt.Println("counter incremented")
	// Output: counter incremented
}

func ExampleGauge() {
	var memoryUsage metrics.Gauge = metrics.NewNopGauge()
	memoryUsage.Set(256 * 1024 * 1024) // 256 MB
	fmt.Println("gauge set")
	// Output: gauge set
}

func ExampleHistogram() {
	var latency metrics.Histogram = metrics.NewNopHistogram()
	latency.With("endpoint", "/api/users").Observe(0.012)
	fmt.Println("latency observed")
	// Output: latency observed
}
