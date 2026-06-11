package log_test

import (
	"bytes"
	"fmt"
	"log/slog"

	harmonilog "github.com/harmonikit/harmoni/log"
)

func ExampleSlogLogger() {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	logger := harmonilog.NewSlogLogger(slog.New(handler))

	logger.Log(harmonilog.LevelInfo, "server started", "port", 8080)
	fmt.Print(buf.String())
	// Output: level=INFO msg="server started" port=8080
}

func ExampleNopLogger() {
	logger := harmonilog.NewNopLogger()

	// All calls are no-ops — useful in tests and benchmarks.
	logger.Log(harmonilog.LevelError, "this will not appear anywhere")
	derived := logger.With("component", "test")
	derived.Log(harmonilog.LevelDebug, "also silent")
	// Output:
}

func ExampleLogger_With() {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})
	base := harmonilog.NewSlogLogger(slog.New(handler))

	serverLogger := base.With("component", "server")
	serverLogger.Log(harmonilog.LevelInfo, "listening", "addr", ":8080")
	fmt.Print(buf.String())
	// Output: level=INFO msg=listening component=server addr=:8080
}
