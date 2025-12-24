# gh-oss-stats

> **Showcase your open source contributions with auto-updating GitHub profile badges**

Track and display your merged PRs, commits, and contributions to external repositories with beautiful, auto-updating SVG badges. Perfect for GitHub profiles, portfolios, and resumes.

## ✨ Features

- 🎨 **Auto-Updating Profile Badges** - Beautiful SVG badges in 4 styles (summary, compact, detailed, minimal)
- 🤖 **GitHub Actions Integration** - Set it and forget it, updates weekly automatically
- 🔍 **External Contribution Tracking** - Discovers all your merged PRs to repos you don't own
- 📊 **Comprehensive Stats** - Total PRs, commits, lines of code, and repository stars
- ⭐ **Smart Filtering** - Filter by minimum stars, exclude organizations
- 🎭 **Dark & Light Themes** - Match your profile's aesthetic
- 📦 **Developer-Friendly** - Use as a Go library or standalone CLI, outputs JSON


| Style |  Output  |
|------------|------------|
| Summary | ![Summary Dark](docs/badges/summary-dark.svg) |
| Detailed | ![Detailed Dark](docs/badges/detailed-dark.svg)  |
| Compact | ![Compact Dark](docs/badges/compact-dark.svg)  |
| Minimal | ![Minimal Dark](docs/badges/minimal-dark.svg)  |



## Quick Start

Add an auto-updating OSS contribution badge to your GitHub profile in a few simple steps:

### 1. Create Your Profile Repository
If you don't have one already, create a repository named `USERNAME/USERNAME` (replace USERNAME with your GitHub username). This is your special [profile repository](https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-github-profile/customizing-your-profile/managing-your-profile-readme).

### 2. Set Up the Workflow
1. In your profile repository, create a new file: `.github/workflows/generate-oss-badge.yaml`
2. Copy the content from [this sample workflow](.github/workflows/generate-oss-badge-sample.yaml) and paste it into the file
3. Create an `images/` directory in your repository root (or use a different path in step 4)

### 3. Commit and Wait
Commit the workflow file. The badge will be generated automatically:
- **First run:** Manually trigger via Actions tab, or wait for the scheduled time (Sundays at midnight)
- **Updates:** Automatically every Sunday at midnight (customizable)

### 4. Add Badge to Your Profile
Add this line to your profile `README.md` where you want the badge to appear:

```markdown
![OSS Contributions](images/oss-badge.svg)
```

**Done!** Your badge will auto-update weekly. 🎉

---

## Customization

### Change Badge Style or Theme

Edit the workflow file (`.github/workflows/generate-oss-badge.yaml`) and modify these flags:

```yaml
--badge-style summary    # Options: summary, compact, detailed, minimal
--badge-theme dark       # Options: dark, light
```

See all badge styles and examples in the [Badge Gallery](docs/badges/README.md).

### Change Output Location

In the workflow file, update the path `images/oss-badge.svg` in two places:
1. The `gh-oss-stats` command's `--badge-output` flag
2. The `git add` command

Then update your README.md to reference the new path.

### Change Update Frequency

Modify the `cron` schedule in the workflow file:

```yaml
schedule:
  - cron: '0 0 * * 0'  # Weekly (Sundays at midnight) - default
```

**Common schedules:**
```yaml
- cron: '0 0 * * *'      # Daily at midnight
- cron: '0 */6 * * *'    # Every 6 hours
- cron: '0 0 * * 1'      # Weekly on Mondays
- cron: '0 0 1 * *'      # Monthly on the 1st
```

### Advanced Options

For filtering, sorting, and other advanced options, see [docs/BADGES.md](docs/BADGES.md) and [docs/TECHNICAL.md](docs/TECHNICAL.md)


---

## 👨‍💻 For Developers

While this tool is optimized for GitHub profile badges, it's also a **full-featured Go library and CLI** for programmatic access to contribution data:

- **CLI Usage:** Fetch contribution stats as JSON for your own tools
- **Go Library:** Import `github.com/mabd-dev/gh-oss-stats/pkg/ossstats` in your projects
- **Local Testing:** `--debug` flag for instant testing with mock data

📖 **Full technical docs:** See [docs/TECHNICAL.md](docs/TECHNICAL.md)

## License

See [LICENSE](LICENSE) file.

## Contributing

Contributions welcome! Please open an issue or PR.
