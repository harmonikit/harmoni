// Package middleware provides composable middleware constructors for harmoni
// endpoints.
//
// Middleware wraps endpoints to add cross-cutting behavior:
//
//	ep = endpoint.Chain(
//	    middleware.Timeout[int, int](5*time.Second),
//	    middleware.Recovery[int, int](),
//	    middleware.Retry[int, int](3, constantBackoff),
//	)(ep)
package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/harmonikit/harmoni/endpoint"
)

// Timeout returns a middleware that sets a deadline on the request context.
// If the parent context is already done, the call returns immediately with
// the context error. If the deadline is exceeded during execution, the
// endpoint returns context.DeadlineExceeded.
func Timeout[Req, Resp any](d time.Duration) endpoint.Middleware[Req, Resp] {
	return func(next endpoint.Endpoint[Req, Resp]) endpoint.Endpoint[Req, Resp] {
		return func(ctx context.Context, req Req) (Resp, error) {
			// Fast path: if the parent context is already done, bail out.
			if err := ctx.Err(); err != nil {
				var zero Resp
				return zero, err
			}
			ctx, cancel := context.WithTimeout(ctx, d)
			defer cancel()
			return next(ctx, req)
		}
	}
}

// Retry returns a middleware that retries the endpoint on error, up to
// attempts retries (not including the initial call). The backoff function
// receives the attempt number (0-based) and returns the duration to wait
// before the next attempt.
//
// If attempts <= 0, Retry delegates to the next endpoint without retries.
func Retry[Req, Resp any](attempts int, backoff func(attempt int) time.Duration) endpoint.Middleware[Req, Resp] {
	return func(next endpoint.Endpoint[Req, Resp]) endpoint.Endpoint[Req, Resp] {
		return func(ctx context.Context, req Req) (resp Resp, err error) {
			for attempt := 0; attempt <= attempts; attempt++ {
				if ctx.Err() != nil {
					return resp, ctx.Err()
				}

				resp, err = next(ctx, req)
				if err == nil {
					return resp, nil
				}

				if attempt < attempts {
					if err := sleep(ctx, backoff(attempt)); err != nil {
						return resp, err
					}
				}
			}
			return resp, err
		}
	}
}

// Recovery returns a middleware that recovers from panics in the endpoint.
// The panic value is converted to an error and returned. The response is
// the zero value of Resp.
func Recovery[Req, Resp any]() endpoint.Middleware[Req, Resp] {
	return func(next endpoint.Endpoint[Req, Resp]) endpoint.Endpoint[Req, Resp] {
		return func(ctx context.Context, req Req) (resp Resp, err error) {
			defer func() {
				if r := recover(); r != nil {
					resp = zero[Resp]()
					err = fmt.Errorf("panic recovered: %v", r)
				}
			}()
			return next(ctx, req)
		}
	}
}

// zero returns the zero value of T. Used by Recovery to return a clean
// response after a panic.
func zero[T any]() T {
	var z T
	return z
}

// sleep waits for d duration or until ctx is done. Returns ctx.Err() if
// the context is cancelled, nil otherwise.
func sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}
