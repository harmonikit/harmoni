package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/harmonikit/harmoni/auth"
	"github.com/harmonikit/harmoni/endpoint"
)

// testAuth implements auth.Auth[int].
type testAuth struct {
	user string
	err  error
}

func (a *testAuth) Authenticate(ctx context.Context, req int) (context.Context, error) {
	if a.err != nil {
		return ctx, a.err
	}
	return auth.SetAuth(ctx, a.user), nil
}

func TestSetAuth_GetAuth(t *testing.T) {
	ctx := context.Background()
	ctx = auth.SetAuth(ctx, "user-123")

	v, ok := auth.GetAuth(ctx)
	if !ok {
		t.Fatal("expected auth value to be present")
	}
	if v != "user-123" {
		t.Fatalf("got %q, want %q", v, "user-123")
	}
}

func TestGetAuth_NotSet(t *testing.T) {
	ctx := context.Background()

	v, ok := auth.GetAuth(ctx)
	if ok {
		t.Fatal("expected auth value to be absent")
	}
	if v != nil {
		t.Fatalf("got %v, want nil", v)
	}
}

func TestClearAuth(t *testing.T) {
	ctx := context.Background()
	ctx = auth.SetAuth(ctx, "user-123")

	// Verify it's set.
	if v, ok := auth.GetAuth(ctx); !ok || v != "user-123" {
		t.Fatal("auth should be set")
	}

	ctx = auth.ClearAuth(ctx)

	v, ok := auth.GetAuth(ctx)
	if ok {
		t.Fatal("auth should be cleared")
	}
	if v != nil {
		t.Fatalf("got %v, want nil", v)
	}
}

func TestAuthMiddleware_Success(t *testing.T) {
	ta := &testAuth{user: "alice"}
	mw := auth.Middleware[int, string](ta)

	ep := endpoint.Endpoint[int, string](func(ctx context.Context, req int) (string, error) {
		v, ok := auth.GetAuth(ctx)
		if !ok {
			t.Fatal("expected auth value in context")
		}
		return v.(string), nil
	})

	wrapped := mw(ep)
	resp, err := wrapped(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "alice" {
		t.Fatalf("got %q, want %q", resp, "alice")
	}
}

func TestAuthMiddleware_Error(t *testing.T) {
	wantErr := errors.New("unauthorized")
	ta := &testAuth{err: wantErr}
	mw := auth.Middleware[int, string](ta)

	called := false
	ep := endpoint.Endpoint[int, string](func(ctx context.Context, req int) (string, error) {
		called = true
		return "ok", nil
	})

	wrapped := mw(ep)
	_, err := wrapped(context.Background(), 1)
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if called {
		t.Fatal("endpoint should not be called on auth failure")
	}
}

func TestAuthMiddleware_PreservesContext(t *testing.T) {
	// Verify that the auth middleware doesn't lose existing context values.
	ctx := context.WithValue(context.Background(), contextKey("existing"), "value")
	ta := &testAuth{user: "bob"}
	mw := auth.Middleware[int, string](ta)

	ep := endpoint.Endpoint[int, string](func(ctx context.Context, req int) (string, error) {
		// Both the existing value and the auth value should be present.
		if v := ctx.Value(contextKey("existing")); v != "value" {
			t.Fatalf("existing context value lost: got %v", v)
		}
		v, ok := auth.GetAuth(ctx)
		if !ok || v != "bob" {
			t.Fatalf("auth value missing: ok=%v v=%v", ok, v)
		}
		return "ok", nil
	})

	wrapped := mw(ep)
	_, err := wrapped(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

type contextKey string
