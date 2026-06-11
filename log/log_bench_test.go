package log_test

import (
	"log/slog"
	"testing"

	harmonilog "github.com/harmonikit/harmoni/log"
)

func BenchmarkSlogLogger_Log(b *testing.B) {
	logger := harmonilog.NewSlogLogger(slog.Default())

	for range b.N {
		logger.Log(harmonilog.LevelInfo, "benchmark message", "key", "value")
	}
}

func BenchmarkSlogLogger_With(b *testing.B) {
	base := harmonilog.NewSlogLogger(slog.Default())

	for range b.N {
		_ = base.With("key", "value")
	}
}

func BenchmarkNopLogger_Log(b *testing.B) {
	logger := harmonilog.NewNopLogger()

	for range b.N {
		logger.Log(harmonilog.LevelInfo, "benchmark message", "key", "value")
	}
}

func BenchmarkNopLogger_With(b *testing.B) {
	base := harmonilog.NewNopLogger()

	for range b.N {
		_ = base.With("key", "value")
	}
}
