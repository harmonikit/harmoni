package tracing_test

import (
	"context"
	"fmt"

	"github.com/harmonikit/harmoni/tracing"
)

func ExampleNopTracer() {
	tracer := tracing.NewNopTracer[int, string]()

	ctx := context.Background()
	_, span := tracer.Start(ctx, "greet", 42)
	defer tracer.End(ctx, span, "hello", nil)

	fmt.Println("operation traced")
	// Output: operation traced
}
