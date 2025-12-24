# gh-oss-stats

A Go library + CLI tool that fetches a GitHub user's open source contributions to external repositories (repos they don't own) and outputs structured JSON.

## Features

- 🔍 **Find External Contributions**: Discovers merged PRs to repositories you don't own
- 📊 **Aggregate Statistics**: Calculates total PRs, commits, lines of code changed
- 🎨 **SVG Badge Generation**: Create beautiful badges in 4 styles (summary, compact, detailed, minimal)
- ⭐ **Repository Filtering**: Filter by minimum star count
- 🚦 **Rate Limit Handling**: Smart rate limit detection with exponential backoff
- 📦 **Library-First Design**: Use as a Go library or standalone CLI

## Usage
Add OSS contribution badge to your github profile in few steps

1. Navigate to your github profile repo
1. Create new file `.github/workflows/generate-oss-badge.yaml`
3. Copy content of [.github/workflows/generate-oss-badge-sample.yaml](.github/workflows/generate-oss-badge-sample.yaml)
4. Commit the changes
5. Reference generated svg image in your `README.md` file
Done

### Workflow Configuration

Samples

| Style |  Output  |
|------------|------------|
| Summary | ![Summary Dark](docs/badges/summary-dark.svg) |
| Detailed | ![Detailed Dark](docs/badges/detailed-dark.svg)  |
| Compact | ![Compact Dark](docs/badges/compact-dark.svg)  |
| Minimal | ![Minimal Dark](docs/badges/minimal-dark.svg)  |


- Change Output Path
You can change generated svg path in lines `39` and `43`

- Change Svg Style And Theme
see [docs/badge/README.md](docs/badge/README.md) for all available options

- How frequent workflow runs:
You can do that at `line 4`

Here is few options
```yaml
# Weekly (Sundays at midnight)
- cron: '0 0 * * 0'  

# Daily (midnight)
- cron: '0 0 * * *' 

# Every 6 hours
- cron: '0 */6 * * *' 

# Hourly
- cron: '0 * * * *' 
```


## Technical Documentation
📖 **Full docs:** See [docs/TECHNICAL.md](docs/TECHNICAL.md)

## License

See [LICENSE](LICENSE) file.

## Contributing

Contributions welcome! Please open an issue or PR.
