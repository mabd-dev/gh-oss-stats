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

#### Summary Badges (400×200)

| Dark Theme | Light Theme |
|------------|-------------|
| ![Summary Dark](default-summary-dark.svg) | ![Summary Light](default-summary-light.svg) |

#### Compact Badges (280×28)

| Dark Theme | Light Theme |
|------------|-------------|
| ![Compact Dark](default-compact-dark.svg) | ![Compact Light](default-compact-light.svg) |

#### Detailed Badges (400×320)

| Dark Theme | Light Theme |
|------------|-------------|
| ![Detailed Dark](default-detailed-dark.svg) | ![Detailed Light](default-detailed-light.svg) |

#### Minimal Badges (120×28)

| Dark Theme | Light Theme |
|------------|-------------|
| ![Minimal Dark](default-minimal-dark.svg) | ![Minimal Light](default-minimal-light.svg) |


 ## Badge Variants
  Variants control the visual design and layout approach:

| Default Variant |  Text Based Variant  |
|------------|------------|
| ![Default Detailed Dark](default-detailed-dark.svg) | ![Text Based Detailed Dark](text-based-detailed-dark.svg) |


Check [All Combos](./BADGE_THEMES.md)

  ### Default Variant
  Modern, card-based designs with gradients, shadows, and rich visual elements.
  - Best for: Modern GitHub profiles, portfolios
  - Styles available: All (summary, compact, detailed, minimal)

  ### Text-Based Variant
  Clean, minimalist text-focused designs with clear typography.
  - Best for: Terminal themes, retro aesthetics, accessibility
  - Styles available: Detailed only (more coming soon)
  - Characteristics:
    - No gradients or shadows
    - Larger text for better readability
    - Cleaner, more spacious layout
    - Lower file size

  **Usage:**
  ```bash
  # Default variant (rich visual design)
  gh-oss-stats --badge --badge-variant default --badge-style detailed

  # Text-based variant (clean typography)
  gh-oss-stats --badge --badge-variant text-based --badge-style detailed
```
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
