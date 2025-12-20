package github

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestParseRateLimitHeaders(t *testing.T) {
	tests := []struct {
		name          string
		headers       http.Header
		wantRemaining int
		wantReset     int64
		wantErr       bool
	}{
		{
			name: "valid headers",
			headers: func() http.Header {
				h := http.Header{}
				h.Set(RateLimitRemainingHeader, "4999")
				h.Set(RateLimitResetHeader, "1672531200")
				return h
			}(),
			wantRemaining: 4999,
			wantReset:     1672531200,
			wantErr:       false,
		},
		{
			name: "zero remaining",
			headers: func() http.Header {
				h := http.Header{}
				h.Set(RateLimitRemainingHeader, "0")
				h.Set(RateLimitResetHeader, "1672531200")
				return h
			}(),
			wantRemaining: 0,
			wantReset:     1672531200,
			wantErr:       false,
		},
		{
			name: "missing remaining header",
			headers: func() http.Header {
				h := http.Header{}
				h.Set(RateLimitResetHeader, "1672531200")
				return h
			}(),
			wantErr: true,
		},
		{
			name: "missing reset header",
			headers: func() http.Header {
				h := http.Header{}
				h.Set(RateLimitRemainingHeader, "4999")
				return h
			}(),
			wantErr: true,
		},
		{
			name:    "missing both headers",
			headers: http.Header{},
			wantErr: true,
		},
		{
			name: "invalid remaining value",
			headers: func() http.Header {
				h := http.Header{}
				h.Set(RateLimitRemainingHeader, "invalid")
				h.Set(RateLimitResetHeader, "1672531200")
				return h
			}(),
			wantErr: true,
		},
		{
			name: "invalid reset value",
			headers: func() http.Header {
				h := http.Header{}
				h.Set(RateLimitRemainingHeader, "4999")
				h.Set(RateLimitResetHeader, "invalid")
				return h
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ParseRateLimitHeaders(tt.headers)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if info.Remaining != tt.wantRemaining {
				t.Errorf("Expected Remaining %d, got %d", tt.wantRemaining, info.Remaining)
			}

			expectedReset := time.Unix(tt.wantReset, 0)
			if !info.Reset.Equal(expectedReset) {
				t.Errorf("Expected Reset %v, got %v", expectedReset, info.Reset)
			}
		})
	}
}

func TestIsRateLimited(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{
			name:       "429 Too Many Requests",
			statusCode: http.StatusTooManyRequests,
			want:       true,
		},
		{
			name:       "403 Forbidden",
			statusCode: http.StatusForbidden,
			want:       true,
		},
		{
			name:       "200 OK",
			statusCode: http.StatusOK,
			want:       false,
		},
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			want:       false,
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
			}

			got := IsRateLimited(resp)
			if got != tt.want {
				t.Errorf("IsRateLimited() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateBackoff(t *testing.T) {
	tests := []struct {
		name    string
		attempt int
		want    time.Duration
	}{
		{
			name:    "first attempt",
			attempt: 0,
			want:    1 * time.Second, // 1 * 2^0
		},
		{
			name:    "second attempt",
			attempt: 1,
			want:    2 * time.Second, // 1 * 2^1
		},
		{
			name:    "third attempt",
			attempt: 2,
			want:    4 * time.Second, // 1 * 2^2
		},
		{
			name:    "fourth attempt",
			attempt: 3,
			want:    8 * time.Second, // 1 * 2^3
		},
		{
			name:    "fifth attempt",
			attempt: 4,
			want:    16 * time.Second, // 1 * 2^4
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateBackoff(tt.attempt)
			if got != tt.want {
				t.Errorf("calculateBackoff(%d) = %v, want %v", tt.attempt, got, tt.want)
			}
		})
	}
}

func TestHandleRateLimitMaxAttempts(t *testing.T) {
	ctx := context.Background()
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{},
	}

	err := HandleRateLimit(ctx, resp, MaxBackoffAttempts)
	if err == nil {
		t.Error("Expected error when max attempts reached, got nil")
	}

	if err != nil && err.Error() != "max retry attempts (5) reached for rate limiting" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestHandleRateLimitContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	headers := http.Header{}
	headers.Set(RateLimitRemainingHeader, "0")
	headers.Set(RateLimitResetHeader, "9999999999")
	
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     headers,
	}

	err := HandleRateLimit(ctx, resp, 0)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

func TestHandleRateLimitWithHeaders(t *testing.T) {
	ctx := context.Background()

	// Set reset time to 1 second in the future
	resetTime := time.Now().Add(1 * time.Second).Unix()

	headers := http.Header{}
	headers.Set(RateLimitRemainingHeader, "0")
	headers.Set(RateLimitResetHeader, string(rune(resetTime)))
	
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     headers,
	}

	start := time.Now()
	err := HandleRateLimit(ctx, resp, 0)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should wait for at least 1 second (plus buffer)
	// But we'll be lenient in the test due to timing variations
	if elapsed < 500*time.Millisecond {
		t.Errorf("Expected to wait, but elapsed time was only %v", elapsed)
	}
}

func TestHandleRateLimitExponentialBackoff(t *testing.T) {
	ctx := context.Background()

	// Response without rate limit headers - should use exponential backoff
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{},
	}

	start := time.Now()
	err := HandleRateLimit(ctx, resp, 1) // Second attempt (2^1 = 2 seconds)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should wait for approximately 2 seconds
	if elapsed < 1500*time.Millisecond || elapsed > 2500*time.Millisecond {
		t.Errorf("Expected ~2s backoff, got %v", elapsed)
	}
}

func TestWaitForSearchAPI(t *testing.T) {
	ctx := context.Background()

	start := time.Now()
	err := WaitForSearchAPI(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should wait for SearchAPIDelay (2 seconds)
	if elapsed < 1900*time.Millisecond || elapsed > 2100*time.Millisecond {
		t.Errorf("Expected ~2s delay, got %v", elapsed)
	}
}

func TestWaitForSearchAPIContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := WaitForSearchAPI(ctx)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

func TestCheckRateLimit(t *testing.T) {
	tests := []struct {
		name      string
		info      *RateLimitInfo
		threshold int
		want      bool
	}{
		{
			name: "above threshold",
			info: &RateLimitInfo{
				Remaining: 100,
				Reset:     time.Now().Add(1 * time.Hour),
			},
			threshold: 50,
			want:      false,
		},
		{
			name: "at threshold",
			info: &RateLimitInfo{
				Remaining: 50,
				Reset:     time.Now().Add(1 * time.Hour),
			},
			threshold: 50,
			want:      true,
		},
		{
			name: "below threshold",
			info: &RateLimitInfo{
				Remaining: 10,
				Reset:     time.Now().Add(1 * time.Hour),
			},
			threshold: 50,
			want:      true,
		},
		{
			name:      "nil info",
			info:      nil,
			threshold: 50,
			want:      false,
		},
		{
			name: "zero remaining",
			info: &RateLimitInfo{
				Remaining: 0,
				Reset:     time.Now().Add(1 * time.Hour),
			},
			threshold: 10,
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckRateLimit(tt.info, tt.threshold)
			if got != tt.want {
				t.Errorf("CheckRateLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetResetTime(t *testing.T) {
	tests := []struct {
		name string
		info *RateLimitInfo
		want string
	}{
		{
			name: "valid info",
			info: &RateLimitInfo{
				Remaining: 100,
				Reset:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			want: "2023-01-01T12:00:00Z",
		},
		{
			name: "nil info",
			info: nil,
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetResetTime(tt.info)
			if got != tt.want {
				t.Errorf("GetResetTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{
			name:       "429 - rate limited",
			statusCode: http.StatusTooManyRequests,
			want:       true,
		},
		{
			name:       "403 - forbidden (rate limit)",
			statusCode: http.StatusForbidden,
			want:       true,
		},
		{
			name:       "500 - internal server error",
			statusCode: http.StatusInternalServerError,
			want:       true,
		},
		{
			name:       "502 - bad gateway",
			statusCode: http.StatusBadGateway,
			want:       true,
		},
		{
			name:       "503 - service unavailable",
			statusCode: http.StatusServiceUnavailable,
			want:       true,
		},
		{
			name:       "504 - gateway timeout",
			statusCode: http.StatusGatewayTimeout,
			want:       true,
		},
		{
			name:       "200 - OK",
			statusCode: http.StatusOK,
			want:       false,
		},
		{
			name:       "404 - not found",
			statusCode: http.StatusNotFound,
			want:       false,
		},
		{
			name:       "400 - bad request",
			statusCode: http.StatusBadRequest,
			want:       false,
		},
		{
			name:       "401 - unauthorized",
			statusCode: http.StatusUnauthorized,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
			}

			got := ShouldRetry(resp)
			if got != tt.want {
				t.Errorf("ShouldRetry() = %v, want %v for status %d", got, tt.want, tt.statusCode)
			}
		})
	}
}

func TestShouldRetryNilResponse(t *testing.T) {
	got := ShouldRetry(nil)
	if got != false {
		t.Error("ShouldRetry(nil) should return false")
	}
}

func TestRateLimitInfoStruct(t *testing.T) {
	resetTime := time.Now().Add(1 * time.Hour)
	info := &RateLimitInfo{
		Remaining: 42,
		Reset:     resetTime,
	}

	if info.Remaining != 42 {
		t.Errorf("Expected Remaining 42, got %d", info.Remaining)
	}

	if !info.Reset.Equal(resetTime) {
		t.Errorf("Expected Reset %v, got %v", resetTime, info.Reset)
	}
}

func TestHandleRateLimitWithPastResetTime(t *testing.T) {
	ctx := context.Background()

	// Set reset time to the past
	resetTime := time.Now().Add(-1 * time.Hour).Unix()

	headers := http.Header{}
	headers.Set(RateLimitRemainingHeader, "0")
	headers.Set(RateLimitResetHeader, string(rune(resetTime)))
	
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     headers,
	}

	start := time.Now()
	err := HandleRateLimit(ctx, resp, 0)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should not wait since reset time is in the past
	// Falls back to exponential backoff
	if elapsed < 900*time.Millisecond || elapsed > 1100*time.Millisecond {
		t.Errorf("Expected ~1s backoff, got %v", elapsed)
	}
}

func TestWaitForSearchAPITimeout(t *testing.T) {
	// Create a context that times out before SearchAPIDelay
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := WaitForSearchAPI(ctx)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}

func TestHandleRateLimitTimeout(t *testing.T) {
	// Create a context that times out quickly
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	headers := http.Header{}
	headers.Set(RateLimitRemainingHeader, "0")
	headers.Set(RateLimitResetHeader, "9999999999")
	
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     headers,
	}

	err := HandleRateLimit(ctx, resp, 0)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}
