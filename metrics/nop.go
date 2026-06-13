package metrics

// NopCounter implements Counter but does nothing.
type NopCounter struct{}

// NewNopCounter returns a Counter that does nothing.
func NewNopCounter() Counter { return &NopCounter{} }

// With returns the NopCounter unchanged.
func (n *NopCounter) With(_ ...string) Counter { return n }

// Add is a no-op.
func (n *NopCounter) Add(_ float64) {}

// NopGauge implements Gauge but does nothing.
type NopGauge struct{}

// NewNopGauge returns a Gauge that does nothing.
func NewNopGauge() Gauge { return &NopGauge{} }

// With returns the NopGauge unchanged.
func (n *NopGauge) With(_ ...string) Gauge { return n }

// Set is a no-op.
func (n *NopGauge) Set(_ float64) {}

// Add is a no-op.
func (n *NopGauge) Add(_ float64) {}

// NopHistogram implements Histogram but does nothing.
type NopHistogram struct{}

// NewNopHistogram returns a Histogram that does nothing.
func NewNopHistogram() Histogram { return &NopHistogram{} }

// With returns the NopHistogram unchanged.
func (n *NopHistogram) With(_ ...string) Histogram { return n }

// Observe is a no-op.
func (n *NopHistogram) Observe(_ float64) {}
