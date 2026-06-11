package metrics

// NopCounter implements Counter but does nothing.
type NopCounter struct{}

// NewNopCounter returns a Counter that does nothing.
func NewNopCounter() Counter { return &NopCounter{} }

func (n *NopCounter) With(labelValues ...string) Counter { return n }
func (n *NopCounter) Add(delta float64)                   {}

// NopGauge implements Gauge but does nothing.
type NopGauge struct{}

// NewNopGauge returns a Gauge that does nothing.
func NewNopGauge() Gauge { return &NopGauge{} }

func (n *NopGauge) With(labelValues ...string) Gauge { return n }
func (n *NopGauge) Set(value float64)                 {}
func (n *NopGauge) Add(delta float64)                 {}

// NopHistogram implements Histogram but does nothing.
type NopHistogram struct{}

// NewNopHistogram returns a Histogram that does nothing.
func NewNopHistogram() Histogram { return &NopHistogram{} }

func (n *NopHistogram) With(labelValues ...string) Histogram { return n }
func (n *NopHistogram) Observe(value float64)                {}
