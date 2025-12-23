# Badge Preview Gallery

This directory contains preview badges for documentation purposes.

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

## Sample Data

All preview badges use the following sample data:
- **Username**: mabd-dev
- **Projects**: 42
- **PRs Merged**: 1,567
- **Commits**: 3,284
- **Lines Changed**: 157,450

**Top Contributions:**
1. kubernetes/kubernetes - ⭐ 108K - 45 PRs
2. facebook/react - ⭐ 220K - 38 PRs
3. microsoft/vscode - ⭐ 158K - 32 PRs
4. golang/go - ⭐ 118K - 28 PRs
5. nodejs/node - ⭐ 102K - 25 PRs

## Generating Your Own

To generate badges with your own data:

```bash
gh-oss-stats --user YOUR_USERNAME --badge --badge-style STYLE --badge-theme THEME
```

See [BADGES.md](../BADGES.md) for complete documentation.
