package log_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	harmonilog "github.com/harmonikit/harmoni/log"
)

func TestSlogLogger_DefaultLevel(t *testing.T) {
	// toSlogLevel with unknown value should map to Info (default).
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := harmonilog.NewSlogLogger(slog.New(handler))

	// Level(99) falls through to default → slog.LevelInfo.
	logger.Log(harmonilog.Level(99), "unknown level")

	out := buf.String()
	if out == "" {
		t.Fatal("expected log output for unknown level")
	}
}

func TestNopLogger_Direct(t *testing.T) {
	// Call NopLogger methods directly (not via interface) for coverage.
	nop := &harmonilog.NopLogger{}
	nop.Log(harmonilog.LevelInfo, "msg", "k", "v")
	nop.Log(harmonilog.LevelError, "also ignored")

	// With returns the same NopLogger.
	derived := nop.With("ctx", "val")
	if derived != nop {
		t.Error("NopLogger.With should return self")
	}

	_ = context.Background()
}
