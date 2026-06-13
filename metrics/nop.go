package metrics

// NopCounter implements Counter but does nothing.
type NopCounter struct{}

// NewNopCounter returns a Counter that does nothing.
func NewNopCounter() Counter { return &NopCounter{} }

func (n *NopCounter) With(_ ...string) Counter { return n }
func (n *NopCounter) Add(_ float64)            {}

// NopGauge implements Gauge but does nothing.
type NopGauge struct{}

// NewNopGauge returns a Gauge that does nothing.
func NewNopGauge() Gauge { return &NopGauge{} }

func (n *NopGauge) With(_ ...string) Gauge { return n }
func (n *NopGauge) Set(_ float64)          {}
func (n *NopGauge) Add(_ float64)          {}

// NopHistogram implements Histogram but does nothing.
type NopHistogram struct{}

// NewNopHistogram returns a Histogram that does nothing.
func NewNopHistogram() Histogram { return &NopHistogram{} }

func (n *NopHistogram) With(_ ...string) Histogram { return n }
func (n *NopHistogram) Observe(_ float64)          {}
