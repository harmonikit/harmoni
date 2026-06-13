package circuitbreaker_test

import (
	"context"
	"testing"

	"github.com/harmonikit/harmoni/circuitbreaker"
	"github.com/harmonikit/harmoni/endpoint"
)

func BenchmarkNopCircuitBreaker_Execute(b *testing.B) {
	cb := circuitbreaker.NewNopCircuitBreaker[int, int]()
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
	})
	ctx := context.Background()

	for range b.N {
		_, _ = cb.Execute(ctx, 0, ep)
	}
}
