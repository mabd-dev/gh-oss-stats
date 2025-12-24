# CLAUDE.md

> This file provides context for Claude Code CLI when working on this project.

## Project Overview

**gh-oss-stats** is a Go library + CLI tool that fetches a GitHub user's open source contributions to external repositories (repos they don't own). It outputs structured JSON for consumption by other tools (websites, badge services, etc.).

**Organization:** `mabd-dev`  
**Repository:** `github.com/mabd-dev/gh-oss-stats`  
**Go Version:** 1.21+

## Architecture

```
gh-oss-stats/
├── cmd/gh-oss-stats/main.go   # CLI entry point (thin wrapper)
├── pkg/ossstats/              # PUBLIC API - importable by external projects
│   ├── client.go              # Client struct + New() constructor
│   ├── contributions.go       # Core logic: GetContributions()
│   ├── types.go               # All exported types (Stats, Contribution, etc.)
│   └── options.go             # Functional options (WithToken, WithLOC, etc.)
├── internal/github/           # PRIVATE - GitHub API implementation
│   ├── api.go                 # HTTP client, request helpers
│   ├── ratelimit.go           # Rate limit handling, backoff
│   └── types.go               # API response structs (internal only)
├── go.mod
└── README.md
```

**Key Principle:** Library-first design. All logic lives in `pkg/ossstats/`. The CLI in `cmd/` is just a thin wrapper that parses flags and calls the library.

## Core Types

```go
// pkg/ossstats/types.go

type Stats struct {
    Username      string         `json:"username"`
    GeneratedAt   time.Time      `json:"generated_at"`
    Summary       Summary        `json:"summary"`
    Contributions []Contribution `json:"contributions"`
}

type Summary struct {
    TotalProjects  int `json:"total_projects"`
    TotalPRsMerged int `json:"total_prs_merged"`
    TotalCommits   int `json:"total_commits"`
    TotalAdditions int `json:"total_additions"`
    TotalDeletions int `json:"total_deletions"`
}

type Contribution struct {
    Repo              string    `json:"repo"`
    Owner             string    `json:"owner"`
    RepoName          string    `json:"repo_name"`
    Description       string    `json:"description"`
    RepoURL           string    `json:"repo_url"`
    Stars             int       `json:"stars"`
    PRsMerged         int       `json:"prs_merged"`
    Commits           int       `json:"commits"`
    Additions         int       `json:"additions"`
    Deletions         int       `json:"deletions"`
    FirstContribution time.Time `json:"first_contribution"`
    LastContribution  time.Time `json:"last_contribution"`
}
```

## GitHub API Strategy

**Step 1: Find merged PRs to external repos**
```
GET /search/issues?q=author:{username}+type:pr+is:merged+-user:{username}&per_page=100
```
*Note: Organizations can be excluded by appending `-org:{orgname}` to the query for each excluded org.*

**Step 2: Get PR details (commits, additions, deletions)**
```
GET /repos/{owner}/{repo}/pulls/{pull_number}
```

**Step 3: Get repo metadata (stars, description)**
```
GET /repos/{owner}/{repo}
```

## Functional Options Pattern

Always use this pattern for client configuration:

```go
client := ossstats.New(
    ossstats.WithToken(token),
    ossstats.WithLOC(true),
    ossstats.WithMinStars(100),
    ossstats.WithExcludeOrgs([]string{"my-org", "my-company"}),
)
```

Required options to implement:
- `WithToken(string)` - GitHub PAT (required for reasonable rate limits)
- `WithLOC(bool)` - Include lines of code metrics (default: true)
- `WithPRDetails(bool)` - Include detailed PR list (default: false)
- `WithMinStars(int)` - Filter repos by minimum stars (default: 0)
- `WithMaxPRs(int)` - Limit PRs fetched (default: 500)
- `WithExcludeOrgs([]string)` - Exclude organizations from the report (default: none)
- `WithTimeout(time.Duration)` - Overall timeout (default: 5m)
- `WithLogger(Logger)` - Custom logger interface
- `WithHTTPClient(*http.Client)` - Custom HTTP client

## CLI Flags

```
--user, -u       string    GitHub username (required)
--token, -t      string    GitHub token (default: $GITHUB_TOKEN)
--include-loc    bool      Include LOC metrics (default: true)
--include-prs    bool      Include PR details (default: false)
--min-stars      int       Minimum repo stars (default: 0)
--max-prs        int       Max PRs to fetch (default: 500)
--exclude-orgs   string    Comma-separated list of organizations to exclude
--output, -o     string    Output file (default: stdout)
--verbose, -v    bool      Verbose logging to stderr
--timeout        duration  Timeout (default: 5m)
--version        bool      Print version
```

## Rate Limit Handling

GitHub limits:
- Core API: 5,000/hour (authenticated)
- Search API: 30/minute (authenticated)

Implementation requirements:
1. Check `X-RateLimit-Remaining` and `X-RateLimit-Reset` headers
2. Exponential backoff on 429 responses
3. 2-second delay between search API calls
4. Return partial results with `ErrPartialResults` if rate limited mid-fetch

## Error Types

```go
type ErrRateLimited struct { ResetAt time.Time; Message string }
type ErrAuthentication struct { Message string }
type ErrNotFound struct { Username string }
type ErrPartialResults struct { Stats *Stats; Errors []error; Message string }
```

## Commands

```bash
# Run tests
go test ./...

# Build CLI
go build -o gh-oss-stats ./cmd/gh-oss-stats

# Install locally
go install ./cmd/gh-oss-stats

# Run CLI
./gh-oss-stats --user mabd-dev --token $GITHUB_TOKEN

# Lint (if golangci-lint installed)
golangci-lint run
```

## Implementation Notes

1. **All API calls must accept `context.Context`** for cancellation/timeout
2. **No external dependencies** in `pkg/ossstats/` - stdlib only
3. **Pagination:** GitHub returns max 100 items per page; handle `Link` header
4. **Concurrency:** Consider parallel PR fetching with `errgroup` (limit to 5 concurrent)
5. **Timestamps:** Always UTC, use `time.RFC3339` for JSON
6. **Logging:** Use the `Logger` interface, never `fmt.Print` in library code

## Testing

- Unit tests: Mock HTTP responses using `httptest.Server`
- Test data: Store fixtures in `pkg/ossstats/testdata/`
- Integration tests: Tag with `//go:build integration` and skip without token

## Do NOT

- Add external dependencies without strong justification
- Put business logic in `cmd/` - it belongs in `pkg/ossstats/`
- Make unauthenticated requests by default (rate limits too restrictive)
- Ignore context cancellation
- Return zero values on error - always return descriptive errors
- Use `log.Fatal` or `os.Exit` in library code
