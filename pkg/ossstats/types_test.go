package ossstats

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestStatsJSONMarshal(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	stats := Stats{
		Username:    "testuser",
		GeneratedAt: now,
		Summary: Summary{
			TotalProjects:  5,
			TotalPRsMerged: 10,
			TotalCommits:   25,
			TotalAdditions: 500,
			TotalDeletions: 200,
		},
		Contributions: []Contribution{
			{
				Repo:              "owner/repo",
				Owner:             "owner",
				RepoName:          "repo",
				Description:       "Test repository",
				RepoURL:           "https://github.com/owner/repo",
				Stars:             100,
				PRsMerged:         2,
				Commits:           5,
				Additions:         50,
				Deletions:         20,
				FirstContribution: now,
				LastContribution:  now,
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(stats)
	if err != nil {
		t.Fatalf("Failed to marshal Stats: %v", err)
	}

	// Unmarshal back
	var decoded Stats
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal Stats: %v", err)
	}

	// Verify
	if decoded.Username != stats.Username {
		t.Errorf("Username mismatch: got %s, want %s", decoded.Username, stats.Username)
	}

	if decoded.Summary.TotalProjects != stats.Summary.TotalProjects {
		t.Errorf("TotalProjects mismatch: got %d, want %d", decoded.Summary.TotalProjects, stats.Summary.TotalProjects)
	}

	if len(decoded.Contributions) != len(stats.Contributions) {
		t.Fatalf("Contributions length mismatch: got %d, want %d", len(decoded.Contributions), len(stats.Contributions))
	}

	if decoded.Contributions[0].Repo != stats.Contributions[0].Repo {
		t.Errorf("Repo mismatch: got %s, want %s", decoded.Contributions[0].Repo, stats.Contributions[0].Repo)
	}
}

func TestErrRateLimitedError(t *testing.T) {
	tests := []struct {
		name    string
		err     *ErrRateLimited
		want    string
		wantMsg string
	}{
		{
			name: "with message",
			err: &ErrRateLimited{
				ResetAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				Message: "search API rate limit exceeded",
			},
			wantMsg: "search API rate limit exceeded",
		},
		{
			name: "without message",
			err: &ErrRateLimited{
				ResetAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			wantMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()

			if got == "" {
				t.Error("Error() returned empty string")
			}

			if tt.wantMsg != "" && !contains(got, tt.wantMsg) {
				t.Errorf("Error() = %q, want to contain %q", got, tt.wantMsg)
			}

			if !contains(got, "2023-01-01T12:00:00Z") {
				t.Errorf("Error() = %q, want to contain timestamp", got)
			}
		})
	}
}

func TestErrAuthenticationError(t *testing.T) {
	tests := []struct {
		name    string
		err     *ErrAuthentication
		wantMsg string
	}{
		{
			name: "with message",
			err: &ErrAuthentication{
				Message: "invalid token",
			},
			wantMsg: "invalid token",
		},
		{
			name:    "without message",
			err:     &ErrAuthentication{},
			wantMsg: "authentication failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()

			if !contains(got, tt.wantMsg) {
				t.Errorf("Error() = %q, want to contain %q", got, tt.wantMsg)
			}

			if !contains(got, "authentication") {
				t.Errorf("Error() = %q, want to contain 'authentication'", got)
			}
		})
	}
}

func TestErrNotFoundError(t *testing.T) {
	err := &ErrNotFound{
		Username: "nonexistent",
	}

	got := err.Error()

	if !contains(got, "nonexistent") {
		t.Errorf("Error() = %q, want to contain username 'nonexistent'", got)
	}

	if !contains(got, "not found") {
		t.Errorf("Error() = %q, want to contain 'not found'", got)
	}
}

func TestErrPartialResultsError(t *testing.T) {
	stats := &Stats{
		Username:    "testuser",
		GeneratedAt: time.Now().UTC(),
		Summary: Summary{
			TotalProjects: 2,
		},
		Contributions: []Contribution{},
	}

	tests := []struct {
		name       string
		err        *ErrPartialResults
		wantMsg    string
		wantErrCnt int
	}{
		{
			name: "with message and errors",
			err: &ErrPartialResults{
				Stats:   stats,
				Errors:  []error{errors.New("error 1"), errors.New("error 2")},
				Message: "collected 2 contributions",
			},
			wantMsg:    "collected 2 contributions",
			wantErrCnt: 2,
		},
		{
			name: "without message",
			err: &ErrPartialResults{
				Stats:  stats,
				Errors: []error{errors.New("error 1")},
			},
			wantErrCnt: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()

			if !contains(got, "partial results") {
				t.Errorf("Error() = %q, want to contain 'partial results'", got)
			}

			if tt.wantMsg != "" && !contains(got, tt.wantMsg) {
				t.Errorf("Error() = %q, want to contain %q", got, tt.wantMsg)
			}

			if tt.err.Stats != stats {
				t.Error("Stats field not preserved")
			}

			if len(tt.err.Errors) != tt.wantErrCnt {
				t.Errorf("Errors count = %d, want %d", len(tt.err.Errors), tt.wantErrCnt)
			}
		})
	}
}

func TestSummaryZeroValues(t *testing.T) {
	summary := Summary{}

	if summary.TotalProjects != 0 {
		t.Errorf("TotalProjects = %d, want 0", summary.TotalProjects)
	}

	if summary.TotalPRsMerged != 0 {
		t.Errorf("TotalPRsMerged = %d, want 0", summary.TotalPRsMerged)
	}

	if summary.TotalCommits != 0 {
		t.Errorf("TotalCommits = %d, want 0", summary.TotalCommits)
	}

	if summary.TotalAdditions != 0 {
		t.Errorf("TotalAdditions = %d, want 0", summary.TotalAdditions)
	}

	if summary.TotalDeletions != 0 {
		t.Errorf("TotalDeletions = %d, want 0", summary.TotalDeletions)
	}
}

func TestContributionJSONMarshal(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	contrib := Contribution{
		Repo:              "owner/repo",
		Owner:             "owner",
		RepoName:          "repo",
		Description:       "A test repository",
		RepoURL:           "https://github.com/owner/repo",
		Stars:             1000,
		PRsMerged:         5,
		Commits:           15,
		Additions:         250,
		Deletions:         50,
		FirstContribution: now,
		LastContribution:  now.Add(24 * time.Hour),
	}

	// Marshal
	data, err := json.Marshal(contrib)
	if err != nil {
		t.Fatalf("Failed to marshal Contribution: %v", err)
	}

	// Unmarshal
	var decoded Contribution
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal Contribution: %v", err)
	}

	// Verify all fields
	if decoded.Repo != contrib.Repo {
		t.Errorf("Repo = %s, want %s", decoded.Repo, contrib.Repo)
	}

	if decoded.Owner != contrib.Owner {
		t.Errorf("Owner = %s, want %s", decoded.Owner, contrib.Owner)
	}

	if decoded.RepoName != contrib.RepoName {
		t.Errorf("RepoName = %s, want %s", decoded.RepoName, contrib.RepoName)
	}

	if decoded.Description != contrib.Description {
		t.Errorf("Description = %s, want %s", decoded.Description, contrib.Description)
	}

	if decoded.Stars != contrib.Stars {
		t.Errorf("Stars = %d, want %d", decoded.Stars, contrib.Stars)
	}

	if decoded.PRsMerged != contrib.PRsMerged {
		t.Errorf("PRsMerged = %d, want %d", decoded.PRsMerged, contrib.PRsMerged)
	}

	if decoded.Commits != contrib.Commits {
		t.Errorf("Commits = %d, want %d", decoded.Commits, contrib.Commits)
	}

	if decoded.Additions != contrib.Additions {
		t.Errorf("Additions = %d, want %d", decoded.Additions, contrib.Additions)
	}

	if decoded.Deletions != contrib.Deletions {
		t.Errorf("Deletions = %d, want %d", decoded.Deletions, contrib.Deletions)
	}

	if !decoded.FirstContribution.Equal(contrib.FirstContribution) {
		t.Errorf("FirstContribution = %v, want %v", decoded.FirstContribution, contrib.FirstContribution)
	}

	if !decoded.LastContribution.Equal(contrib.LastContribution) {
		t.Errorf("LastContribution = %v, want %v", decoded.LastContribution, contrib.LastContribution)
	}
}

func TestErrorTypeImplementsError(t *testing.T) {
	// Verify all error types implement error interface
	var _ error = &ErrRateLimited{}
	var _ error = &ErrAuthentication{}
	var _ error = &ErrNotFound{}
	var _ error = &ErrPartialResults{}
}

func TestStatsEmptyContributions(t *testing.T) {
	stats := Stats{
		Username:      "testuser",
		GeneratedAt:   time.Now().UTC(),
		Summary:       Summary{},
		Contributions: []Contribution{},
	}

	data, err := json.Marshal(stats)
	if err != nil {
		t.Fatalf("Failed to marshal empty Stats: %v", err)
	}

	var decoded Stats
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty Stats: %v", err)
	}

	if decoded.Contributions == nil {
		t.Error("Contributions should be empty slice, not nil")
	}

	if len(decoded.Contributions) != 0 {
		t.Errorf("Contributions length = %d, want 0", len(decoded.Contributions))
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
