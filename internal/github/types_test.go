package github

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSearchIssuesResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"total_count": 42,
		"incomplete_results": false,
		"items": [
			{
				"number": 123,
				"title": "Test PR",
				"state": "closed",
				"created_at": "2023-01-01T00:00:00Z",
				"updated_at": "2023-01-02T00:00:00Z",
				"closed_at": "2023-01-03T00:00:00Z",
				"repository_url": "https://api.github.com/repos/owner/repo",
				"html_url": "https://github.com/owner/repo/pull/123",
				"user": {
					"login": "testuser",
					"id": 12345,
					"type": "User"
				},
				"pull_request": {
					"url": "https://api.github.com/repos/owner/repo/pulls/123",
					"html_url": "https://github.com/owner/repo/pull/123",
					"diff_url": "https://github.com/owner/repo/pull/123.diff",
					"patch_url": "https://github.com/owner/repo/pull/123.patch",
					"merged_at": "2023-01-03T00:00:00Z"
				}
			}
		]
	}`

	var resp SearchIssuesResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal SearchIssuesResponse: %v", err)
	}

	if resp.TotalCount != 42 {
		t.Errorf("Expected TotalCount 42, got %d", resp.TotalCount)
	}

	if resp.IncompleteResults {
		t.Error("Expected IncompleteResults to be false")
	}

	if len(resp.Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(resp.Items))
	}

	item := resp.Items[0]
	if item.Number != 123 {
		t.Errorf("Expected number 123, got %d", item.Number)
	}

	if item.Title != "Test PR" {
		t.Errorf("Expected title 'Test PR', got %s", item.Title)
	}

	if item.User.Login != "testuser" {
		t.Errorf("Expected user login 'testuser', got %s", item.User.Login)
	}

	if item.PullRequest == nil {
		t.Fatal("Expected PullRequest to be non-nil")
	}

	if item.PullRequest.MergedAt == nil {
		t.Fatal("Expected MergedAt to be non-nil")
	}
}

func TestIssueWithNullClosedAt(t *testing.T) {
	jsonData := `{
		"number": 456,
		"title": "Open Issue",
		"state": "open",
		"created_at": "2023-01-01T00:00:00Z",
		"updated_at": "2023-01-02T00:00:00Z",
		"closed_at": null,
		"repository_url": "https://api.github.com/repos/owner/repo",
		"html_url": "https://github.com/owner/repo/issues/456",
		"user": {
			"login": "testuser",
			"id": 12345,
			"type": "User"
		}
	}`

	var issue Issue
	err := json.Unmarshal([]byte(jsonData), &issue)
	if err != nil {
		t.Fatalf("Failed to unmarshal Issue: %v", err)
	}

	if issue.ClosedAt != nil {
		t.Error("Expected ClosedAt to be nil for open issue")
	}

	if issue.PullRequest != nil {
		t.Error("Expected PullRequest to be nil for regular issue")
	}
}

func TestPullRequestUnmarshal(t *testing.T) {
	jsonData := `{
		"number": 789,
		"state": "closed",
		"title": "Add feature X",
		"user": {
			"login": "contributor",
			"id": 67890,
			"type": "User"
		},
		"created_at": "2023-01-01T00:00:00Z",
		"updated_at": "2023-01-02T00:00:00Z",
		"closed_at": "2023-01-03T00:00:00Z",
		"merged_at": "2023-01-03T00:00:00Z",
		"merged": true,
		"commits": 5,
		"additions": 120,
		"deletions": 30,
		"changed_files": 3,
		"html_url": "https://github.com/owner/repo/pull/789"
	}`

	var pr PullRequest
	err := json.Unmarshal([]byte(jsonData), &pr)
	if err != nil {
		t.Fatalf("Failed to unmarshal PullRequest: %v", err)
	}

	if pr.Number != 789 {
		t.Errorf("Expected number 789, got %d", pr.Number)
	}

	if !pr.Merged {
		t.Error("Expected Merged to be true")
	}

	if pr.Commits != 5 {
		t.Errorf("Expected 5 commits, got %d", pr.Commits)
	}

	if pr.Additions != 120 {
		t.Errorf("Expected 120 additions, got %d", pr.Additions)
	}

	if pr.Deletions != 30 {
		t.Errorf("Expected 30 deletions, got %d", pr.Deletions)
	}

	if pr.ChangedFiles != 3 {
		t.Errorf("Expected 3 changed files, got %d", pr.ChangedFiles)
	}

	if pr.MergedAt == nil {
		t.Error("Expected MergedAt to be non-nil")
	}
}

func TestRepositoryUnmarshal(t *testing.T) {
	jsonData := `{
		"name": "awesome-project",
		"full_name": "owner/awesome-project",
		"owner": {
			"login": "owner",
			"id": 11111,
			"type": "Organization"
		},
		"description": "An awesome open source project",
		"html_url": "https://github.com/owner/awesome-project",
		"fork": false,
		"created_at": "2020-01-01T00:00:00Z",
		"updated_at": "2023-01-01T00:00:00Z",
		"pushed_at": "2023-01-02T00:00:00Z",
		"stargazers_count": 1500,
		"language": "Go",
		"forks_count": 200,
		"open_issues_count": 25,
		"default_branch": "main"
	}`

	var repo Repository
	err := json.Unmarshal([]byte(jsonData), &repo)
	if err != nil {
		t.Fatalf("Failed to unmarshal Repository: %v", err)
	}

	if repo.Name != "awesome-project" {
		t.Errorf("Expected name 'awesome-project', got %s", repo.Name)
	}

	if repo.FullName != "owner/awesome-project" {
		t.Errorf("Expected full_name 'owner/awesome-project', got %s", repo.FullName)
	}

	if repo.Owner.Login != "owner" {
		t.Errorf("Expected owner login 'owner', got %s", repo.Owner.Login)
	}

	if repo.Owner.Type != "Organization" {
		t.Errorf("Expected owner type 'Organization', got %s", repo.Owner.Type)
	}

	if repo.Description != "An awesome open source project" {
		t.Errorf("Expected specific description, got %s", repo.Description)
	}

	if repo.Fork {
		t.Error("Expected Fork to be false")
	}

	if repo.StargazersCount != 1500 {
		t.Errorf("Expected 1500 stars, got %d", repo.StargazersCount)
	}

	if repo.Language != "Go" {
		t.Errorf("Expected language 'Go', got %s", repo.Language)
	}

	if repo.DefaultBranch != "main" {
		t.Errorf("Expected default branch 'main', got %s", repo.DefaultBranch)
	}

	if repo.PushedAt == nil {
		t.Error("Expected PushedAt to be non-nil")
	}
}

func TestRepositoryWithNullPushedAt(t *testing.T) {
	jsonData := `{
		"name": "empty-repo",
		"full_name": "owner/empty-repo",
		"owner": {
			"login": "owner",
			"id": 11111,
			"type": "User"
		},
		"description": "",
		"html_url": "https://github.com/owner/empty-repo",
		"fork": false,
		"created_at": "2023-01-01T00:00:00Z",
		"updated_at": "2023-01-01T00:00:00Z",
		"pushed_at": null,
		"stargazers_count": 0,
		"language": "",
		"forks_count": 0,
		"open_issues_count": 0,
		"default_branch": "main"
	}`

	var repo Repository
	err := json.Unmarshal([]byte(jsonData), &repo)
	if err != nil {
		t.Fatalf("Failed to unmarshal Repository: %v", err)
	}

	if repo.PushedAt != nil {
		t.Error("Expected PushedAt to be nil for empty repo")
	}

	if repo.StargazersCount != 0 {
		t.Errorf("Expected 0 stars, got %d", repo.StargazersCount)
	}

	if repo.Language != "" {
		t.Errorf("Expected empty language, got %s", repo.Language)
	}
}

func TestRateLimitResponseUnmarshal(t *testing.T) {
	jsonData := `{
		"resources": {
			"core": {
				"limit": 5000,
				"remaining": 4999,
				"reset": 1672531200,
				"used": 1
			},
			"search": {
				"limit": 30,
				"remaining": 29,
				"reset": 1672531200,
				"used": 1
			}
		}
	}`

	var rateLimitResp RateLimitResponse
	err := json.Unmarshal([]byte(jsonData), &rateLimitResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal RateLimitResponse: %v", err)
	}

	if rateLimitResp.Resources.Core.Limit != 5000 {
		t.Errorf("Expected core limit 5000, got %d", rateLimitResp.Resources.Core.Limit)
	}

	if rateLimitResp.Resources.Core.Remaining != 4999 {
		t.Errorf("Expected core remaining 4999, got %d", rateLimitResp.Resources.Core.Remaining)
	}

	if rateLimitResp.Resources.Core.Reset != 1672531200 {
		t.Errorf("Expected core reset 1672531200, got %d", rateLimitResp.Resources.Core.Reset)
	}

	if rateLimitResp.Resources.Core.Used != 1 {
		t.Errorf("Expected core used 1, got %d", rateLimitResp.Resources.Core.Used)
	}

	if rateLimitResp.Resources.Search.Limit != 30 {
		t.Errorf("Expected search limit 30, got %d", rateLimitResp.Resources.Search.Limit)
	}

	if rateLimitResp.Resources.Search.Remaining != 29 {
		t.Errorf("Expected search remaining 29, got %d", rateLimitResp.Resources.Search.Remaining)
	}
}

func TestUserUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     User
	}{
		{
			name: "regular user",
			jsonData: `{
				"login": "johndoe",
				"id": 12345,
				"type": "User"
			}`,
			want: User{
				Login: "johndoe",
				ID:    12345,
				Type:  "User",
			},
		},
		{
			name: "organization",
			jsonData: `{
				"login": "myorg",
				"id": 67890,
				"type": "Organization"
			}`,
			want: User{
				Login: "myorg",
				ID:    67890,
				Type:  "Organization",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var user User
			err := json.Unmarshal([]byte(tt.jsonData), &user)
			if err != nil {
				t.Fatalf("Failed to unmarshal User: %v", err)
			}

			if user.Login != tt.want.Login {
				t.Errorf("Expected login %s, got %s", tt.want.Login, user.Login)
			}

			if user.ID != tt.want.ID {
				t.Errorf("Expected ID %d, got %d", tt.want.ID, user.ID)
			}

			if user.Type != tt.want.Type {
				t.Errorf("Expected type %s, got %s", tt.want.Type, user.Type)
			}
		})
	}
}

func TestPullRequestRefUnmarshal(t *testing.T) {
	jsonData := `{
		"url": "https://api.github.com/repos/owner/repo/pulls/123",
		"html_url": "https://github.com/owner/repo/pull/123",
		"diff_url": "https://github.com/owner/repo/pull/123.diff",
		"patch_url": "https://github.com/owner/repo/pull/123.patch",
		"merged_at": "2023-01-03T12:00:00Z"
	}`

	var prRef PullRequestRef
	err := json.Unmarshal([]byte(jsonData), &prRef)
	if err != nil {
		t.Fatalf("Failed to unmarshal PullRequestRef: %v", err)
	}

	if prRef.URL != "https://api.github.com/repos/owner/repo/pulls/123" {
		t.Errorf("Expected specific URL, got %s", prRef.URL)
	}

	if prRef.HTMLURL != "https://github.com/owner/repo/pull/123" {
		t.Errorf("Expected specific HTMLURL, got %s", prRef.HTMLURL)
	}

	if prRef.DiffURL != "https://github.com/owner/repo/pull/123.diff" {
		t.Errorf("Expected specific DiffURL, got %s", prRef.DiffURL)
	}

	if prRef.PatchURL != "https://github.com/owner/repo/pull/123.patch" {
		t.Errorf("Expected specific PatchURL, got %s", prRef.PatchURL)
	}

	if prRef.MergedAt == nil {
		t.Fatal("Expected MergedAt to be non-nil")
	}

	expectedTime := time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC)
	if !prRef.MergedAt.Equal(expectedTime) {
		t.Errorf("Expected MergedAt %v, got %v", expectedTime, prRef.MergedAt)
	}
}

func TestPullRequestRefWithNullMergedAt(t *testing.T) {
	jsonData := `{
		"url": "https://api.github.com/repos/owner/repo/pulls/123",
		"html_url": "https://github.com/owner/repo/pull/123",
		"diff_url": "https://github.com/owner/repo/pull/123.diff",
		"patch_url": "https://github.com/owner/repo/pull/123.patch",
		"merged_at": null
	}`

	var prRef PullRequestRef
	err := json.Unmarshal([]byte(jsonData), &prRef)
	if err != nil {
		t.Fatalf("Failed to unmarshal PullRequestRef: %v", err)
	}

	if prRef.MergedAt != nil {
		t.Error("Expected MergedAt to be nil for unmerged PR")
	}
}

func TestStructMarshalUnmarshalRoundTrip(t *testing.T) {
	original := SearchIssuesResponse{
		TotalCount:        10,
		IncompleteResults: false,
		Items: []Issue{
			{
				Number:        1,
				Title:         "Test Issue",
				State:         "open",
				CreatedAt:     time.Now().UTC().Truncate(time.Second),
				UpdatedAt:     time.Now().UTC().Truncate(time.Second),
				ClosedAt:      nil,
				RepositoryURL: "https://api.github.com/repos/test/repo",
				HTMLURL:       "https://github.com/test/repo/issues/1",
				User: User{
					Login: "testuser",
					ID:    123,
					Type:  "User",
				},
				PullRequest: nil,
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal back
	var decoded SearchIssuesResponse
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Compare
	if decoded.TotalCount != original.TotalCount {
		t.Errorf("TotalCount mismatch: got %d, want %d", decoded.TotalCount, original.TotalCount)
	}

	if decoded.IncompleteResults != original.IncompleteResults {
		t.Errorf("IncompleteResults mismatch: got %v, want %v", decoded.IncompleteResults, original.IncompleteResults)
	}

	if len(decoded.Items) != len(original.Items) {
		t.Fatalf("Items length mismatch: got %d, want %d", len(decoded.Items), len(original.Items))
	}

	if decoded.Items[0].Title != original.Items[0].Title {
		t.Errorf("Title mismatch: got %s, want %s", decoded.Items[0].Title, original.Items[0].Title)
	}
}
