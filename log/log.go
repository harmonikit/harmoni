// Package log defines a structured logging interface for harmoni services.
//
// It provides a Logger interface that is library-owned (not tied to slog.Handler),
// allowing it to evolve independently of the standard library. A SlogLogger
// adapter wraps *slog.Logger so users who just want stdlib logging get it for free.
//
// Example:
//
//	logger := log.NewSlogLogger(slog.Default())
//	logger.Info("request processed", "method", "GET", "duration", 12*time.Millisecond)
package log

// Level represents the severity of a log message.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// String returns the level as a human-readable string.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger is a structured logging interface. Implementations must be safe for
// concurrent use.
type Logger interface {
	// Log emits a log message at the given level with structured key-value pairs.
	// keysAndValues must be an even number of arguments, alternating key, value.
	Log(level Level, msg string, keysAndValues ...any)

	// With returns a new Logger with the given key-value pairs added to the
	// logging context. The returned logger derives from the receiver.
	With(keysAndValues ...any) Logger
}
