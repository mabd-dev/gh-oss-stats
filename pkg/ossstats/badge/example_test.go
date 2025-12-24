package badge_test

import (
	"fmt"
	"time"

	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats"
	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats/badge"
)

func ExampleRenderSVG_summary() {
	stats := &ossstats.Stats{
		Username:    "mabd-dev",
		GeneratedAt: time.Now(),
		Summary: ossstats.Summary{
			TotalProjects:  42,
			TotalPRsMerged: 156,
			TotalCommits:   328,
			TotalAdditions: 12450,
			TotalDeletions: 3200,
		},
	}

	opts := badge.BadgeOptions{
		Style: badge.StyleSummary,
		Theme: badge.ThemeDark,
	}

	svg, err := badge.RenderSVG(stats, opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Generated summary badge (%d bytes)\n", len(svg))
	// Output: Generated summary badge (1435 bytes)
}

func ExampleRenderSVG_compact() {
	stats := &ossstats.Stats{
		Username:    "mabd-dev",
		GeneratedAt: time.Now(),
		Summary: ossstats.Summary{
			TotalProjects:  42,
			TotalPRsMerged: 156,
			TotalCommits:   328,
		},
	}

	opts := badge.BadgeOptions{
		Style: badge.StyleCompact,
		Theme: badge.ThemeLight,
	}

	svg, err := badge.RenderSVG(stats, opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Generated compact badge (%d bytes)\n", len(svg))
	// Output: Generated compact badge (766 bytes)
}

func ExampleRenderSVG_minimal() {
	stats := &ossstats.Stats{
		Username: "mabd-dev",
		Summary: ossstats.Summary{
			TotalProjects: 42,
		},
	}

	opts := badge.BadgeOptions{
		Style: badge.StyleMinimal,
		Theme: badge.ThemeDark,
	}

	svg, err := badge.RenderSVG(stats, opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Generated minimal badge (%d bytes)\n", len(svg))
	// Output: Generated minimal badge (448 bytes)
}

func ExampleRenderSVG_detailed() {
	stats := &ossstats.Stats{
		Username:    "mabd-dev",
		GeneratedAt: time.Now(),
		Summary: ossstats.Summary{
			TotalProjects:  42,
			TotalPRsMerged: 156,
			TotalCommits:   328,
		},
		Contributions: []ossstats.Contribution{
			{RepoName: "kubernetes/kubernetes", Stars: 108000, PRsMerged: 45},
			{RepoName: "facebook/react", Stars: 220000, PRsMerged: 38},
			{RepoName: "microsoft/vscode", Stars: 158000, PRsMerged: 32},
		},
	}

	opts := badge.BadgeOptions{
		Style:  badge.StyleDetailed,
		Theme:  badge.ThemeDark,
		SortBy: badge.SortByPRs,
		Limit:  3,
	}

	svg, err := badge.RenderSVG(stats, opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Generated detailed badge (%d bytes)\n", len(svg))
	// Output: Generated detailed badge (2158 bytes)
}
