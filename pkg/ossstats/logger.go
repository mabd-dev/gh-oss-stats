package ossstats

// Logger is the interface for logging within the library.
// Implementations can provide custom logging behavior.
type Logger interface {
	Printf(format string, v ...any)
}

// defaultLogger is a no-op logger that discards all log messages.
type defaultLogger struct{}

func (defaultLogger) Printf(format string, v ...any) {}
