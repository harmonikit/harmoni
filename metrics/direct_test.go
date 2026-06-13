package metrics_test

import (
	"testing"

	"github.com/harmonikit/harmoni/metrics"
)

func TestNopCounter_Direct(t *testing.T) {
	c := &metrics.NopCounter{}
	c.Add(1)
	c.Add(100)

	labeled := c.With("method", "GET")
	if _, ok := labeled.(*metrics.NopCounter); !ok {
		t.Error("With should return *NopCounter")
	}
}

func TestNopGauge_Direct(t *testing.T) {
	g := &metrics.NopGauge{}
	g.Set(3.14)
	g.Add(-1.0)
	g.Add(5.0)
}

func TestNopHistogram_Direct(t *testing.T) {
	h := &metrics.NopHistogram{}
	h.Observe(1.0)
	h.Observe(100.0)
	h.With("path", "/api").Observe(0.5)
}
