// Package circuitbreaker defines the circuit breaker interface and state
// machine for harmoni.
//
// Circuit breakers prevent cascading failures by detecting when a downstream
// service is unhealthy and failing fast instead of waiting for timeouts.
//
// Implementations (Hystrix-style, Google SRE-style) live in kit/circuitbreaker/.
package circuitbreaker

import (
	"context"

	"github.com/harmonikit/harmoni/endpoint"
)

// State represents the circuit breaker's current state.
type State int

const (
	// StateClosed is the normal state — requests flow through.
	StateClosed State = iota

	// StateHalfOpen allows a limited number of requests to test if the
	// downstream service has recovered.
	StateHalfOpen

	// StateOpen rejects requests immediately without calling the endpoint.
	StateOpen
)

// String returns a human-readable state name.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return "unknown"
	}
}

// CircuitBreaker wraps an endpoint with circuit breaking logic.
type CircuitBreaker[Req, Resp any] interface {
	// Execute runs the endpoint through the circuit breaker.
	// If the circuit is open, it returns an error without calling endpoint.
	Execute(ctx context.Context, req Req, ep endpoint.Endpoint[Req, Resp]) (Resp, error)

	// State returns the current state of the circuit breaker.
	State() State
}
