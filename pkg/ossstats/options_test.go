package ossstats

import (
	"log"
	"net/http"
	"testing"
	"time"
)

func TestWithToken(t *testing.T) {
	client := &Client{}
	token := "test-token-123"

	opt := WithToken(token)
	opt(client)

	if client.token != token {
		t.Errorf("token = %s, want %s", client.token, token)
	}
}

func TestWithLOC(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{"enabled", true},
		{"disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}

			opt := WithLOC(tt.enabled)
			opt(client)

			if client.includeLOC != tt.enabled {
				t.Errorf("includeLOC = %v, want %v", client.includeLOC, tt.enabled)
			}
		})
	}
}

func TestWithPRDetails(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{"enabled", true},
		{"disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}

			opt := WithPRDetails(tt.enabled)
			opt(client)

			if client.includePRDetails != tt.enabled {
				t.Errorf("includePRDetails = %v, want %v", client.includePRDetails, tt.enabled)
			}
		})
	}
}

func TestWithMinStars(t *testing.T) {
	tests := []struct {
		name  string
		stars int
	}{
		{"zero", 0},
		{"hundred", 100},
		{"thousand", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}

			opt := WithMinStars(tt.stars)
			opt(client)

			if client.minStars != tt.stars {
				t.Errorf("minStars = %d, want %d", client.minStars, tt.stars)
			}
		})
	}
}

func TestWithMaxPRs(t *testing.T) {
	tests := []struct {
		name string
		max  int
	}{
		{"default", 500},
		{"hundred", 100},
		{"unlimited", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}

			opt := WithMaxPRs(tt.max)
			opt(client)

			if client.maxPRs != tt.max {
				t.Errorf("maxPRs = %d, want %d", client.maxPRs, tt.max)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{"1 minute", 1 * time.Minute},
		{"5 minutes", 5 * time.Minute},
		{"10 minutes", 10 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{}

			opt := WithTimeout(tt.timeout)
			opt(client)

			if client.timeout != tt.timeout {
				t.Errorf("timeout = %v, want %v", client.timeout, tt.timeout)
			}
		})
	}
}

func TestWithLogger(t *testing.T) {
	client := &Client{}
	logger := &mockLogger{}

	opt := WithLogger(logger)
	opt(client)

	if client.logger != logger {
		t.Error("logger not set correctly")
	}

	// Test that logger works
	client.logger.Printf("test message")
	if len(logger.messages) != 1 {
		t.Errorf("logger messages count = %d, want 1", len(logger.messages))
	}
}

func TestWithHTTPClient(t *testing.T) {
	client := &Client{}
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	opt := WithHTTPClient(httpClient)
	opt(client)

	if client.httpClient != httpClient {
		t.Error("httpClient not set correctly")
	}

	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("httpClient.Timeout = %v, want %v", client.httpClient.Timeout, 30*time.Second)
	}
}

func TestWithVerbose(t *testing.T) {
	client := &Client{
		logger: defaultLogger{}, // Start with default
	}

	opt := WithVerbose()
	opt(client)

	// Should replace with log.Default()
	if client.logger == nil {
		t.Error("logger should not be nil")
	}

	// Verify it's not the default logger anymore
	_, isDefault := client.logger.(defaultLogger)
	if isDefault {
		t.Error("logger should not be defaultLogger after WithVerbose")
	}

	// Should be *log.Logger
	if _, ok := client.logger.(*log.Logger); !ok {
		t.Error("logger should be *log.Logger")
	}
}

func TestMultipleOptions(t *testing.T) {
	client := &Client{}

	token := "test-token"
	minStars := 100
	timeout := 10 * time.Minute
	logger := &mockLogger{}

	opts := []Option{
		WithToken(token),
		WithLOC(true),
		WithPRDetails(true),
		WithMinStars(minStars),
		WithTimeout(timeout),
		WithLogger(logger),
	}

	for _, opt := range opts {
		opt(client)
	}

	// Verify all options were applied
	if client.token != token {
		t.Errorf("token = %s, want %s", client.token, token)
	}

	if !client.includeLOC {
		t.Error("includeLOC should be true")
	}

	if !client.includePRDetails {
		t.Error("includePRDetails should be true")
	}

	if client.minStars != minStars {
		t.Errorf("minStars = %d, want %d", client.minStars, minStars)
	}

	if client.timeout != timeout {
		t.Errorf("timeout = %v, want %v", client.timeout, timeout)
	}

	if client.logger != logger {
		t.Error("logger not set correctly")
	}
}

func TestOptionOverride(t *testing.T) {
	client := &Client{}

	// Apply option twice
	WithMinStars(100)(client)
	WithMinStars(200)(client)

	// Second call should override
	if client.minStars != 200 {
		t.Errorf("minStars = %d, want 200 (should be overridden)", client.minStars)
	}
}
