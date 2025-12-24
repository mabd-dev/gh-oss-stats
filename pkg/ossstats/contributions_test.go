package ossstats

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mabd-dev/gh-oss-stats/internal/github"
)

func TestGetContributionsNoContributions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/search/issues") {
			// Return empty search results
			resp := github.SearchIssuesResponse{
				TotalCount: 0,
				Items:      []github.Issue{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	httpClient := &http.Client{}
	client := New(
		WithHTTPClient(httpClient),
		WithToken("test-token"),
	)

	// Override the base URL in the client's HTTP client
	client.httpClient.Transport = &mockTransport{server: server}

	stats, err := client.GetContributions(context.Background(), "testuser")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if stats == nil {
		t.Fatal("stats should not be nil")
	}

	if stats.Username != "testuser" {
		t.Errorf("Username = %s, want testuser", stats.Username)
	}

	if len(stats.Contributions) != 0 {
		t.Errorf("Contributions count = %d, want 0", len(stats.Contributions))
	}

	if stats.Summary.TotalProjects != 0 {
		t.Errorf("TotalProjects = %d, want 0", stats.Summary.TotalProjects)
	}
}

func TestGetContributionsWithContributions(t *testing.T) {
	mergedAt := time.Now().UTC()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/search/issues") {
			resp := github.SearchIssuesResponse{
				TotalCount: 1,
				Items: []github.Issue{
					{
						Number:        123,
						Title:         "Test PR",
						State:         "closed",
						RepositoryURL: "https://api.github.com/repos/owner/repo",
						PullRequest: &github.PullRequestRef{
							MergedAt: &mergedAt,
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		if strings.Contains(r.URL.Path, "/repos/owner/repo/pulls") {
			resp := github.PullRequest{
				Number:    123,
				Merged:    true,
				Commits:   5,
				Additions: 100,
				Deletions: 20,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/repos/owner/repo") {
			resp := github.Repository{
				Name:            "repo",
				FullName:        "owner/repo",
				Description:     "Test repository",
				HTMLURL:         "https://github.com/owner/repo",
				StargazersCount: 100,
				Owner: github.User{
					Login: "owner",
					ID:    123,
					Type:  "User",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client := New(
		WithToken("test-token"),
		WithLOC(true),
	)
	client.httpClient.Transport = &mockTransport{server: server}

	stats, err := client.GetContributions(context.Background(), "testuser")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if stats == nil {
		t.Fatal("stats should not be nil")
	}

	if len(stats.Contributions) != 1 {
		t.Fatalf("Contributions count = %d, want 1", len(stats.Contributions))
	}

	contrib := stats.Contributions[0]
	if contrib.Repo != "owner/repo" {
		t.Errorf("Repo = %s, want owner/repo", contrib.Repo)
	}

	if contrib.Stars != 100 {
		t.Errorf("Stars = %d, want 100", contrib.Stars)
	}

	if contrib.Commits != 5 {
		t.Errorf("Commits = %d, want 5", contrib.Commits)
	}

	if contrib.Additions != 100 {
		t.Errorf("Additions = %d, want 100", contrib.Additions)
	}

	if contrib.Deletions != 20 {
		t.Errorf("Deletions = %d, want 20", contrib.Deletions)
	}

	// Verify summary
	if stats.Summary.TotalProjects != 1 {
		t.Errorf("TotalProjects = %d, want 1", stats.Summary.TotalProjects)
	}

	if stats.Summary.TotalPRsMerged != 1 {
		t.Errorf("TotalPRsMerged = %d, want 1", stats.Summary.TotalPRsMerged)
	}

	if stats.Summary.TotalCommits != 5 {
		t.Errorf("TotalCommits = %d, want 5", stats.Summary.TotalCommits)
	}

	if stats.Summary.TotalAdditions != 100 {
		t.Errorf("TotalAdditions = %d, want 100", stats.Summary.TotalAdditions)
	}

	if stats.Summary.TotalDeletions != 20 {
		t.Errorf("TotalDeletions = %d, want 20", stats.Summary.TotalDeletions)
	}
}

func TestGetContributionsWithoutLOC(t *testing.T) {
	mergedAt := time.Now().UTC()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/search/issues") {
			resp := github.SearchIssuesResponse{
				TotalCount: 1,
				Items: []github.Issue{
					{
						Number:        123,
						RepositoryURL: "https://api.github.com/repos/owner/repo",
						PullRequest: &github.PullRequestRef{
							MergedAt: &mergedAt,
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/repos/owner/repo") && !strings.Contains(r.URL.Path, "/pulls") {
			resp := github.Repository{
				Name:            "repo",
				FullName:        "owner/repo",
				StargazersCount: 50,
				Owner: github.User{
					Login: "owner",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Should NOT call /pulls endpoint when LOC is disabled
		if strings.Contains(r.URL.Path, "/pulls") {
			t.Error("Should not fetch PR details when LOC is disabled")
			http.NotFound(w, r)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client := New(
		WithToken("test-token"),
		WithLOC(false), // Disable LOC fetching
	)
	client.httpClient.Transport = &mockTransport{server: server}

	stats, err := client.GetContributions(context.Background(), "testuser")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(stats.Contributions) != 1 {
		t.Fatalf("Contributions count = %d, want 1", len(stats.Contributions))
	}

	contrib := stats.Contributions[0]

	// Should default to 1 commit when not fetching details
	if contrib.Commits != 1 {
		t.Errorf("Commits = %d, want 1 (default)", contrib.Commits)
	}

	// LOC should be 0 when not fetched
	if contrib.Additions != 0 {
		t.Errorf("Additions = %d, want 0", contrib.Additions)
	}

	if contrib.Deletions != 0 {
		t.Errorf("Deletions = %d, want 0", contrib.Deletions)
	}
}

func TestGetContributionsWithMinStarsFilter(t *testing.T) {
	mergedAt := time.Now().UTC()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/search/issues") {
			resp := github.SearchIssuesResponse{
				TotalCount: 2,
				Items: []github.Issue{
					{
						Number:        1,
						RepositoryURL: "https://api.github.com/repos/popular/repo",
						PullRequest:   &github.PullRequestRef{MergedAt: &mergedAt},
					},
					{
						Number:        2,
						RepositoryURL: "https://api.github.com/repos/unpopular/repo",
						PullRequest:   &github.PullRequestRef{MergedAt: &mergedAt},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		if strings.Contains(r.URL.Path, "/repos/popular/repo") && !strings.Contains(r.URL.Path, "/pulls") {
			resp := github.Repository{
				Name:            "repo",
				FullName:        "popular/repo",
				StargazersCount: 1000, // Above filter
				Owner:           github.User{Login: "popular"},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		if strings.Contains(r.URL.Path, "/repos/unpopular/repo") && !strings.Contains(r.URL.Path, "/pulls") {
			resp := github.Repository{
				Name:            "repo",
				FullName:        "unpopular/repo",
				StargazersCount: 10, // Below filter
				Owner:           github.User{Login: "unpopular"},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client := New(
		WithToken("test-token"),
		WithMinStars(100), // Filter repos with < 100 stars
	)
	client.httpClient.Transport = &mockTransport{server: server}

	stats, err := client.GetContributions(context.Background(), "testuser")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should only include the repo with >= 100 stars
	if len(stats.Contributions) != 1 {
		t.Fatalf("Contributions count = %d, want 1 (filtered)", len(stats.Contributions))
	}

	if stats.Contributions[0].Repo != "popular/repo" {
		t.Errorf("Repo = %s, want popular/repo", stats.Contributions[0].Repo)
	}

	if stats.Contributions[0].Stars != 1000 {
		t.Errorf("Stars = %d, want 1000", stats.Contributions[0].Stars)
	}
}

func TestGetContributionsRateLimited(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/search/issues") {
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", "9999999999")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message":"API rate limit exceeded"}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client := New(
		WithToken("test-token"),
	)
	client.httpClient.Transport = &mockTransport{server: server}

	_, err := client.GetContributions(context.Background(), "testuser")

	if err == nil {
		t.Fatal("Expected rate limit error, got nil")
	}

	rateLimitErr, ok := err.(*ErrRateLimited)
	if !ok {
		t.Fatalf("Expected *ErrRateLimited, got %T", err)
	}

	if !contains(rateLimitErr.Message, "rate limit") {
		t.Errorf("Error message = %q, want to contain 'rate limit'", rateLimitErr.Message)
	}
}

func TestGetContributionsUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/search/issues") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message":"Bad credentials"}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client := New(
		WithToken("invalid-token"),
	)
	client.httpClient.Transport = &mockTransport{server: server}

	_, err := client.GetContributions(context.Background(), "testuser")

	if err == nil {
		t.Fatal("Expected authentication error, got nil")
	}

	authErr, ok := err.(*ErrAuthentication)
	if !ok {
		t.Fatalf("Expected *ErrAuthentication, got %T", err)
	}

	if !contains(authErr.Message, "token") {
		t.Errorf("Error message = %q, want to contain 'token'", authErr.Message)
	}
}

func TestGetContributionsUserNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/search/issues") {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message":"Not Found"}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client := New(
		WithToken("test-token"),
	)
	client.httpClient.Transport = &mockTransport{server: server}

	_, err := client.GetContributions(context.Background(), "nonexistent")

	if err == nil {
		t.Fatal("Expected not found error, got nil")
	}

	notFoundErr, ok := err.(*ErrNotFound)
	if !ok {
		t.Fatalf("Expected *ErrNotFound, got %T", err)
	}

	if notFoundErr.Username != "nonexistent" {
		t.Errorf("Username = %q, want nonexistent", notFoundErr.Username)
	}
}

func TestCalculateSummary(t *testing.T) {
	client := New()

	contributions := []Contribution{
		{
			PRsMerged: 5,
			Commits:   15,
			Additions: 100,
			Deletions: 20,
		},
		{
			PRsMerged: 3,
			Commits:   10,
			Additions: 50,
			Deletions: 10,
		},
	}

	summary := client.calculateSummary(contributions)

	if summary.TotalProjects != 2 {
		t.Errorf("TotalProjects = %d, want 2", summary.TotalProjects)
	}

	if summary.TotalPRsMerged != 8 {
		t.Errorf("TotalPRsMerged = %d, want 8", summary.TotalPRsMerged)
	}

	if summary.TotalCommits != 25 {
		t.Errorf("TotalCommits = %d, want 25", summary.TotalCommits)
	}

	if summary.TotalAdditions != 150 {
		t.Errorf("TotalAdditions = %d, want 150", summary.TotalAdditions)
	}

	if summary.TotalDeletions != 30 {
		t.Errorf("TotalDeletions = %d, want 30", summary.TotalDeletions)
	}
}

func TestApplyFilters(t *testing.T) {
	tests := []struct {
		name          string
		minStars      int
		contributions []Contribution
		wantCount     int
	}{
		{
			name:     "no filter",
			minStars: 0,
			contributions: []Contribution{
				{Repo: "a/b", Stars: 10},
				{Repo: "c/d", Stars: 100},
			},
			wantCount: 2,
		},
		{
			name:     "filter by 50 stars",
			minStars: 50,
			contributions: []Contribution{
				{Repo: "a/b", Stars: 10},
				{Repo: "c/d", Stars: 100},
				{Repo: "e/f", Stars: 50},
			},
			wantCount: 2, // Only 100 and 50
		},
		{
			name:     "filter removes all",
			minStars: 1000,
			contributions: []Contribution{
				{Repo: "a/b", Stars: 10},
				{Repo: "c/d", Stars: 100},
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(WithMinStars(tt.minStars))

			filtered := client.applyFilters(tt.contributions)

			if len(filtered) != tt.wantCount {
				t.Errorf("filtered count = %d, want %d", len(filtered), tt.wantCount)
			}

			// Verify all filtered items meet the criteria
			for _, contrib := range filtered {
				if contrib.Stars < tt.minStars {
					t.Errorf("Contribution %s has %d stars, below minimum %d", contrib.Repo, contrib.Stars, tt.minStars)
				}
			}
		})
	}
}

func TestGetContributionsMultiplePRsSameRepo(t *testing.T) {
	mergedAt1 := time.Now().UTC().Add(-48 * time.Hour)
	mergedAt2 := time.Now().UTC().Add(-24 * time.Hour)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/search/issues") {
			resp := github.SearchIssuesResponse{
				TotalCount: 2,
				Items: []github.Issue{
					{
						Number:        1,
						RepositoryURL: "https://api.github.com/repos/owner/repo",
						PullRequest:   &github.PullRequestRef{MergedAt: &mergedAt1},
					},
					{
						Number:        2,
						RepositoryURL: "https://api.github.com/repos/owner/repo",
						PullRequest:   &github.PullRequestRef{MergedAt: &mergedAt2},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/repos/owner/repo") && !strings.Contains(r.URL.Path, "/pulls") {
			resp := github.Repository{
				Name:            "repo",
				FullName:        "owner/repo",
				StargazersCount: 100,
				Owner:           github.User{Login: "owner"},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client := New(
		WithToken("test-token"),
		WithLOC(false),
	)
	client.httpClient.Transport = &mockTransport{server: server}

	stats, err := client.GetContributions(context.Background(), "testuser")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should aggregate into single contribution
	if len(stats.Contributions) != 1 {
		t.Fatalf("Contributions count = %d, want 1 (aggregated)", len(stats.Contributions))
	}

	contrib := stats.Contributions[0]

	if contrib.PRsMerged != 2 {
		t.Errorf("PRsMerged = %d, want 2", contrib.PRsMerged)
	}

	if contrib.Commits != 2 {
		t.Errorf("Commits = %d, want 2 (1 per PR)", contrib.Commits)
	}

	// Should have first and last contribution times
	if !contrib.FirstContribution.Equal(mergedAt1) {
		t.Errorf("FirstContribution = %v, want %v", contrib.FirstContribution, mergedAt1)
	}

	if !contrib.LastContribution.Equal(mergedAt2) {
		t.Errorf("LastContribution = %v, want %v", contrib.LastContribution, mergedAt2)
	}
}

// mockTransport redirects requests to test server
type mockTransport struct {
	server *httptest.Server
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite the URL to point to our test server
	req.URL.Scheme = "http"
	req.URL.Host = strings.TrimPrefix(t.server.URL, "http://")

	return http.DefaultTransport.RoundTrip(req)
}
