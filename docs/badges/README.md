# Badge Generation Guide

Complete guide to generating SVG badges for your OSS contributions.

## Quick Start

```bash
# Generate a badge (uses default: summary style, dark theme)
gh-oss-stats --user YOUR_USERNAME --badge

# Customize style and theme
gh-oss-stats --user YOUR_USERNAME --badge \
  --badge-style compact \
  --badge-theme light \
  --badge-output my-badge.svg
```

## Badge Previews

### Summary Badges (400×200)

| Dark Theme | Light Theme |
|------------|-------------|
| ![Summary Dark](summary-dark.svg) | ![Summary Light](summary-light.svg) |

### Compact Badges (280×28)

| Dark Theme | Light Theme |
|------------|-------------|
| ![Compact Dark](compact-dark.svg) | ![Compact Light](compact-light.svg) |

### Detailed Badges (400×320)

| Dark Theme | Light Theme |
|------------|-------------|
| ![Detailed Dark](detailed-dark.svg) | ![Detailed Light](detailed-light.svg) |

### Minimal Badges (120×28)

| Dark Theme | Light Theme |
|------------|-------------|
| ![Minimal Dark](minimal-dark.svg) | ![Minimal Light](minimal-light.svg) |


---

## Advanced Options

### Sorting (Detailed Badge Only)

Control how contributions are sorted in the detailed badge:

```bash
# Sort by PRs merged (default)
gh-oss-stats --user mabd-dev --badge --badge-style detailed --badge-sort prs

# Sort by repository stars
gh-oss-stats --user mabd-dev --badge --badge-style detailed --badge-sort stars

# Sort by total commits
gh-oss-stats --user mabd-dev --badge --badge-style detailed --badge-sort commits
```

### Contribution Limit (Detailed Badge Only)

Control how many contributions to show:

```bash
# Show top 3 contributions
gh-oss-stats --user mabd-dev --badge --badge-style detailed --badge-limit 3

# Show top 10 contributions
gh-oss-stats --user mabd-dev --badge --badge-style detailed --badge-limit 10
```

### Custom Output Path

```bash
# Save to custom location
gh-oss-stats --user mabd-dev --badge --badge-output path/to/badge.svg

# Generate both JSON stats and badge
gh-oss-stats --user mabd-dev \
  -o stats.json \
  --badge --badge-output badge.svg
```

## Theme Comparison

| Feature | Dark Theme | Light Theme |
|---------|------------|-------------|
| Background | `#0d1117` (GitHub dark) | `#ffffff` (white) |
| Text | `#c9d1d9` (light gray) | `#24292f` (dark gray) |
| Accent | `#58a6ff` (GitHub blue) | `#0969da` (GitHub blue) |
| Use Case | Dark mode sites, modern look | Light mode sites, traditional |

## All CLI Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--badge` | bool | `false` | Enable badge generation |
| `--badge-style` | string | `"summary"` | Badge style: summary, compact, detailed, minimal |
| `--badge-theme` | string | `"dark"` | Color theme: dark, light |
| `--badge-output` | string | `"badge.svg"` | Output file path |
| `--badge-sort` | string | `"prs"` | Sort by: prs, stars, commits (detailed only) |
| `--badge-limit` | int | `5` | Max contributions shown (detailed only) |

## Embedding in README

### Standard Markdown

```markdown
![OSS Contributions](badge.svg)
```

### With Link

```markdown
[![OSS Contributions](badge.svg)](https://github.com/YOUR_USERNAME)
```

### HTML (with size control)

```html
<img src="badge.svg" alt="OSS Contributions" width="400">
```


## Contributing

Found a bug or have a feature request? Open an issue at [gh-oss-stats](https://github.com/mabd-dev/gh-oss-stats/issues).

---

**Need help?** See the main [README](../README.md) or [open an issue](https://github.com/mabd-dev/gh-oss-stats/issues).
