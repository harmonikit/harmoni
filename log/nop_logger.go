package log

// NopLogger is a Logger that discards all messages. Safe for concurrent use.
type NopLogger struct{}

// NewNopLogger returns a Logger that discards all messages.
func NewNopLogger() Logger {
	return &NopLogger{}
}

// Log discards the log message.
func (n *NopLogger) Log(_ Level, _ string, _ ...any) {}

// With returns the same NopLogger (no allocation).
func (n *NopLogger) With(_ ...any) Logger {
	return n
}
