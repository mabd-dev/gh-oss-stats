# gh-oss-stats

A Go library + CLI tool that fetches a GitHub user's open source contributions to external repositories (repos they don't own) and outputs structured JSON.

## Features

- 🔍 **Find External Contributions**: Discovers merged PRs to repositories you don't own
- 📊 **Aggregate Statistics**: Calculates total PRs, commits, lines of code changed
- ⭐ **Repository Filtering**: Filter by minimum star count
- 🚦 **Rate Limit Handling**: Smart rate limit detection with exponential backoff
- 📦 **Library-First Design**: Use as a Go library or standalone CLI

## Installation

### From Source

```bash
go install github.com/gh-oss-tools/gh-oss-stats/cmd/gh-oss-stats@latest
```

### Build Locally

```bash
git clone https://github.com/gh-oss-tools/gh-oss-stats.git
cd gh-oss-stats
go build -o gh-oss-stats ./cmd/gh-oss-stats
```

## Usage

### CLI

Basic usage:

```bash
# Fetch contributions for a user
gh-oss-stats --user github-username --token $GITHUB_TOKEN

# Filter by stars and limit PRs
gh-oss-stats -u github-username -t $GITHUB_TOKEN --min-stars 100 --max-prs 200

# Exclude your own organizations
gh-oss-stats -u github-username -t $GITHUB_TOKEN --exclude-orgs "my-org,my-company"

# Save to file with verbose logging
gh-oss-stats -u github-username -t $GITHUB_TOKEN -o output.json -v

# Show version
gh-oss-stats --version
```

### CLI Flags

```
--user, -u       GitHub username (required)
--token, -t      GitHub token (default: $GITHUB_TOKEN)
--include-loc    Include LOC metrics (default: true)
--include-prs    Include PR details (default: false)
--min-stars      Minimum repo stars (default: 0)
--max-prs        Max PRs to fetch (default: 500)
--exclude-orgs   Comma-separated list of organizations to exclude
--output, -o     Output file (default: stdout)
--verbose, -v    Verbose logging to stderr
--timeout        Timeout in seconds (default: 300)
--version        Print version
```

### As a Library

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/gh-oss-tools/gh-oss-stats/pkg/ossstats"
)

func main() {
    // Create client with options
    client := ossstats.New(
        ossstats.WithToken("your-github-token"),
        ossstats.WithMinStars(100),
        ossstats.WithExcludeOrgs([]string{"my-org", "my-company"}),
        ossstats.WithVerbose(),
    )
    
    // Fetch contributions
    stats, err := client.GetContributions(context.Background(), "github-username")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Total projects: %d\n", stats.Summary.TotalProjects)
    fmt.Printf("Total PRs merged: %d\n", stats.Summary.TotalPRsMerged)
}
```

## Output Format

```json
{
  "username": "github-username",
  "generatedAt": "2025-01-15T10:30:00Z",
  "summary": {
    "totalProjects": 42,
    "totalPRsMerged": 127,
    "totalCommits": 203,
    "totalAdditions": 5420,
    "totalDeletions": 2134
  },
  "contributions": [
    {
      "repo": "owner/repo-name",
      "owner": "owner",
      "repoName": "repo-name",
      "description": "An awesome project",
      "repoURL": "https://github.com/owner/repo-name",
      "stars": 1234,
      "prsMerged": 5,
      "commits": 12,
      "additions": 450,
      "deletions": 120,
      "firstContribution": "2024-01-10T08:20:00Z",
      "lastContribution": "2024-12-15T16:45:00Z"
    }
  ]
}
```

## GitHub Token

### Quick Setup

**Required for non-trivial usage** due to GitHub's rate limits:
- ❌ Without token: 60 requests/hour
- ✅ With token: 5,000 requests/hour

**1. Create a token:**
   - Go to https://github.com/settings/tokens
   - Generate new token (classic)
   - **No scopes needed** for public contributions (read-only access is sufficient)

**2. Set environment variable:**

Add to your `~/.bashrc` or `~/.zshrc`:
```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Reload: `source ~/.bashrc`

**3. Done!** The tool automatically uses `$GITHUB_TOKEN`:
```bash
gh-oss-stats --user YOUR_USERNAME
```

### Alternative: CLI Flag

For one-time use or CI/CD:
```bash
gh-oss-stats --user YOUR_USERNAME --token ghp_xxx...
```

### CI/CD (GitHub Actions)

GitHub Actions automatically provides `GITHUB_TOKEN`:
```yaml
- name: Fetch stats
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: gh-oss-stats --user ${{ github.actor }}
```

📖 **Full setup guide:** See [docs/TOKEN_SETUP.md](docs/TOKEN_SETUP.md)

## Rate Limiting

The tool implements smart rate limit handling:
- Respects GitHub's rate limits (5,000/hour core API, 30/min search API)
- Automatically waits when rate limited
- Returns partial results if rate limited mid-fetch
- Uses exponential backoff for retries

## Architecture

```
gh-oss-stats/
├── cmd/gh-oss-stats/     # CLI entry point
├── pkg/ossstats/         # Public API (importable)
│   ├── client.go         # Client + New()
│   ├── contributions.go  # GetContributions() logic
│   ├── types.go          # Exported types
│   └── options.go        # Functional options
└── internal/github/      # GitHub API client (private)
    ├── api.go            # HTTP client
    ├── ratelimit.go      # Rate limit handling
    └── types.go          # API response types
```

## Development

```bash
# Run tests
go test ./...

# Build
go build -o gh-oss-stats ./cmd/gh-oss-stats

# Lint (requires golangci-lint)
golangci-lint run
```

## License

See [LICENSE](LICENSE) file.

## Contributing

Contributions welcome! Please open an issue or PR.
