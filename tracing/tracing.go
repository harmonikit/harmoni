// Package tracing defines distributed tracing interfaces for harmoni.
//
// Example:
//
//	tracer := myTracingLib.NewTracer[MyReq, MyResp]()
//	ctx, span := tracer.Start(ctx, "createUser", req)
//	defer tracer.End(ctx, span, resp, err)
package tracing

import "context"

// Span represents a single operation within a trace.
type Span interface {
	// End completes the span.
	End()

	// SetAttributes sets key-value pairs on the span.
	SetAttributes(attrs ...any)
}

// Tracer creates spans for request/response operations.
type Tracer[Req, Resp any] interface {
	// Start begins a new span. It returns a context carrying the span and the
	// span itself for later use with End.
	Start(ctx context.Context, operationName string, req Req) (context.Context, Span)

	// End completes the span, recording the response and any error.
	End(ctx context.Context, span Span, resp Resp, err error)
}
