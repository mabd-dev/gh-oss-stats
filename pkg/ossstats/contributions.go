package ossstats

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/mabd-dev/gh-oss-stats/internal/github"
)

// GetContributions fetches and aggregates a user's open source contributions
// to external repositories (repos they don't own).
//
// Returns Stats containing the aggregated contribution data, or an error.
// If rate limiting occurs mid-fetch, returns ErrPartialResults with whatever
// data was collected before the rate limit.
func (c *Client) GetContributions(ctx context.Context, username string) (*Stats, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	c.logger.Printf("Fetching contributions for user: %s", username)

	// Initialize GitHub API client
	var apiClient github.GithubAPI
	if c.debug {
		c.logger.Printf("DEBUG MODE: Using mock API client")
		apiClient = github.NewMockAPIClient()
	} else {
		apiClient = github.NewAPIClient(c.httpClient, c.token)
	}

	// Step 1: Search for merged PRs to external repos
	c.logger.Printf("Searching for merged PRs...")
	issues, err := c.searchMergedPRs(ctx, apiClient, username)
	if err != nil {
		return nil, err
	}

	if len(issues) == 0 {
		c.logger.Printf("No contributions found")
		return &Stats{
			Username:      username,
			GeneratedAt:   time.Now().UTC(),
			Summary:       Summary{},
			Contributions: []Contribution{},
		}, nil
	}

	c.logger.Printf("Found %d merged PRs", len(issues))

	// Step 2: Fetch PR details and aggregate by repository
	c.logger.Printf("Fetching PR details...")
	contributions, errors := c.fetchPRDetails(ctx, apiClient, issues)

	// Step 3: Fetch repository metadata
	c.logger.Printf("Fetching repository metadata...")
	contributions = c.enrichWithRepoData(ctx, apiClient, contributions)

	// Step 4: Apply filters
	contributions = c.applyFilters(contributions)

	slices.SortFunc(contributions, func(a, b Contribution) int {
		if a.FirstContribution.Before(b.FirstContribution) {
			return 1
		}
		if a.FirstContribution.After(b.FirstContribution) {
			return -1
		}
		return 0
	})

	// Step 5: Calculate summary
	summary := c.calculateSummary(contributions)

	stats := &Stats{
		Username:      username,
		GeneratedAt:   time.Now().UTC(),
		Summary:       summary,
		Contributions: contributions,
	}

	// If there were errors during fetching, return partial results
	if len(errors) > 0 {
		c.logger.Printf("Completed with %d errors", len(errors))
		return stats, &ErrPartialResults{
			Stats:   stats,
			Errors:  errors,
			Message: fmt.Sprintf("collected %d contributions with errors", len(contributions)),
		}
	}

	c.logger.Printf("Successfully fetched %d contributions", len(contributions))
	return stats, nil
}

// searchMergedPRs searches for all merged PRs authored by the user to external repos.
func (c *Client) searchMergedPRs(ctx context.Context, api github.GithubAPI, username string) ([]github.Issue, error) {
	// Build search query: merged PRs by user, excluding their own repos
	query := fmt.Sprintf("author:%s type:pr is:merged -user:%s", username, username)

	// Exclude specified organizations
	for _, org := range c.excludeOrgs {
		if org != "" {
			query += fmt.Sprintf(" -org:%s", org)
		}
	}

	var allIssues []github.Issue
	page := 1
	perPage := 100

	for {
		// Respect search API rate limits
		if page > 1 {
			if err := github.WaitForSearchAPI(ctx); err != nil {
				return nil, fmt.Errorf("waiting for search API: %w", err)
			}
		}

		result, resp, err := api.SearchIssues(ctx, query, page, perPage)
		if err != nil {
			if github.IsRateLimited(resp) {
				resetTime := time.Now().Add(time.Minute)
				if info, err := github.ParseRateLimitHeaders(resp.Header); err == nil {
					resetTime = info.Reset
				}
				return nil, &ErrRateLimited{
					ResetAt: resetTime,
					Message: "search API rate limit exceeded",
				}
			}
			if resp != nil && resp.StatusCode == http.StatusUnauthorized {
				return nil, &ErrAuthentication{Message: "invalid or missing token"}
			}
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return nil, &ErrNotFound{Username: username}
			}
			return nil, fmt.Errorf("searching issues: %w", err)
		}

		allIssues = append(allIssues, result.Items...)

		// Check if we've hit the max PRs limit
		if c.maxPRs > 0 && len(allIssues) >= c.maxPRs {
			allIssues = allIssues[:c.maxPRs]
			c.logger.Printf("Reached max PRs limit (%d)", c.maxPRs)
			break
		}

		// Check if there are more pages
		if len(result.Items) < perPage {
			break
		}

		page++
	}

	return allIssues, nil
}

// fetchPRDetails fetches detailed information for each PR and aggregates by repository.
func (c *Client) fetchPRDetails(ctx context.Context, api github.GithubAPI, issues []github.Issue) ([]Contribution, []error) {
	// Map to aggregate PRs by repository
	repoMap := make(map[string]*Contribution)
	var mu sync.Mutex
	var errors []error

	// Process PRs with limited concurrency
	semaphore := make(chan struct{}, 5) // Limit to 5 concurrent requests
	var wg sync.WaitGroup

	for _, issue := range issues {
		// Skip if not a PR
		if issue.PullRequest == nil {
			continue
		}

		// Skip if not merged
		if issue.PullRequest.MergedAt == nil {
			continue
		}

		wg.Add(1)
		go func(iss github.Issue) {
			defer wg.Done()

			// Acquire semaphore
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				return
			}

			// Parse repository URL
			owner, repo, err := github.ParseRepoURL(iss.RepositoryURL)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("parsing repo URL: %w", err))
				mu.Unlock()
				return
			}

			// Fetch PR details if LOC is enabled
			var additions, deletions, commits int
			if c.includeLOC {
				pr, resp, err := api.GetPullRequest(ctx, owner, repo, iss.Number)
				if err != nil {
					if !github.IsRateLimited(resp) {
						mu.Lock()
						errors = append(errors, fmt.Errorf("fetching PR %s/%s#%d: %w", owner, repo, iss.Number, err))
						mu.Unlock()
					}
					return
				}
				additions = pr.Additions
				deletions = pr.Deletions
				commits = pr.Commits
			} else {
				commits = 1 // Default to 1 commit per PR if not fetching details
			}

			// Aggregate by repository
			repoKey := owner + "/" + repo
			mu.Lock()
			defer mu.Unlock()

			if contrib, exists := repoMap[repoKey]; exists {
				// Update existing contribution
				contrib.PRsMerged++
				contrib.Commits += commits
				contrib.Additions += additions
				contrib.Deletions += deletions

				// Update first/last contribution times
				mergedAt := *iss.PullRequest.MergedAt
				if mergedAt.Before(contrib.FirstContribution) {
					contrib.FirstContribution = mergedAt
				}
				if mergedAt.After(contrib.LastContribution) {
					contrib.LastContribution = mergedAt
				}
			} else {
				// Create new contribution entry
				repoMap[repoKey] = &Contribution{
					Repo:              repoKey,
					Owner:             owner,
					RepoName:          repo,
					PRsMerged:         1,
					Commits:           commits,
					Additions:         additions,
					Deletions:         deletions,
					FirstContribution: *iss.PullRequest.MergedAt,
					LastContribution:  *iss.PullRequest.MergedAt,
				}
			}
		}(issue)
	}

	wg.Wait()

	// Convert map to slice
	contributions := make([]Contribution, 0, len(repoMap))
	for _, contrib := range repoMap {
		contributions = append(contributions, *contrib)
	}

	return contributions, errors
}

// enrichWithRepoData fetches repository metadata and enriches contributions.
func (c *Client) enrichWithRepoData(ctx context.Context, api github.GithubAPI, contributions []Contribution) []Contribution {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit concurrent requests

	for i := range contributions {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			// Acquire semaphore
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				return
			}

			contrib := &contributions[idx]
			repo, _, err := api.GetRepository(ctx, contrib.Owner, contrib.RepoName)
			if err != nil {
				c.logger.Printf("Failed to fetch repo %s: %v", contrib.Repo, err)
				return
			}

			contrib.Description = repo.Description
			contrib.RepoURL = repo.HTMLURL
			contrib.Stars = repo.StargazersCount
		}(i)
	}

	wg.Wait()
	return contributions
}

// applyFilters applies client filters to contributions.
func (c *Client) applyFilters(contributions []Contribution) []Contribution {
	if c.minStars == 0 {
		return contributions
	}

	filtered := make([]Contribution, 0, len(contributions))
	for _, contrib := range contributions {
		if contrib.Stars >= c.minStars {
			filtered = append(filtered, contrib)
		}
	}

	return filtered
}

// calculateSummary calculates aggregate statistics.
func (c *Client) calculateSummary(contributions []Contribution) Summary {
	summary := Summary{
		TotalProjects: len(contributions),
	}

	for _, contrib := range contributions {
		summary.TotalPRsMerged += contrib.PRsMerged
		summary.TotalCommits += contrib.Commits
		summary.TotalAdditions += contrib.Additions
		summary.TotalDeletions += contrib.Deletions
	}

	return summary
}
