package github

import (
	"context"
	"net/http"
)

// This allows for both real and mock implementations.
// GithubAPI defines the interface for GitHub API operations.
type GithubAPI interface {
	// SearchIssues searches for issues/PRs matching the given query.
	SearchIssues(ctx context.Context, query string, page, perPage int) (*SearchIssuesResponse, *http.Response, error)

	// GetPullRequest fetches detailed information about a pull request.
	GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, *http.Response, error)

	// GetRepository fetches information about a repository.
	GetRepository(ctx context.Context, owner, repo string) (*Repository, *http.Response, error)

	// GetRateLimit fetches the current rate limit status.
	GetRateLimit(ctx context.Context) (*RateLimitResponse, error)
}
