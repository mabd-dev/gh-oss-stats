package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// MockAPIClient is a mock GitHub API client that reads from local JSON files.
// Used for testing and development without hitting the real GitHub API.
type MockAPIClient struct {
	mockDataDir string
}

// NewMockAPIClient creates a new mock API client.
func NewMockAPIClient() *MockAPIClient {
	// Get the path to the mock data directory
	// Assuming we're running from the project root
	mockDataDir := filepath.Join("internal", "github", "mockResponses")

	return &MockAPIClient{
		mockDataDir: mockDataDir,
	}
}

// readMockFile reads and unmarshals a mock JSON file.
func (c *MockAPIClient) readMockFile(filename string, result interface{}) error {
	path := filepath.Join(c.mockDataDir, filename)

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading mock file %s: %w", filename, err)
	}

	if err := json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("unmarshaling mock file %s: %w", filename, err)
	}

	return nil
}

// SearchIssues returns mock search results for merged PRs.
func (c *MockAPIClient) SearchIssues(ctx context.Context, query string, page, perPage int) (*SearchIssuesResponse, *http.Response, error) {
	var result SearchIssuesResponse

	if err := c.readMockFile("searchForMergedPrs.json", &result); err != nil {
		return nil, nil, err
	}

	// Create a mock response
	mockResp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
	}
	mockResp.Header.Set("X-RateLimit-Remaining", "5000")
	mockResp.Header.Set("X-RateLimit-Limit", "5000")

	return &result, mockResp, nil
}

// GetPullRequest returns mock PR details.
func (c *MockAPIClient) GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, *http.Response, error) {
	var result PullRequest

	if err := c.readMockFile("pullRequest.json", &result); err != nil {
		return nil, nil, err
	}

	// Create a mock response
	mockResp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
	}
	mockResp.Header.Set("X-RateLimit-Remaining", "5000")
	mockResp.Header.Set("X-RateLimit-Limit", "5000")

	return &result, mockResp, nil
}

// GetRepository returns mock repository information.
func (c *MockAPIClient) GetRepository(ctx context.Context, owner, repo string) (*Repository, *http.Response, error) {
	var result Repository

	if err := c.readMockFile("repository.json", &result); err != nil {
		return nil, nil, err
	}

	// Create a mock response
	mockResp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
	}
	mockResp.Header.Set("X-RateLimit-Remaining", "5000")
	mockResp.Header.Set("X-RateLimit-Limit", "5000")

	return &result, mockResp, nil
}

// GetRateLimit returns mock rate limit information.
func (c *MockAPIClient) GetRateLimit(ctx context.Context) (*RateLimitResponse, error) {
	return &RateLimitResponse{
		Resources: RateLimitResources{
			Core: RateLimit{
				Limit:     5000,
				Remaining: 5000,
				Reset:     1735689600, // Some future timestamp
				Used:      0,
			},
			Search: RateLimit{
				Limit:     30,
				Remaining: 30,
				Reset:     1735689600,
				Used:      0,
			},
		},
	}, nil
}
