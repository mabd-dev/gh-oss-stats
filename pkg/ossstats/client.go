package ossstats

import (
	"net/http"
	"time"
)

var (
	DefaultIncludeLOC       bool          = false
	DefaultIncludePRDetails bool          = false
	DefaultMinStars         int           = 0
	DefaultMaxPRS           int           = 500
	DefaultTimeout          time.Duration = 5 * time.Minute
)

// Client represents a GitHub OSS stats client.
// It is safe for concurrent use by multiple goroutines.
type Client struct {
	// Authentication
	token string

	// Configuration options
	includeLOC       bool
	includePRDetails bool
	minStars         int
	maxPRs           int
	timeout          time.Duration
	excludeOrgs      []string

	// HTTP client
	httpClient *http.Client

	// Logger
	logger Logger

	debug bool
}

// New creates a new Client with the provided options.
// The client is configured with sensible defaults that can be overridden
// using functional options.
//
// Example:
//
//	client := ossstats.New(
//	    ossstats.WithToken(token),
//	    ossstats.WithMinStars(100),
//	    ossstats.WithVerbose(),
//	)
func New(opts ...Option) *Client {
	// Create client with default values
	client := &Client{
		includeLOC:       DefaultIncludeLOC,
		includePRDetails: DefaultIncludePRDetails,
		minStars:         DefaultMinStars,
		maxPRs:           DefaultMaxPRS,
		timeout:          DefaultTimeout,
		httpClient:       &http.Client{},
		logger:           defaultLogger{},
	}

	// Apply all provided options
	for _, opt := range opts {
		opt(client)
	}

	// Configure HTTP client timeout if not already set
	if client.httpClient.Timeout == 0 {
		client.httpClient.Timeout = client.timeout
	}

	return client
}
