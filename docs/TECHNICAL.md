# Technical Documentation

## Installation

### From Source

```bash
go install github.com/mabd-dev/mabd-dev/cmd/gh-oss-stats@latest
```

### Build Locally

```bash
git clone https://github.com/mabd-dev/gh-oss-stats.git
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

**Data Fetching:**

| Flag | Type | Default | Description |
|-------|-----------|-------------|-------------|
| --user, -u | string | "" | Github username |
| --token, -t | string | $GITHUB_TOKEN | Github token |
| --include-loc | bool | false | Include LOC metrics (line of code) |
| --include-prs | bool | false | Include PR details |
| --min-stars | int | 0 | Minimum repo stars |
| --max-prs | int | 500 | Max PRs to fetch |
| --exclude-orgs | string | "" | Comma-separated list of organizations to exclude |
| --output, -o | string | "" | Output file path |
| --verbose, -v | bool | false | Verbose logging |
| -- timeout | int | 300 | Timeout in **seconds** |
| --version | bool | false | Print version |


**Development:**

| Flag | Type | Default | Description |
|-------|-----------|-------------|-------------|
| --debug | boolean | false | Uses fake data when true |


**Badge Generation:**

| Flag | Type | Default | Description |
|-------|-----------|-------------|-------------|
| --badge | boolean | false | Generate SVG Badge |
| --badge-style | string | summary | Available options: summary, compact, detailed |
| --badge-variant | string | default | Available options: default, text-based |
| --badge-theme | string | dark | Available options: dark, light, nord, dracula, gruvbox-light, gruvbox-dar |
| --badge-output | string | ./badge.svg | Badge output file path |
| --badge-sort | string | prs | Sort contributions by: prs, stars, commits |
| --badge-limit | int | 5 | Number of contributions to show in detailed badge |




### Badge Generation

Generate beautiful SVG badges from your contribution stats:

```bash
# Generate a summary badge (400x200)
gh-oss-stats --user mabd-dev --badge

# Generate a compact shields.io style badge (280x28)
gh-oss-stats --user mabd-dev --badge --badge-style compact --badge-theme light

# Generate a detailed badge with top 10 repos sorted by stars (400x320)
gh-oss-stats --user mabd-dev --badge --badge-style detailed --badge-sort stars --badge-limit 10
```

**Badge Styles:**

| Style | Dimensions | Description |
|-------|-----------|-------------|
| `summary` | 400×200 | Key metrics: projects, PRs, commits, lines |
| `compact` | 280×28 | Shields.io style: "42 projects \| 1.6K PRs" |
| `detailed` | 400×320 | Summary + top N contributions with stars & PRs |

Check [All Combos](/badges/BADGE_THEMES.md)

**Example:**
```bash
# Fetch stats + generate both JSON and badge
gh-oss-stats --user mabd-dev \
  -o stats.json \
  --badge \
  --badge-style detailed \
  --badge-theme dark \
  --badge-output badge.svg
```

### Local Development & Testing

Use debug mode to test the tool locally without hitting GitHub API:

```bash
# Fast local testing with mock data (no token required)
gh-oss-stats --user test-user --debug

# Generate badge from mock data
gh-oss-stats --user test-user --debug --badge --badge-output test-badge.svg

# Test different badge styles with mock data
gh-oss-stats --user test-user --debug \
  --badge \
  --badge-style detailed \
  --badge-theme light \
  --badge-output test-detailed.svg
```

**Debug Mode Benefits:**
- ✅ No API rate limits
- ✅ Instant results (<1 second)
- ✅ No GitHub token required
- ✅ Uses static mock data from `internal/github/mockResponses/`
- ✅ Perfect for development and CI testing

### As a Library

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/mabd-dev/gh-oss-stats/pkg/ossstats"
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


## Prerequisites

### For GitHub Actions (Recommended)
✅ **No setup required!** GitHub Actions automatically provides a token with `${{ secrets.GITHUB_TOKEN }}`.

### For Local CLI Usage
You have two options:

**Option 1: GitHub Token (for real data)**
- Required for fetching your actual contribution data
- See [TOKEN_SETUP.md](TOKEN_SETUP.md) for setup instructions

**Option 2: Debug Mode (for testing)**
- Use `--debug` flag for testing without a token
- Uses mock data, perfect for development

## GitHub Token

**Required for non-trivial usage** due to GitHub's rate limits:
- ❌ Without token: 60 requests/hour
- ✅ With token: 5,000 requests/hour

📖 **Full setup guide:** See [TOKEN_SETUP.md](TOKEN_SETUP.md)

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


## Rate Limiting

The tool implements smart rate limit handling:
- Respects GitHub's rate limits (5,000/hour core API, 30/min search API)
- Automatically waits when rate limited
- Returns partial results if rate limited mid-fetch
- Uses exponential backoff for retries

## Architecture

```
gh-oss-stats/
├── cmd/gh-oss-stats/           # CLI entry point
├── pkg/ossstats/               # Public API (importable)
│   ├── badge/                  # Badge generation folder
    │   ├── badgeTemplates/     # Defines all badge svg templates
    │   ├── badge.go            # Generate and save badge
    │   ├── badgeSortBy.go      # Defines sorting types
    │   ├── badgeStyle.go       # Defines all badge styles + helper function
    │   ├── badgeTheme.go       # Defines all badge themes + helper function
    │   ├── badgeVariant.go     # Defines all badge variants + helper function
    │   └── types.go            # Client + New()
│   ├── client.go               # Client + New()
│   ├── contributions.go        # GetContributions() logic
│   ├── types.go                # Exported types
│   └── options.go              # Functional options
└── internal/github/            # GitHub API client (private)
    ├── mockResponses/          # Fake github API responses for debug mode
    ├── interface.go            # HTTP client interface
    ├── api.go                  # Real Github HTTP client
    ├── mock_client.go          # Mock Github HTTP client
    ├── ratelimit.go            # Rate limit handling
    └── types.go                # API response types
```

## Development

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Build
go build -o gh-oss-stats ./cmd/gh-oss-stats

# Test locally with debug mode
./gh-oss-stats --user test --debug

# Lint (requires golangci-lint)
golangci-lint run
```
