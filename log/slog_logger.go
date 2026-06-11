package log

import (
	"context"
	"log/slog"
)

// SlogLogger adapts *slog.Logger to the Logger interface.
// It maps harmoni log levels to slog levels:
//
//	LevelDebug → slog.LevelDebug
//	LevelInfo  → slog.LevelInfo
//	LevelWarn  → slog.LevelWarn
//	LevelError → slog.LevelError
type SlogLogger struct {
	logger *slog.Logger
}

// NewSlogLogger wraps *slog.Logger as a harmoni Logger.
func NewSlogLogger(logger *slog.Logger) *SlogLogger {
	return &SlogLogger{logger: logger}
}

// Log emits a log message at the given level.
func (l *SlogLogger) Log(level Level, msg string, keysAndValues ...any) {
	l.logger.Log(context.Background(), toSlogLevel(level), msg, keysAndValues...)
}

// With returns a new Logger with additional context.
func (l *SlogLogger) With(keysAndValues ...any) Logger {
	return &SlogLogger{logger: l.logger.With(keysAndValues...)}
}

// toSlogLevel converts a harmoni Level to a slog.Level.
func toSlogLevel(l Level) slog.Level {
	switch l {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
