package github

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewAPIClient(t *testing.T) {
	httpClient := &http.Client{}
	token := "test-token"

	client := NewAPIClient(httpClient, token)

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.httpClient != httpClient {
		t.Error("httpClient not set correctly")
	}

	if client.token != token {
		t.Errorf("Expected token %s, got %s", token, client.token)
	}

	if client.baseURL != GitHubAPIBaseURL {
		t.Errorf("Expected baseURL %s, got %s", GitHubAPIBaseURL, client.baseURL)
	}
}

func TestNewAPIClientWithoutToken(t *testing.T) {
	httpClient := &http.Client{}

	client := NewAPIClient(httpClient, "")

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.token != "" {
		t.Errorf("Expected empty token, got %s", client.token)
	}
}

func TestAPIClientDoRequest(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		wantAuth   bool
		wantHeader map[string]string
	}{
		{
			name:     "with token",
			token:    "test-token-123",
			wantAuth: true,
			wantHeader: map[string]string{
				"Accept":              "application/vnd.github+json",
				"X-GitHub-Api-Version": APIVersion,
				"Authorization":        "Bearer test-token-123",
			},
		},
		{
			name:     "without token",
			token:    "",
			wantAuth: false,
			wantHeader: map[string]string{
				"Accept":              "application/vnd.github+json",
				"X-GitHub-Api-Version": APIVersion,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that validates headers
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for key, expectedValue := range tt.wantHeader {
					gotValue := r.Header.Get(key)
					if gotValue != expectedValue {
						t.Errorf("Header %s: expected %s, got %s", key, expectedValue, gotValue)
					}
				}

				if tt.wantAuth {
					auth := r.Header.Get("Authorization")
					if auth == "" {
						t.Error("Expected Authorization header, got empty")
					}
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"ok"}`))
			}))
			defer server.Close()

			client := NewAPIClient(&http.Client{}, tt.token)
			client.baseURL = server.URL

			ctx := context.Background()
			resp, err := client.doRequest(ctx, "GET", "/test", nil)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}
		})
	}
}

func TestAPIClientDoRequestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewAPIClient(&http.Client{}, "token")
	client.baseURL = server.URL

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.doRequest(ctx, "GET", "/test", nil)
	if err == nil {
		t.Error("Expected error due to context cancellation")
	}
}

func TestAPIClientGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "success",
		})
	}))
	defer server.Close()

	client := NewAPIClient(&http.Client{}, "token")
	client.baseURL = server.URL

	var result map[string]string
	ctx := context.Background()
	resp, err := client.get(ctx, "/test", &result)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if result["message"] != "success" {
		t.Errorf("Expected message 'success', got %s", result["message"])
	}
}

func TestAPIClientGetHTTPError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
	}{
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			body:       `{"message":"Not Found"}`,
		},
		{
			name:       "401 Unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       `{"message":"Bad credentials"}`,
		},
		{
			name:       "403 Forbidden",
			statusCode: http.StatusForbidden,
			body:       `{"message":"Rate limit exceeded"}`,
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			body:       `{"message":"Internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			}))
			defer server.Close()

			client := NewAPIClient(&http.Client{}, "token")
			client.baseURL = server.URL

			var result map[string]string
			ctx := context.Background()
			_, err := client.get(ctx, "/test", &result)

			if err == nil {
				t.Error("Expected error for HTTP error response")
			}

			if !strings.Contains(err.Error(), "HTTP") {
				t.Errorf("Expected HTTP error message, got: %v", err)
			}
		})
	}
}

func TestAPIClientGetInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client := NewAPIClient(&http.Client{}, "token")
	client.baseURL = server.URL

	var result map[string]string
	ctx := context.Background()
	_, err := client.get(ctx, "/test", &result)

	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "decoding") {
		t.Errorf("Expected decoding error, got: %v", err)
	}
}

func TestAPIClientSearchIssues(t *testing.T) {
	expectedQuery := "author:testuser type:pr is:merged"
	expectedPage := 1
	expectedPerPage := 50

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/search/issues") {
			t.Errorf("Expected path /search/issues, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("q") != expectedQuery {
			t.Errorf("Expected query %s, got %s", expectedQuery, query.Get("q"))
		}

		if query.Get("page") != "1" {
			t.Errorf("Expected page 1, got %s", query.Get("page"))
		}

		if query.Get("per_page") != "50" {
			t.Errorf("Expected per_page 50, got %s", query.Get("per_page"))
		}

		if query.Get("sort") != "updated" {
			t.Errorf("Expected sort 'updated', got %s", query.Get("sort"))
		}

		if query.Get("order") != "desc" {
			t.Errorf("Expected order 'desc', got %s", query.Get("order"))
		}

		response := SearchIssuesResponse{
			TotalCount:        1,
			IncompleteResults: false,
			Items: []Issue{
				{
					Number: 123,
					Title:  "Test PR",
					State:  "closed",
					User: User{
						Login: "testuser",
						ID:    12345,
						Type:  "User",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAPIClient(&http.Client{}, "token")
	client.baseURL = server.URL

	ctx := context.Background()
	result, resp, err := client.SearchIssues(ctx, expectedQuery, expectedPage, expectedPerPage)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if result.TotalCount != 1 {
		t.Errorf("Expected TotalCount 1, got %d", result.TotalCount)
	}

	if len(result.Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(result.Items))
	}

	if result.Items[0].Number != 123 {
		t.Errorf("Expected issue number 123, got %d", result.Items[0].Number)
	}
}

func TestAPIClientGetPullRequest(t *testing.T) {
	expectedOwner := "testowner"
	expectedRepo := "testrepo"
	expectedNumber := 456

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/repos/testowner/testrepo/pulls/456"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		mergedAt := time.Now().UTC()
		response := PullRequest{
			Number:    456,
			State:     "closed",
			Title:     "Test Pull Request",
			Merged:    true,
			MergedAt:  &mergedAt,
			Commits:   5,
			Additions: 100,
			Deletions: 20,
			User: User{
				Login: "contributor",
				ID:    67890,
				Type:  "User",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAPIClient(&http.Client{}, "token")
	client.baseURL = server.URL

	ctx := context.Background()
	result, resp, err := client.GetPullRequest(ctx, expectedOwner, expectedRepo, expectedNumber)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if result.Number != 456 {
		t.Errorf("Expected number 456, got %d", result.Number)
	}

	if !result.Merged {
		t.Error("Expected Merged to be true")
	}

	if result.Commits != 5 {
		t.Errorf("Expected 5 commits, got %d", result.Commits)
	}

	if result.Additions != 100 {
		t.Errorf("Expected 100 additions, got %d", result.Additions)
	}

	if result.Deletions != 20 {
		t.Errorf("Expected 20 deletions, got %d", result.Deletions)
	}
}

func TestAPIClientGetRepository(t *testing.T) {
	expectedOwner := "testowner"
	expectedRepo := "testrepo"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/repos/testowner/testrepo"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		response := Repository{
			Name:            "testrepo",
			FullName:        "testowner/testrepo",
			Description:     "A test repository",
			StargazersCount: 1000,
			Language:        "Go",
			Fork:            false,
			Owner: User{
				Login: "testowner",
				ID:    11111,
				Type:  "Organization",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAPIClient(&http.Client{}, "token")
	client.baseURL = server.URL

	ctx := context.Background()
	result, resp, err := client.GetRepository(ctx, expectedOwner, expectedRepo)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if result.Name != "testrepo" {
		t.Errorf("Expected name 'testrepo', got %s", result.Name)
	}

	if result.FullName != "testowner/testrepo" {
		t.Errorf("Expected full_name 'testowner/testrepo', got %s", result.FullName)
	}

	if result.StargazersCount != 1000 {
		t.Errorf("Expected 1000 stars, got %d", result.StargazersCount)
	}

	if result.Language != "Go" {
		t.Errorf("Expected language 'Go', got %s", result.Language)
	}
}

func TestAPIClientGetRateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rate_limit" {
			t.Errorf("Expected path /rate_limit, got %s", r.URL.Path)
		}

		response := RateLimitResponse{
			Resources: RateLimitResources{
				Core: RateLimit{
					Limit:     5000,
					Remaining: 4999,
					Reset:     1672531200,
					Used:      1,
				},
				Search: RateLimit{
					Limit:     30,
					Remaining: 29,
					Reset:     1672531200,
					Used:      1,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAPIClient(&http.Client{}, "token")
	client.baseURL = server.URL

	ctx := context.Background()
	result, err := client.GetRateLimit(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Resources.Core.Limit != 5000 {
		t.Errorf("Expected core limit 5000, got %d", result.Resources.Core.Limit)
	}

	if result.Resources.Core.Remaining != 4999 {
		t.Errorf("Expected core remaining 4999, got %d", result.Resources.Core.Remaining)
	}

	if result.Resources.Search.Limit != 30 {
		t.Errorf("Expected search limit 30, got %d", result.Resources.Search.Limit)
	}

	if result.Resources.Search.Remaining != 29 {
		t.Errorf("Expected search remaining 29, got %d", result.Resources.Search.Remaining)
	}
}

func TestParseRepoURL(t *testing.T) {
	tests := []struct {
		name      string
		repoURL   string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{
			name:      "valid API URL",
			repoURL:   "https://api.github.com/repos/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "valid API URL with trailing slash",
			repoURL:   "https://api.github.com/repos/owner/repo/",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "valid web URL",
			repoURL:   "https://github.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "simple owner/repo",
			repoURL:   "owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:    "invalid - single part",
			repoURL: "repo",
			wantErr: true,
		},
		{
			name:    "invalid - empty string",
			repoURL: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := ParseRepoURL(tt.repoURL)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if owner != tt.wantOwner {
				t.Errorf("Expected owner %s, got %s", tt.wantOwner, owner)
			}

			if repo != tt.wantRepo {
				t.Errorf("Expected repo %s, got %s", tt.wantRepo, repo)
			}
		})
	}
}

func TestParseLinkHeader(t *testing.T) {
	tests := []struct {
		name       string
		linkHeader string
		want       map[string]string
	}{
		{
			name:       "single next link",
			linkHeader: `<https://api.github.com/search/issues?page=2>; rel="next"`,
			want: map[string]string{
				"next": "https://api.github.com/search/issues?page=2",
			},
		},
		{
			name:       "next and last links",
			linkHeader: `<https://api.github.com/search/issues?page=2>; rel="next", <https://api.github.com/search/issues?page=10>; rel="last"`,
			want: map[string]string{
				"next": "https://api.github.com/search/issues?page=2",
				"last": "https://api.github.com/search/issues?page=10",
			},
		},
		{
			name:       "prev, next, first, last",
			linkHeader: `<https://api.github.com/search/issues?page=3>; rel="prev", <https://api.github.com/search/issues?page=5>; rel="next", <https://api.github.com/search/issues?page=1>; rel="first", <https://api.github.com/search/issues?page=10>; rel="last"`,
			want: map[string]string{
				"prev":  "https://api.github.com/search/issues?page=3",
				"next":  "https://api.github.com/search/issues?page=5",
				"first": "https://api.github.com/search/issues?page=1",
				"last":  "https://api.github.com/search/issues?page=10",
			},
		},
		{
			name:       "empty header",
			linkHeader: "",
			want:       map[string]string{},
		},
		{
			name:       "malformed - missing semicolon",
			linkHeader: `<https://api.github.com/search/issues?page=2> rel="next"`,
			want:       map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLinkHeader(tt.linkHeader)

			if len(got) != len(tt.want) {
				t.Errorf("Expected %d links, got %d", len(tt.want), len(got))
			}

			for key, expectedValue := range tt.want {
				gotValue, ok := got[key]
				if !ok {
					t.Errorf("Expected key %s not found in result", key)
					continue
				}

				if gotValue != expectedValue {
					t.Errorf("For key %s: expected %s, got %s", key, expectedValue, gotValue)
				}
			}
		})
	}
}

func TestAPIClientGetNilResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test":"value"}`))
	}))
	defer server.Close()

	client := NewAPIClient(&http.Client{}, "token")
	client.baseURL = server.URL

	ctx := context.Background()
	resp, err := client.get(ctx, "/test", nil)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestAPIClientDoRequestWithBody(t *testing.T) {
	requestBody := "test body content"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}

		if string(body) != requestBody {
			t.Errorf("Expected body %s, got %s", requestBody, string(body))
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewAPIClient(&http.Client{}, "token")
	client.baseURL = server.URL

	ctx := context.Background()
	resp, err := client.doRequest(ctx, "POST", "/test", strings.NewReader(requestBody))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestAPIClientSearchIssuesError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message":"API rate limit exceeded"}`))
	}))
	defer server.Close()

	client := NewAPIClient(&http.Client{}, "token")
	client.baseURL = server.URL

	ctx := context.Background()
	result, resp, err := client.SearchIssues(ctx, "test query", 1, 50)

	if err == nil {
		t.Error("Expected error for rate limit response")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}

	if resp == nil {
		t.Error("Expected non-nil response even on error")
	}
}

func TestAPIClientGetPullRequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer server.Close()

	client := NewAPIClient(&http.Client{}, "token")
	client.baseURL = server.URL

	ctx := context.Background()
	result, resp, err := client.GetPullRequest(ctx, "owner", "repo", 999)

	if err == nil {
		t.Error("Expected error for not found response")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}

	if resp == nil {
		t.Error("Expected non-nil response even on error")
	}
}

func TestAPIClientGetRepositoryError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Bad credentials"}`))
	}))
	defer server.Close()

	client := NewAPIClient(&http.Client{}, "token")
	client.baseURL = server.URL

	ctx := context.Background()
	result, resp, err := client.GetRepository(ctx, "owner", "repo")

	if err == nil {
		t.Error("Expected error for unauthorized response")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}

	if resp == nil {
		t.Error("Expected non-nil response even on error")
	}
}
