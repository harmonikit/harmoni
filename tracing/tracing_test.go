package tracing_test

import (
	"context"
	"errors"
	"testing"

	"github.com/harmonikit/harmoni/tracing"
)

func TestNopSpan_End(t *testing.T) {
	span := tracing.NopSpan{}
	span.End() // should not panic
}

func TestNopSpan_SetAttributes(t *testing.T) {
	span := tracing.NopSpan{}
	span.SetAttributes("key", "value", "count", 42) // should not panic
}

func TestNopTracer_Start(t *testing.T) {
	tracer := tracing.NewNopTracer[int, string]()

	ctx := context.Background()
	newCtx, span := tracer.Start(ctx, "operation", 42)

	if newCtx != ctx {
		t.Error("NopTracer should return the same context")
	}
	if _, ok := span.(tracing.NopSpan); !ok {
		t.Errorf("expected NopSpan, got %T", span)
	}
}

func TestNopTracer_End(t *testing.T) {
	tracer := tracing.NewNopTracer[int, string]()

	ctx := context.Background()
	_, span := tracer.Start(ctx, "operation", 42)

	tracer.End(ctx, span, "response", nil) // should not panic
}

func TestNopTracer_End_WithError(t *testing.T) {
	tracer := tracing.NewNopTracer[int, string]()

	ctx := context.Background()
	_, span := tracer.Start(ctx, "operation", 42)

	err := errors.New("something failed")
	tracer.End(ctx, span, "partial response", err) // should not panic
}

func TestNopTracer_StructTypes(t *testing.T) {
	type Req struct{ Name string }
	type Resp struct{ ID int }

	tracer := tracing.NewNopTracer[Req, Resp]()
	ctx := context.Background()
	_, span := tracer.Start(ctx, "createUser", Req{Name: "test"})
	tracer.End(ctx, span, Resp{ID: 1}, nil)
}

func TestTracer_Interface(t *testing.T) {
	var _ tracing.Tracer[int, string] = tracing.NewNopTracer[int, string]()
}

func TestSpan_Interface(t *testing.T) {
	var _ tracing.Span = tracing.NopSpan{}
}
