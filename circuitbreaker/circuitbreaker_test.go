package circuitbreaker_test

import (
	"context"
	"errors"
	"testing"

	"github.com/harmonikit/harmoni/circuitbreaker"
	"github.com/harmonikit/harmoni/endpoint"
)

func TestState_String(t *testing.T) {
	tests := []struct {
		state circuitbreaker.State
		want  string
	}{
		{circuitbreaker.StateClosed, "closed"},
		{circuitbreaker.StateHalfOpen, "half-open"},
		{circuitbreaker.StateOpen, "open"},
		{circuitbreaker.State(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNopCircuitBreaker_Execute(t *testing.T) {
	cb := circuitbreaker.NewNopCircuitBreaker[int, int]()

	ep := endpoint.Endpoint[int, int](func(ctx context.Context, req int) (int, error) {
		return req * 2, nil
	})

	resp, err := cb.Execute(context.Background(), 21, ep)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Fatalf("got %d, want 42", resp)
	}
}

func TestNopCircuitBreaker_Execute_Error(t *testing.T) {
	cb := circuitbreaker.NewNopCircuitBreaker[int, string]()

	wantErr := errors.New("downstream failure")
	ep := endpoint.Endpoint[int, string](func(ctx context.Context, req int) (string, error) {
		return "", wantErr
	})

	_, err := cb.Execute(context.Background(), 1, ep)
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
}

func TestNopCircuitBreaker_State(t *testing.T) {
	cb := circuitbreaker.NewNopCircuitBreaker[int, int]()

	if state := cb.State(); state != circuitbreaker.StateClosed {
		t.Fatalf("got %v, want StateClosed", state)
	}
}

func TestCircuitBreaker_Interface(t *testing.T) {
	var _ circuitbreaker.CircuitBreaker[int, int] = circuitbreaker.NewNopCircuitBreaker[int, int]()
}
