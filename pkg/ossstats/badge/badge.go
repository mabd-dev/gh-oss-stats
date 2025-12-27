package badge

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"text/template"

	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats"
	bt "github.com/mabd-dev/gh-oss-stats/pkg/ossstats/badge/badgeTemplates"
)

var DefaultPRsLimit = 5

// templateData holds the data passed to SVG templates
type templateData struct {
	Stats            *ossstats.Stats
	Colors           ThemeColors
	TotalProjects    string
	TotalPRs         string
	TotalCommits     string
	TotalLines       string
	CompactText      string // For compact badge: "n projects | m PRs"
	TopContributions []contributionData
}

// contributionData holds formatted contribution data for templates
type contributionData struct {
	RepoName string
	Stars    string
	PRs      string
}

// RenderSVG generates an SVG badge from the given stats
func RenderSVG(stats *ossstats.Stats, opts BadgeOptions) (string, error) {
	if stats == nil {
		return "", errors.New("stats cannot be nil")
	}

	// Set defaults
	if opts.SortBy == "" {
		opts.SortBy = DefaultSortBy
	}
	if opts.Limit == 0 {
		opts.Limit = DefaultPRsLimit
	}

	// Get theme colors
	colors := GetThemeColors(opts.Theme)

	// Prepare base template data
	data := templateData{
		Stats:         stats,
		Colors:        colors,
		TotalProjects: formatNumber(stats.Summary.TotalProjects),
		TotalPRs:      formatNumber(stats.Summary.TotalPRsMerged),
		TotalCommits:  formatNumber(stats.Summary.TotalCommits),
		TotalLines:    formatNumber(stats.Summary.TotalAdditions + stats.Summary.TotalDeletions),
		CompactText:   fmt.Sprintf("%s projects | %s PRs", formatNumber(stats.Summary.TotalProjects), formatNumber(stats.Summary.TotalPRsMerged)),
	}

	// Add top contributions for detailed view
	if opts.Style == StyleDetailed {
		data.TopContributions = getTopContributions(stats, opts.SortBy, opts.Limit)
	}

	// Select template based on style
	tmplStr, err := getTemplateStr(opts.Style, opts.Variant)
	if err != nil {
		return "", err
	}

	// Parse and execute template with custom functions
	tmpl, err := template.New("badge").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"mod": func(a, b int) int { return a % b },
		"div": func(a, b int) int { return a / b },
	}).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// formatNumber formats an integer with appropriate suffix (K, M)
func formatNumber(n int) string {
	if n >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	}
	return fmt.Sprintf("%d", n)
}

// formatStars formats a star count with appropriate suffix
func formatStars(n int) string {
	return formatNumber(n)
}

// getTopContributions returns the top N contributions sorted by the specified criteria
func getTopContributions(stats *ossstats.Stats, sortBy SortBy, limit int) []contributionData {
	// Make a copy of contributions for sorting
	contributions := make([]ossstats.Contribution, len(stats.Contributions))
	copy(contributions, stats.Contributions)

	// Sort based on criteria
	sort.Slice(contributions, func(i, j int) bool {
		switch sortBy {
		case SortByStars:
			return contributions[i].Stars > contributions[j].Stars
		case SortByCommits:
			return contributions[i].Commits > contributions[j].Commits
		case SortByPRs:
			fallthrough
		default:
			return contributions[i].PRsMerged > contributions[j].PRsMerged
		}
	})

	// Get top N
	if limit > len(contributions) {
		limit = len(contributions)
	}

	// Format for template
	result := make([]contributionData, limit)
	for i := 0; i < limit; i++ {
		result[i] = contributionData{
			RepoName: contributions[i].RepoName,
			Stars:    formatStars(contributions[i].Stars),
			PRs:      formatNumber(contributions[i].PRsMerged),
		}
	}

	return result
}

func getTemplateStr(
	style BadgeStyle,
	variant BadgeVariant,
) (string, error) {
	switch variant {
	case VariantDefault:
		switch style {
		case StyleSummary:
			return bt.DefaultSummary, nil
		case StyleCompact:
			return bt.DefaultCompact, nil
		case StyleDetailed:
			return bt.DefaultDetailed, nil
		}
	case VariantTextBased:
		switch style {
		case StyleDetailed:
			return bt.TextBasedDetailed, nil
		}
	}

	err := fmt.Errorf("unsupported badge variant: %s, and style: %s combinations", variant, style)
	return "", err
}
