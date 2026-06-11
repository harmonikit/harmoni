package middleware_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/harmonikit/harmoni/endpoint"
	"github.com/harmonikit/harmoni/middleware"
)

var errTest = errors.New("test error")

func epOK(ctx context.Context, req int) (int, error) {
	return req + 1, nil
}

func epFail(ctx context.Context, req int) (int, error) {
	return 0, errTest
}

// go test -bench=. -benchtime=100000x works with Go 1.22+ for range b.N

func BenchmarkTimeout_NoDeadlineExceeded(b *testing.B) {
	ep := endpoint.Endpoint[int, int](epOK)
	wrapped := middleware.Timeout[int, int](5 * time.Second)(ep)
	ctx := context.Background()

	for range b.N {
		_, _ = wrapped(ctx, 0)
	}
}

func BenchmarkRetry_SuccessPath(b *testing.B) {
	ep := endpoint.Endpoint[int, int](epOK)
	backoff := func(attempt int) time.Duration { return 0 }
	wrapped := middleware.Retry[int, int](3, backoff)(ep)
	ctx := context.Background()

	for range b.N {
		_, _ = wrapped(ctx, 0)
	}
}

func BenchmarkRetry_FailurePath(b *testing.B) {
	ep := endpoint.Endpoint[int, int](epFail)
	backoff := func(attempt int) time.Duration { return 0 }
	wrapped := middleware.Retry[int, int](3, backoff)(ep)
	ctx := context.Background()

	for range b.N {
		_, _ = wrapped(ctx, 0)
	}
}

func BenchmarkRecovery_NoPanic(b *testing.B) {
	ep := endpoint.Endpoint[int, int](epOK)
	wrapped := middleware.Recovery[int, int]()(ep)
	ctx := context.Background()

	for range b.N {
		_, _ = wrapped(ctx, 0)
	}
}
