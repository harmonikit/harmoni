package tracing_test

import (
	"context"
	"testing"

	"github.com/harmonikit/harmoni/tracing"
)

func BenchmarkNopTracer_StartEnd(b *testing.B) {
	tracer := tracing.NewNopTracer[int, string]()
	ctx := context.Background()

	for range b.N {
		_, span := tracer.Start(ctx, "op", 0)
		tracer.End(ctx, span, "ok", nil)
	}
}
