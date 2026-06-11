package circuitbreaker

import (
	"context"

	"github.com/harmonikit/harmoni/endpoint"
)

// NopCircuitBreaker is a pass-through circuit breaker that never opens.
// Useful as a default or for testing.
type NopCircuitBreaker[Req, Resp any] struct{}

// NewNopCircuitBreaker returns a NopCircuitBreaker.
func NewNopCircuitBreaker[Req, Resp any]() *NopCircuitBreaker[Req, Resp] {
	return &NopCircuitBreaker[Req, Resp]{}
}

// Execute delegates directly to the endpoint.
func (cb *NopCircuitBreaker[Req, Resp]) Execute(
	ctx context.Context,
	req Req,
	ep endpoint.Endpoint[Req, Resp],
) (Resp, error) {
	return ep(ctx, req)
}

// State always returns StateClosed.
func (cb *NopCircuitBreaker[Req, Resp]) State() State {
	return StateClosed
}
