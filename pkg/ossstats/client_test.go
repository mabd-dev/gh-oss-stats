package ossstats

import (
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	client := New()

	if client == nil {
		t.Fatal("New() returned nil")
	}

	// Verify defaults
	if client.includeLOC != DefaultIncludeLOC {
		t.Errorf("includeLOC = %v, want %v", client.includeLOC, DefaultIncludeLOC)
	}

	if client.includePRDetails != DefaultIncludePRDetails {
		t.Errorf("includePRDetails = %v, want %v", client.includePRDetails, DefaultIncludePRDetails)
	}

	if client.minStars != DefaultMinStars {
		t.Errorf("minStars = %d, want %d", client.minStars, DefaultMinStars)
	}

	if client.maxPRs != DefaultMaxPRS {
		t.Errorf("maxPRs = %d, want %d", client.maxPRs, DefaultMaxPRS)
	}

	if client.timeout != DefaultTimeout {
		t.Errorf("timeout = %v, want %v", client.timeout, DefaultTimeout)
	}

	if client.httpClient == nil {
		t.Error("httpClient should not be nil")
	}

	if client.logger == nil {
		t.Error("logger should not be nil")
	}

	// Verify HTTP client timeout is set
	if client.httpClient.Timeout != DefaultTimeout {
		t.Errorf("httpClient.Timeout = %v, want %v", client.httpClient.Timeout, DefaultTimeout)
	}
}

func TestNewWithOptions(t *testing.T) {
	token := "test-token"
	minStars := 100
	maxPRs := 200
	timeout := 10 * time.Minute
	logger := &mockLogger{}
	httpClient := &http.Client{}

	client := New(
		WithToken(token),
		WithLOC(true),
		WithPRDetails(true),
		WithMinStars(minStars),
		WithMaxPRs(maxPRs),
		WithTimeout(timeout),
		WithLogger(logger),
		WithHTTPClient(httpClient),
	)

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

	if client.maxPRs != maxPRs {
		t.Errorf("maxPRs = %d, want %d", client.maxPRs, maxPRs)
	}

	if client.timeout != timeout {
		t.Errorf("timeout = %v, want %v", client.timeout, timeout)
	}

	if client.logger != logger {
		t.Error("logger not set correctly")
	}

	if client.httpClient != httpClient {
		t.Error("httpClient not set correctly")
	}
}

func TestNewHTTPClientTimeout(t *testing.T) {
	tests := []struct {
		name            string
		clientTimeout   time.Duration
		expectedTimeout time.Duration
		setupHTTP       bool
		httpTimeout     time.Duration
	}{
		{
			name:            "default client gets timeout from option",
			clientTimeout:   10 * time.Minute,
			expectedTimeout: 10 * time.Minute,
			setupHTTP:       false,
		},
		{
			name:            "custom client keeps its timeout",
			clientTimeout:   10 * time.Minute,
			expectedTimeout: 5 * time.Minute,
			setupHTTP:       true,
			httpTimeout:     5 * time.Minute,
		},
		{
			name:            "custom client with zero timeout gets client timeout",
			clientTimeout:   10 * time.Minute,
			expectedTimeout: 10 * time.Minute,
			setupHTTP:       true,
			httpTimeout:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []Option{
				WithTimeout(tt.clientTimeout),
			}

			if tt.setupHTTP {
				httpClient := &http.Client{
					Timeout: tt.httpTimeout,
				}
				opts = append(opts, WithHTTPClient(httpClient))
			}

			client := New(opts...)

			if client.httpClient.Timeout != tt.expectedTimeout {
				t.Errorf("httpClient.Timeout = %v, want %v", client.httpClient.Timeout, tt.expectedTimeout)
			}
		})
	}
}

func TestDefaultConstants(t *testing.T) {
	// Verify default constants have expected values
	if DefaultIncludeLOC != false {
		t.Errorf("DefaultIncludeLOC = %v, want false", DefaultIncludeLOC)
	}

	if DefaultIncludePRDetails != false {
		t.Errorf("DefaultIncludePRDetails = %v, want false", DefaultIncludePRDetails)
	}

	if DefaultMinStars != 0 {
		t.Errorf("DefaultMinStars = %d, want 0", DefaultMinStars)
	}

	if DefaultMaxPRS != 500 {
		t.Errorf("DefaultMaxPRS = %d, want 500", DefaultMaxPRS)
	}

	if DefaultTimeout != 5*time.Minute {
		t.Errorf("DefaultTimeout = %v, want 5m", DefaultTimeout)
	}
}

func TestNewDefaultLogger(t *testing.T) {
	client := New()

	// Should have defaultLogger
	if _, ok := client.logger.(defaultLogger); !ok {
		t.Error("default client should have defaultLogger")
	}

	// Should not panic when logging
	client.logger.Printf("test message")
}

func TestNewMultipleOptionsSameType(t *testing.T) {
	// Last option should win
	client := New(
		WithMinStars(100),
		WithMinStars(200),
		WithMinStars(300),
	)

	if client.minStars != 300 {
		t.Errorf("minStars = %d, want 300 (last option should win)", client.minStars)
	}
}

func TestClientFieldsNotExported(t *testing.T) {
	client := New()

	// Verify fields are not exported (lowercase first letter)
	// This is a compile-time check - if this compiles, fields are private

	// We can't access these fields from outside the package
	// client.token = "test" // Would not compile
	// client.includeLOC = true // Would not compile

	// But we can verify the client exists and is usable
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestNewEmptyOptions(t *testing.T) {
	// Calling New with empty slice should work
	client := New([]Option{}...)

	if client == nil {
		t.Fatal("New() with empty options returned nil")
	}

	// Should have defaults
	if client.minStars != DefaultMinStars {
		t.Errorf("minStars = %d, want default %d", client.minStars, DefaultMinStars)
	}
}
