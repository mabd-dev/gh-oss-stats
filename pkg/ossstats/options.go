package ossstats

import (
	"log"
	"net/http"
	"time"
)

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithToken sets the GitHub Personal Access Token for authentication.
// Required for reasonable rate limits (5,000/hour vs 60/hour unauthenticated).
func WithToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

// WithLOC enables or disables fetching lines of code metrics (additions/deletions).
// Default: false
func WithLOC(enabled bool) Option {
	return func(c *Client) {
		c.includeLOC = enabled
	}
}

// WithPRDetails enables or disables including detailed PR information.
// When enabled, includes a list of individual PR details for each contribution.
// Default: false
func WithPRDetails(enabled bool) Option {
	return func(c *Client) {
		c.includePRDetails = enabled
	}
}

// WithMinStars filters repositories by minimum star count.
// Only contributions to repositories with at least this many stars will be included.
// Default: 0 (no filtering)
func WithMinStars(stars int) Option {
	return func(c *Client) {
		c.minStars = stars
	}
}

// WithMaxPRs limits the maximum number of PRs to fetch.
// Useful for large contributors to avoid excessive API calls.
// Default: 500
func WithMaxPRs(max int) Option {
	return func(c *Client) {
		c.maxPRs = max
	}
}

// WithTimeout sets the overall timeout for the entire operation.
// Default: 5 minutes
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithExcludeOrgs excludes contributions to repositories owned by the specified organizations.
// This is useful for excluding your own organizations from the report.
func WithExcludeOrgs(orgs []string) Option {
	return func(c *Client) {
		c.excludeOrgs = orgs
	}
}

// WithLogger sets a custom logger for the client.
// The logger will receive informational messages about the operation progress.
// Default: no-op logger that discards all messages
func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithHTTPClient sets a custom HTTP client.
// Useful for testing or custom transport configuration.
// Default: http.DefaultClient with timeout
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithVerbose enables verbose logging to the default logger.
// This is a convenience option that sets up a standard logger.
func WithVerbose() Option {
	return func(c *Client) {
		c.logger = log.Default()
	}
}

// WithDebug enable/disable debug mode. When enabled, fake api client is used
// and mock api data is returned
func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}
