package endpoint_test

import (
	"context"
	"testing"

	"github.com/harmonikit/harmoni/endpoint"
)

// Benchmarks for the hot path: endpoint invocation with middleware chains.

func BenchmarkEndpoint_Direct(b *testing.B) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
	})
	ctx := context.Background()

	for range b.N {
		_, _ = ep(ctx, 0)
	}
}

func BenchmarkEndpoint_ThroughMiddleware(b *testing.B) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
	})

	nopMW := endpoint.Middleware[int, int](func(next endpoint.Endpoint[int, int]) endpoint.Endpoint[int, int] {
		return func(ctx context.Context, req int) (int, error) {
			return next(ctx, req)
		}
	})

	ep = nopMW(ep)
	ctx := context.Background()

	for range b.N {
		_, _ = ep(ctx, 0)
	}
}

func BenchmarkEndpoint_Chain5(b *testing.B) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
	})

	nopMW := endpoint.Middleware[int, int](func(next endpoint.Endpoint[int, int]) endpoint.Endpoint[int, int] {
		return func(ctx context.Context, req int) (int, error) {
			return next(ctx, req)
		}
	})

	wrapped := endpoint.Chain(nopMW, nopMW, nopMW, nopMW, nopMW)(ep)
	ctx := context.Background()

	for range b.N {
		_, _ = wrapped(ctx, 0)
	}
}

func BenchmarkEndpoint_Chain10(b *testing.B) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
	})

	nopMW := endpoint.Middleware[int, int](func(next endpoint.Endpoint[int, int]) endpoint.Endpoint[int, int] {
		return func(ctx context.Context, req int) (int, error) {
			return next(ctx, req)
		}
	})

	wrapped := endpoint.Chain(nopMW, nopMW, nopMW, nopMW, nopMW, nopMW, nopMW, nopMW, nopMW, nopMW)(ep)
	ctx := context.Background()

	for range b.N {
		_, _ = wrapped(ctx, 0)
	}
}

func BenchmarkNop(b *testing.B) {
	ep := endpoint.Nop[int, int]()
	ctx := context.Background()

	for range b.N {
		_, _ = ep(ctx, 0)
	}
}
