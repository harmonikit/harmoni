package tracing_test

import (
	"testing"

	"github.com/harmonikit/harmoni/tracing"
)

func TestNopSpan_Direct(t *testing.T) {
	s := tracing.NopSpan{}
	s.End()
	s.SetAttributes("k1", "v1", "k2", 42)
}

func TestNopTracer_End_Direct(t *testing.T) {
	tr := &tracing.NopTracer[int, string]{}
	tr.End(nil, tracing.NopSpan{}, "ok", nil)
}
