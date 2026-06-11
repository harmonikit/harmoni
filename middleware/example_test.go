package middleware_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/harmonikit/harmoni/endpoint"
	"github.com/harmonikit/harmoni/middleware"
)

func ExampleTimeout() {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req * 2, nil
	})

	withTimeout := middleware.Timeout[int, int](1 * time.Second)(ep)
	resp, err := withTimeout(context.Background(), 21)
	if err != nil {
		fmt.Println("timeout:", err)
	} else {
		fmt.Println(resp)
	}
	// Output: 42
}

func ExampleRetry() {
	var calls int
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		calls++
		if calls < 3 {
			return 0, errors.New("transient error")
		}
		return req, nil
	})

	backoff := func(attempt int) time.Duration { return 1 * time.Millisecond }
	withRetry := middleware.Retry[int, int](3, backoff)(ep)
	resp, err := withRetry(context.Background(), 42)
	if err != nil {
		fmt.Println("failed:", err)
	} else {
		fmt.Printf("resp=%d calls=%d\n", resp, calls)
	}
	// Output: resp=42 calls=3
}

func ExampleRecovery() {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		if req < 0 {
			panic("negative input")
		}
		return req * 2, nil
	})

	withRecovery := middleware.Recovery[int, int]()(ep)
	_, err := withRecovery(context.Background(), -1)
	fmt.Println(err)
	// Output: panic recovered: negative input
}
