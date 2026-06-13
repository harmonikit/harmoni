package middleware_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/harmonikit/harmoni/endpoint"
	"github.com/harmonikit/harmoni/middleware"
)

func TestTimeout_Success(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
	})

	wrapped := middleware.Timeout[int, int](5 * time.Second)(ep)
	resp, err := wrapped(context.Background(), 41)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestTimeout_Exceeded(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return req + 1, nil
		}
	})

	wrapped := middleware.Timeout[int, int](1 * time.Nanosecond)(ep)
	_, err := wrapped(context.Background(), 41)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got error %v, want context.DeadlineExceeded", err)
	}
}

func TestTimeout_SetsDeadline(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		if _, ok := ctx.Deadline(); !ok {
			t.Error("context should have a deadline")
		}
		return req + 1, nil
	})

	wrapped := middleware.Timeout[int, int](5 * time.Second)(ep)
	resp, err := wrapped(context.Background(), 41)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestTimeout_ExpiredContext(t *testing.T) {
	called := false
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		called = true
		return req + 1, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	wrapped := middleware.Timeout[int, int](5 * time.Second)(ep)
	_, err := wrapped(ctx, 41)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got error %v, want context.Canceled", err)
	}
	if called {
		t.Fatal("endpoint should not have been called")
	}
}

func TestRetry_Success_FirstAttempt(t *testing.T) {
	var attempts int
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		attempts++
		return req + 1, nil
	})

	backoff := func(attempt int) time.Duration { return 1 * time.Millisecond }
	wrapped := middleware.Retry[int, int](3, backoff)(ep)

	resp, err := wrapped(context.Background(), 41)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
	if attempts != 1 {
		t.Fatalf("got %d attempts, want 1", attempts)
	}
}

func TestRetry_Success_AfterRetries(t *testing.T) {
	var attempts int
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		attempts++
		if attempts < 3 {
			return 0, errors.New("temporary failure")
		}
		return req + 1, nil
	})

	backoff := func(attempt int) time.Duration { return 1 * time.Millisecond }
	wrapped := middleware.Retry[int, int](5, backoff)(ep)

	resp, err := wrapped(context.Background(), 41)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
	if attempts != 3 {
		t.Fatalf("got %d attempts, want 3", attempts)
	}
}

func TestRetry_AllFailed(t *testing.T) {
	wantErr := errors.New("permanent failure")
	var attempts int
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		attempts++
		return 0, wantErr
	})

	backoff := func(attempt int) time.Duration { return 1 * time.Millisecond }
	wrapped := middleware.Retry[int, int](3, backoff)(ep)

	_, err := wrapped(context.Background(), 41)
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if attempts != 4 {
		t.Fatalf("got %d attempts, want 4", attempts)
	}
}

func TestRetry_ZeroMax(t *testing.T) {
	var attempts int
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		attempts++
		return 0, errors.New("fail")
	})

	backoff := func(attempt int) time.Duration { return 1 * time.Millisecond }
	wrapped := middleware.Retry[int, int](0, backoff)(ep)

	_, err := wrapped(context.Background(), 41)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if attempts != 1 {
		t.Fatalf("got %d attempts, want 1", attempts)
	}
}

func TestRetry_NegativeMax(t *testing.T) {
	// Negative max means no attempts at all — returns zero values.
	var called bool
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		called = true
		return req, nil
	})

	backoff := func(attempt int) time.Duration { return 0 }
	wrapped := middleware.Retry[int, int](-1, backoff)(ep)

	resp, err := wrapped(context.Background(), 41)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 0 {
		t.Fatalf("got %d, want 0 (zero value)", resp)
	}
	if called {
		t.Fatal("endpoint should not have been called with negative max")
	}
}

func TestRetry_ZeroBackoff(t *testing.T) {
	// Zero backoff — no sleep between retries, just immediate retry.
	var attempts int
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		attempts++
		if attempts < 3 {
			return 0, errors.New("fail")
		}
		return req, nil
	})

	backoff := func(attempt int) time.Duration { return 0 }
	wrapped := middleware.Retry[int, int](5, backoff)(ep)

	resp, err := wrapped(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestRetry_ContextCancelled(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return 0, errors.New("fail")
	})

	ctx, cancel := context.WithCancel(context.Background())

	backoff := func(attempt int) time.Duration {
		cancel() // Cancel on first retry.
		return 10 * time.Millisecond
	}

	wrapped := middleware.Retry[int, int](5, backoff)(ep)
	_, err := wrapped(ctx, 41)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got error %v, want context.Canceled", err)
	}
}

func TestRecovery_NoPanic(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req * 2, nil
	})

	wrapped := middleware.Recovery[int, int]()(ep)
	resp, err := wrapped(context.Background(), 21)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestRecovery_Panic_String(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		panic("something went wrong")
	})

	wrapped := middleware.Recovery[int, int]()(ep)
	resp, err := wrapped(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error from panic, got nil")
	}
	if !strings.Contains(err.Error(), "panic recovered") {
		t.Fatalf("got error %v, want panic recovered error", err)
	}
	if resp != 0 {
		t.Fatalf("got %d, want zero value", resp)
	}
}

func TestRecovery_Panic_Error(t *testing.T) {
	wantErr := errors.New("custom panic")
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		panic(wantErr)
	})

	wrapped := middleware.Recovery[int, int]()(ep)
	_, err := wrapped(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error from panic, got nil")
	}
	if !strings.Contains(err.Error(), wantErr.Error()) {
		t.Fatalf("got error %v, want error containing %v", err, wantErr)
	}
}

func TestRecovery_Panic_Int(t *testing.T) {
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		panic(42)
	})

	wrapped := middleware.Recovery[int, int]()(ep)
	_, err := wrapped(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error from panic, got nil")
	}
	if !strings.Contains(err.Error(), "42") {
		t.Fatalf("got error %v, want error containing 42", err)
	}
}

func TestRecovery_StructType(t *testing.T) {
	type Req struct{ Name string }
	type Resp struct{ ID int }

	ep := endpoint.Endpoint[Req, Resp](func(ctx context.Context, req Req) (Resp, error) {
		return Resp{ID: 1}, nil
	})

	wrapped := middleware.Recovery[Req, Resp]()(ep)
	resp, err := wrapped(context.Background(), Req{Name: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != 1 {
		t.Fatalf("got %d, want 1", resp.ID)
	}
}

func TestRecovery_StructPanic_ZeroResp(t *testing.T) {
	type Req struct{ Name string }
	type Resp struct {
		ID   int
		Name string
	}

	ep := endpoint.Endpoint[Req, Resp](func(ctx context.Context, req Req) (Resp, error) {
		panic("boom")
	})

	wrapped := middleware.Recovery[Req, Resp]()(ep)
	resp, err := wrapped(context.Background(), Req{Name: "test"})
	if err == nil {
		t.Fatal("expected error from panic, got nil")
	}
	if resp.ID != 0 || resp.Name != "" {
		t.Fatalf("got %+v, want zero Resp", resp)
	}
}
