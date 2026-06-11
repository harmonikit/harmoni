package endpoint_test

import (
	"context"
	"errors"
	"testing"

	"github.com/harmonikit/harmoni/endpoint"
)

func TestEndpoint_Nop(t *testing.T) {
	nop := endpoint.Nop[string, int]()

	resp, err := nop(context.Background(), "ignored")
	if err != nil {
		t.Fatalf("Nop should not return error, got: %v", err)
	}
	if resp != 0 {
		t.Fatalf("Nop should return zero value, got: %d", resp)
	}
}

func TestEndpoint_Execute(t *testing.T) {
	wantResp := 42

	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req * 2, nil
	})

	resp, err := ep(context.Background(), 21)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != wantResp {
		t.Fatalf("got %d, want %d", resp, wantResp)
	}
}

func TestEndpoint_Error(t *testing.T) {
	wantErr := errors.New("something went wrong")

	ep := endpoint.Endpoint[int, string](func(ctx context.Context, req int) (string, error) {
		return "", wantErr
	})

	_, err := ep(context.Background(), 1)
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
}

func TestMiddleware_Identity(t *testing.T) {
	// Identity middleware: passes through unchanged.
	identity := endpoint.Middleware[int, int](func(next endpoint.Endpoint[int, int]) endpoint.Endpoint[int, int] {
		return next
	})

	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req + 1, nil
	})

	wrapped := identity(ep)
	resp, err := wrapped(context.Background(), 41)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestMiddleware_Transform(t *testing.T) {
	// Middleware that doubles the request before passing it on.
	doubleReq := endpoint.Middleware[int, int](func(next endpoint.Endpoint[int, int]) endpoint.Endpoint[int, int] {
		return func(ctx context.Context, req int) (int, error) {
			return next(ctx, req*2)
		}
	})

	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	wrapped := doubleReq(ep)
	resp, err := wrapped(context.Background(), 21)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestChain_Single(t *testing.T) {
	// A single middleware in Chain is equivalent to applying it directly.
	mw := endpoint.Middleware[int, int](func(next endpoint.Endpoint[int, int]) endpoint.Endpoint[int, int] {
		return func(ctx context.Context, req int) (int, error) {
			resp, err := next(ctx, req)
			return resp + 1, err
		}
	})

	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	chained := endpoint.Chain(mw)(ep)
	resp, err := chained(context.Background(), 41)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestChain_Ordering(t *testing.T) {
	// Chain(A, B, C) should apply A outermost, then B, then C.
	// So A receives the request first, then B, then C, then the endpoint.
	// On the way back, C's response goes to B, then A.
	//
	// We track order with a slice to verify left-to-right wrapping.
	var order []string

	makeMW := func(name string) endpoint.Middleware[int, int] {
		return func(next endpoint.Endpoint[int, int]) endpoint.Endpoint[int, int] {
			return func(ctx context.Context, req int) (int, error) {
				order = append(order, name+":before")
				resp, err := next(ctx, req)
				order = append(order, name+":after")
				return resp, err
			}
		}
	}

	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		order = append(order, "endpoint")
		return req, nil
	})

	wrapped := endpoint.Chain(
		makeMW("outer"),
		makeMW("middle"),
		makeMW("inner"),
	)(ep)

	_, err := wrapped(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Expected: outer:before, middle:before, inner:before, endpoint,
	//           inner:after, middle:after, outer:after
	want := []string{
		"outer:before", "middle:before", "inner:before",
		"endpoint",
		"inner:after", "middle:after", "outer:after",
	}

	if len(order) != len(want) {
		t.Fatalf("got order %v (len %d), want %v (len %d)", order, len(order), want, len(want))
	}
	for i := range order {
		if order[i] != want[i] {
			t.Fatalf("at position %d: got %q, want %q\nfull order: %v", i, order[i], want[i], order)
		}
	}
}

func TestChain_EmptyRest(t *testing.T) {
	// Chain with no rest middlewares should just apply the outer one.
	called := false
	mw := endpoint.Middleware[int, int](func(next endpoint.Endpoint[int, int]) endpoint.Endpoint[int, int] {
		called = true
		return next
	})

	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req, nil
	})

	_ = endpoint.Chain(mw)(ep)
	if !called {
		t.Fatal("middleware was not applied")
	}
}

func TestNop_DifferentTypes(t *testing.T) {
	// Nop should work with any type parameters.
	t.Run("string", func(t *testing.T) {
		nop := endpoint.Nop[string, string]()
		resp, err := nop(context.Background(), "hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp != "" {
			t.Fatalf("got %q, want empty string", resp)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type Req struct{ Name string }
		type Resp struct{ ID int }
		nop := endpoint.Nop[Req, Resp]()
		resp, err := nop(context.Background(), Req{Name: "test"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID != 0 {
			t.Fatalf("got %d, want 0", resp.ID)
		}
	})

	t.Run("pointer", func(t *testing.T) {
		nop := endpoint.Nop[*int, *string]()
		resp, err := nop(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp != nil {
			t.Fatalf("got %v, want nil", resp)
		}
	})
}

func TestEndpoint_ContextPropagation(t *testing.T) {
	type ctxKey string
	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		v := ctx.Value(ctxKey("key")).(string)
		if v != "value" {
			t.Fatalf("got %q, want %q", v, "value")
		}
		return req, nil
	})

	ctx := context.WithValue(context.Background(), ctxKey("key"), "value")
	_, err := ep(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
