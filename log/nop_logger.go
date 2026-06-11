package log

// NopLogger is a Logger that discards all messages. Safe for concurrent use.
type NopLogger struct{}

// NewNopLogger returns a Logger that discards all messages.
func NewNopLogger() Logger {
	return &NopLogger{}
}

// Log discards the log message.
func (n *NopLogger) Log(level Level, msg string, keysAndValues ...any) {}

// With returns the same NopLogger (no allocation).
func (n *NopLogger) With(keysAndValues ...any) Logger {
	return n
}
