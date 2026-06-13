package tracing

import "context"

// NopSpan is a no-op span.
type NopSpan struct{}

// End is a no-op.
func (NopSpan) End() {}

// SetAttributes is a no-op.
func (NopSpan) SetAttributes(_ ...any) {}

// NopTracer is a no-op tracer that creates NopSpans.
type NopTracer[Req, Resp any] struct{}

// NewNopTracer returns a NopTracer.
func NewNopTracer[Req, Resp any]() *NopTracer[Req, Resp] {
	return &NopTracer[Req, Resp]{}
}

// Start returns the unchanged context and a NopSpan.
func (t *NopTracer[Req, Resp]) Start(ctx context.Context, _ string, _ Req) (context.Context, Span) {
	return ctx, NopSpan{}
}

// End is a no-op.
func (t *NopTracer[Req, Resp]) End(_ context.Context, _ Span, _ Resp, _ error) {}
