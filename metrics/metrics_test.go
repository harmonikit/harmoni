package metrics_test

import (
	"sync"
	"testing"

	"github.com/harmonikit/harmoni/metrics"
)

func TestNopCounter_Add(t *testing.T) {
	c := metrics.NewNopCounter()
	c.Add(1)
	c.Add(100)
	// No panic.
}

func TestNopCounter_With(t *testing.T) {
	c := metrics.NewNopCounter()
	labeled := c.With("method", "GET", "status", "200")
	labeled.Add(1)
	// No panic. Labeled counter should be the same NopCounter.
}

func TestNopCounter_With_Chaining(t *testing.T) {
	c := metrics.NewNopCounter()
	c.With("a", "1").With("b", "2").Add(42)
}

func TestNopGauge_Set(t *testing.T) {
	g := metrics.NewNopGauge()
	g.Set(3.14)
	g.Set(-1.0)
}

func TestNopGauge_Add(t *testing.T) {
	g := metrics.NewNopGauge()
	g.Add(1.0)
	g.Add(-0.5)
}

func TestNopGauge_With(t *testing.T) {
	g := metrics.NewNopGauge()
	g.With("region", "us-east").Set(100)
}

func TestNopHistogram_Observe(t *testing.T) {
	h := metrics.NewNopHistogram()
	h.Observe(0.001)
	h.Observe(1000.0)
}

func TestNopHistogram_With(t *testing.T) {
	h := metrics.NewNopHistogram()
	h.With("path", "/api/users").Observe(12.5)
}

func TestNop_Concurrent(t *testing.T) {
	// All Nop implementations should be safe for concurrent use.
	var wg sync.WaitGroup

	c := metrics.NewNopCounter()
	g := metrics.NewNopGauge()
	h := metrics.NewNopHistogram()

	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Add(1)
			g.Set(1.0)
			h.Observe(1.0)
		}()
	}
	wg.Wait()
}

func TestNopCounter_Interface(t *testing.T) {
	var _ metrics.Counter = metrics.NewNopCounter()
}

func TestNopGauge_Interface(t *testing.T) {
	var _ metrics.Gauge = metrics.NewNopGauge()
}

func TestNopHistogram_Interface(t *testing.T) {
	var _ metrics.Histogram = metrics.NewNopHistogram()
}
