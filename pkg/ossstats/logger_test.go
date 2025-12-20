package ossstats

import (
	"testing"
)

// mockLogger is a test logger that captures log messages
type mockLogger struct {
	messages []string
}

func (m *mockLogger) Printf(format string, v ...any) {
	// Store the formatted message
	m.messages = append(m.messages, format)
}

func TestLoggerInterface(t *testing.T) {
	// Verify mockLogger implements Logger
	var _ Logger = &mockLogger{}
}

func TestDefaultLogger(t *testing.T) {
	logger := defaultLogger{}

	// Should not panic when called
	logger.Printf("test message")
	logger.Printf("test with args: %s %d", "hello", 42)
}

func TestDefaultLoggerImplementsInterface(t *testing.T) {
	// Verify defaultLogger implements Logger
	var _ Logger = defaultLogger{}
}

func TestMockLogger(t *testing.T) {
	mock := &mockLogger{}

	mock.Printf("first message")
	mock.Printf("second message with %s", "args")

	if len(mock.messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(mock.messages))
	}

	if mock.messages[0] != "first message" {
		t.Errorf("First message = %q, want %q", mock.messages[0], "first message")
	}

	if mock.messages[1] != "second message with %s" {
		t.Errorf("Second message = %q, want %q", mock.messages[1], "second message with %s")
	}
}
