package endpoint

import "context"

// Endpoint is a single RPC method. It is the universal abstraction for any
// request/response interaction — HTTP handler, gRPC method, or in-process call.
type Endpoint[Req, Resp any] func(ctx context.Context, req Req) (Resp, error)

// Middleware wraps an Endpoint, adding cross-cutting behavior.
// Middleware composes via Chain.
type Middleware[Req, Resp any] func(Endpoint[Req, Resp]) Endpoint[Req, Resp]

// Chain composes middleware left-to-right. The first middleware in the list
// is the outermost — it executes first on the way in and last on the way out.
//
//	Chain(logMW, metricsMW, timeoutMW)
//	// Results in: logMW(metricsMW(timeoutMW(endpoint)))
func Chain[Req, Resp any](outer Middleware[Req, Resp], rest ...Middleware[Req, Resp]) Middleware[Req, Resp] {
	return func(next Endpoint[Req, Resp]) Endpoint[Req, Resp] {
		// Apply right-to-left so leftmost is outermost.
		wrapped := next
		for i := len(rest) - 1; i >= 0; i-- {
			wrapped = rest[i](wrapped)
		}
		return outer(wrapped)
	}
}

// Nop returns an Endpoint that does nothing — it returns the zero value of Resp
// and nil error. Useful as a placeholder or for testing middleware chains.
func Nop[Req, Resp any]() Endpoint[Req, Resp] {
	return func(_ context.Context, _ Req) (Resp, error) {
		var zero Resp
		return zero, nil
	}
}
