package log_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	harmonilog "github.com/harmonikit/harmoni/log"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level harmonilog.Level
		want  string
	}{
		{harmonilog.LevelDebug, "DEBUG"},
		{harmonilog.LevelInfo, "INFO"},
		{harmonilog.LevelWarn, "WARN"},
		{harmonilog.LevelError, "ERROR"},
		{harmonilog.Level(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSlogLogger_Log(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := harmonilog.NewSlogLogger(slog.New(handler))

	logger.Log(harmonilog.LevelInfo, "hello", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "msg=hello") {
		t.Errorf("expected log output to contain msg=hello, got: %s", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("expected log output to contain key=value, got: %s", output)
	}
}

func TestSlogLogger_LevelMapping(t *testing.T) {
	tests := []struct {
		name  string
		level harmonilog.Level
		want  string
	}{
		{"debug", harmonilog.LevelDebug, "DEBUG"},
		{"info", harmonilog.LevelInfo, "INFO"},
		{"warn", harmonilog.LevelWarn, "WARN"},
		{"error", harmonilog.LevelError, "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
			logger := harmonilog.NewSlogLogger(slog.New(handler))

			logger.Log(tt.level, "test message")

			output := buf.String()
			if !strings.Contains(output, "level="+tt.want) {
				t.Errorf("expected level %s, got: %s", tt.want, output)
			}
		})
	}
}

func TestSlogLogger_With(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	base := harmonilog.NewSlogLogger(slog.New(handler))

	derived := base.With("component", "server")
	derived.Log(harmonilog.LevelInfo, "started")

	output := buf.String()
	if !strings.Contains(output, "component=server") {
		t.Errorf("expected component=server in output, got: %s", output)
	}
}

func TestSlogLogger_With_Multiple(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	base := harmonilog.NewSlogLogger(slog.New(handler))

	derived := base.With("a", "1").With("b", "2")
	derived.Log(harmonilog.LevelInfo, "test")

	output := buf.String()
	if !strings.Contains(output, "a=1") {
		t.Errorf("expected a=1 in output, got: %s", output)
	}
	if !strings.Contains(output, "b=2") {
		t.Errorf("expected b=2 in output, got: %s", output)
	}
}

func TestSlogLogger_BaseUnchanged(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	base := harmonilog.NewSlogLogger(slog.New(handler))

	// With returns a derived logger, base should be unchanged.
	_ = base.With("component", "server")
	base.Log(harmonilog.LevelInfo, "base message")

	output := buf.String()
	if strings.Contains(output, "component=server") {
		t.Errorf("base logger should not have component=server, got: %s", output)
	}
}

func TestNopLogger(t *testing.T) {
	logger := harmonilog.NewNopLogger()

	// Should not panic.
	logger.Log(harmonilog.LevelInfo, "ignored", "key", "value")
	logger.Log(harmonilog.LevelError, "also ignored")

	// With should return a logger that also doesn't panic.
	derived := logger.With("ctx", "value")
	derived.Log(harmonilog.LevelDebug, "nop derived")
}

func TestNopLogger_With_ReturnsNop(t *testing.T) {
	logger := harmonilog.NewNopLogger()
	derived := logger.With("key", "value")

	if _, ok := derived.(*harmonilog.NopLogger); !ok {
		t.Error("NopLogger.With should return a *NopLogger")
	}
}

func TestNopLogger_Log_Direct(t *testing.T) {
	// Direct call (not via interface) to ensure coverage.
	nop := &harmonilog.NopLogger{}
	nop.Log(harmonilog.LevelInfo, "ignored", "key", "value")
	nop.Log(harmonilog.LevelDebug, "also ignored")
}

func TestLogger_Interface(t *testing.T) {
	// Compile-time check: SlogLogger implements Logger.
	var _ harmonilog.Logger = harmonilog.NewSlogLogger(slog.Default())

	// Compile-time check: NopLogger implements Logger.
	var _ harmonilog.Logger = harmonilog.NewNopLogger()
}
