package metrics_test

import (
	"testing"

	"github.com/harmonikit/harmoni/metrics"
)

func BenchmarkNopCounter_Add(b *testing.B) {
	c := metrics.NewNopCounter()
	for range b.N {
		c.Add(1)
	}
}

func BenchmarkNopCounter_With(b *testing.B) {
	c := metrics.NewNopCounter()
	for range b.N {
		_ = c.With("key", "value")
	}
}

func BenchmarkNopGauge_Set(b *testing.B) {
	g := metrics.NewNopGauge()
	for range b.N {
		g.Set(1.0)
	}
}

func BenchmarkNopHistogram_Observe(b *testing.B) {
	h := metrics.NewNopHistogram()
	for range b.N {
		h.Observe(1.0)
	}
}
