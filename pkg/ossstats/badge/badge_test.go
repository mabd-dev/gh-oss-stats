package badge

import (
	"strings"
	"testing"
	"time"

	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats"
)

func TestRenderSVG_AllStyles(t *testing.T) {
	stats := &ossstats.Stats{
		Username:    "testuser",
		GeneratedAt: time.Now(),
		Summary: ossstats.Summary{
			TotalProjects:  42,
			TotalPRsMerged: 1567,
			TotalCommits:   3284,
			TotalAdditions: 125450,
			TotalDeletions: 32000,
		},
		Contributions: []ossstats.Contribution{
			{RepoName: "kubernetes/kubernetes", Stars: 108000, PRsMerged: 45, Commits: 120},
			{RepoName: "facebook/react", Stars: 220000, PRsMerged: 38, Commits: 95},
			{RepoName: "microsoft/vscode", Stars: 158000, PRsMerged: 32, Commits: 87},
		},
	}

	tests := []struct {
		name      string
		style     BadgeStyle
		variant   BadgeVariant
		theme     BadgeTheme
		wantWidth string
		wantErr   bool
	}{
		{
			name:      "summary_dark",
			style:     StyleSummary,
			variant:   VariantDefault,
			theme:     ThemeGithubDark,
			wantWidth: `width="400"`,
			wantErr:   false,
		},
		{
			name:      "summary_light",
			style:     StyleSummary,
			variant:   VariantDefault,
			theme:     ThemeGithubLight,
			wantWidth: `width="400"`,
			wantErr:   false,
		},
		{
			name:      "compact_dark",
			style:     StyleCompact,
			variant:   VariantDefault,
			theme:     ThemeGithubDark,
			wantWidth: `width="280"`,
			wantErr:   false,
		},
		{
			name:      "compact_light",
			style:     StyleCompact,
			variant:   VariantDefault,
			theme:     ThemeGithubLight,
			wantWidth: `width="280"`,
			wantErr:   false,
		},
		{
			name:      "detailed_dark",
			style:     StyleDetailed,
			variant:   VariantDefault,
			theme:     ThemeGithubDark,
			wantWidth: `width="900"`,
			wantErr:   false,
		},
		{
			name:      "detailed_light",
			style:     StyleDetailed,
			variant:   VariantDefault,
			theme:     ThemeGithubLight,
			wantWidth: `width="900"`,
			wantErr:   false,
		},
		{
			name:      "minimal_dark",
			style:     StyleMinimal,
			variant:   VariantDefault,
			theme:     ThemeGithubDark,
			wantWidth: `width="120"`,
			wantErr:   false,
		},
		{
			name:      "minimal_light",
			style:     StyleMinimal,
			variant:   VariantDefault,
			theme:     ThemeGithubLight,
			wantWidth: `width="120"`,
			wantErr:   false,
		},
		{
			name:      "summary_dracula",
			style:     StyleSummary,
			variant:   VariantDefault,
			theme:     ThemeDracula,
			wantWidth: `width="400"`,
			wantErr:   false,
		},
		{
			name:      "compact_nord",
			style:     StyleCompact,
			variant:   VariantDefault,
			theme:     ThemeNord,
			wantWidth: `width="280"`,
			wantErr:   false,
		},
		{
			name:      "detailed_nord",
			style:     StyleDetailed,
			variant:   VariantTextBased,
			theme:     ThemeNord,
			wantWidth: `width="720"`,
			wantErr:   false,
		},
		{
			name:      "compact_nord",
			style:     StyleCompact,
			variant:   VariantTextBased,
			theme:     ThemeNord,
			wantWidth: `width="720"`,
			wantErr:   true,
		},
		{
			name:      "summary_nord",
			style:     StyleSummary,
			variant:   VariantTextBased,
			theme:     ThemeNord,
			wantWidth: `width="720"`,
			wantErr:   true,
		},
		{
			name:      "minimal_nord",
			style:     StyleMinimal,
			variant:   VariantTextBased,
			theme:     ThemeNord,
			wantWidth: `width="720"`,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := BadgeOptions{
				Style:   tt.style,
				Theme:   tt.theme,
				Variant: tt.variant,
			}

			svg, err := RenderSVG(stats, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderSVG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !strings.Contains(svg, tt.wantWidth) {
					t.Errorf("RenderSVG() missing expected width %s", tt.wantWidth)
				}
				if !strings.Contains(svg, "<svg") {
					t.Errorf("RenderSVG() doesn't contain <svg tag")
				}
				if !strings.Contains(svg, "</svg>") {
					t.Errorf("RenderSVG() doesn't contain closing </svg> tag")
				}
			}
		})
	}
}

func TestRenderSVG_NilStats(t *testing.T) {
	opts := BadgeOptions{
		Style:   StyleSummary,
		Variant: VariantDefault,
		Theme:   ThemeGithubDark,
	}

	_, err := RenderSVG(nil, opts)
	if err == nil {
		t.Error("RenderSVG() expected error for nil stats, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "cannot be nil") {
		t.Errorf("RenderSVG() unexpected error message: %v", err)
	}
}

func TestRenderSVG_EmptyContributions(t *testing.T) {
	stats := &ossstats.Stats{
		Username: "testuser",
		Summary: ossstats.Summary{
			TotalProjects:  0,
			TotalPRsMerged: 0,
			TotalCommits:   0,
		},
		Contributions: []ossstats.Contribution{},
	}

	tests := []BadgeStyle{StyleSummary, StyleCompact, StyleDetailed, StyleMinimal}

	for _, style := range tests {
		t.Run(string(style), func(t *testing.T) {
			opts := BadgeOptions{
				Style:   style,
				Variant: VariantDefault,
				Theme:   ThemeGithubDark,
			}

			svg, err := RenderSVG(stats, opts)
			if err != nil {
				t.Errorf("RenderSVG() unexpected error for empty contributions: %v", err)
				return
			}

			if !strings.Contains(svg, "<svg") {
				t.Error("RenderSVG() should still generate valid SVG for empty stats")
			}
		})
	}
}

func TestRenderSVG_InvalidStyle(t *testing.T) {
	stats := &ossstats.Stats{
		Username: "testuser",
		Summary:  ossstats.Summary{TotalProjects: 1},
	}

	opts := BadgeOptions{
		Style:   BadgeStyle("invalid"),
		Variant: VariantDefault,
		Theme:   ThemeGithubDark,
	}

	_, err := RenderSVG(stats, opts)
	if err == nil {
		t.Error("RenderSVG() expected error for invalid style, got nil")
	}
}

func TestRenderSVG_ThemeColors(t *testing.T) {
	stats := &ossstats.Stats{
		Username: "testuser",
		Summary:  ossstats.Summary{TotalProjects: 42},
	}

	tests := []struct {
		name       string
		style      BadgeStyle
		theme      BadgeTheme
		wantColor  string
		wantAccent string
	}{
		{
			name:       "dark_theme_compact",
			style:      StyleCompact,
			theme:      ThemeGithubDark,
			wantColor:  "#e6edf3", // text color
			wantAccent: "#58a6ff",
		},
		{
			name:       "light_theme_compact",
			style:      StyleCompact,
			theme:      ThemeGithubLight,
			wantColor:  "#1f2328", // text color
			wantAccent: "#0969da",
		},
		{
			name:       "dark_theme_summary",
			style:      StyleSummary,
			theme:      ThemeGithubDark,
			wantColor:  "#0d1117", // background color
			wantAccent: "",        // summary doesn't use accent in visible way
		},
		{
			name:       "light_theme_summary",
			style:      StyleSummary,
			theme:      ThemeGithubLight,
			wantColor:  "#ffffff", // background color
			wantAccent: "",        // summary doesn't use accent in visible way
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := BadgeOptions{
				Style:   tt.style,
				Variant: VariantDefault,
				Theme:   tt.theme,
			}

			svg, err := RenderSVG(stats, opts)
			if err != nil {
				t.Fatalf("RenderSVG() unexpected error: %v", err)
			}

			if !strings.Contains(svg, tt.wantColor) {
				t.Errorf("RenderSVG() missing theme color %s", tt.wantColor)
			}
			if tt.wantAccent != "" && !strings.Contains(svg, tt.wantAccent) {
				t.Errorf("RenderSVG() missing accent color %s", tt.wantAccent)
			}
		})
	}
}

func TestRenderSVG_DetailedBadgeSorting(t *testing.T) {
	stats := &ossstats.Stats{
		Username: "testuser",
		Summary: ossstats.Summary{
			TotalProjects:  3,
			TotalPRsMerged: 100,
			TotalCommits:   200,
		},
		Contributions: []ossstats.Contribution{
			{RepoName: "repo-a", Stars: 1000, PRsMerged: 5, Commits: 50},
			{RepoName: "repo-b", Stars: 500, PRsMerged: 10, Commits: 30},
			{RepoName: "repo-c", Stars: 2000, PRsMerged: 3, Commits: 100},
		},
	}

	tests := []struct {
		name      string
		sortBy    SortBy
		wantFirst string
	}{
		{
			name:      "sort_by_prs",
			sortBy:    SortByPRs,
			wantFirst: "repo-b", // 10 PRs
		},
		{
			name:      "sort_by_stars",
			sortBy:    SortByStars,
			wantFirst: "repo-c", // 2000 stars
		},
		{
			name:      "sort_by_commits",
			sortBy:    SortByCommits,
			wantFirst: "repo-c", // 100 commits
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := BadgeOptions{
				Style:   StyleDetailed,
				Variant: VariantDefault,
				Theme:   ThemeGithubDark,
				SortBy:  tt.sortBy,
				Limit:   3,
			}

			svg, err := RenderSVG(stats, opts)
			if err != nil {
				t.Fatalf("RenderSVG() unexpected error: %v", err)
			}

			if !strings.Contains(svg, tt.wantFirst) {
				t.Errorf("RenderSVG() expected first repo to be %s", tt.wantFirst)
			}

			// Check that the first occurrence of wantFirst comes before others
			firstIdx := strings.Index(svg, tt.wantFirst)
			if firstIdx == -1 {
				t.Errorf("RenderSVG() missing repo %s", tt.wantFirst)
			}
		})
	}
}

func TestRenderSVG_DetailedBadgeLimit(t *testing.T) {
	contributions := make([]ossstats.Contribution, 10)
	for i := 0; i < 10; i++ {
		contributions[i] = ossstats.Contribution{
			RepoName:  string(rune('a' + i)),
			Stars:     1000 * (10 - i),
			PRsMerged: 10 - i,
			Commits:   100,
		}
	}

	stats := &ossstats.Stats{
		Username:      "testuser",
		Summary:       ossstats.Summary{TotalProjects: 10},
		Contributions: contributions,
	}

	tests := []struct {
		name      string
		limit     int
		wantCount int
	}{
		{
			name:      "limit_3",
			limit:     3,
			wantCount: 3,
		},
		{
			name:      "limit_5",
			limit:     5,
			wantCount: 5,
		},
		{
			name:      "limit_10",
			limit:     10,
			wantCount: 10,
		},
		{
			name:      "limit_exceeds_total",
			limit:     20,
			wantCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := BadgeOptions{
				Style:   StyleDetailed,
				Variant: VariantDefault,
				Theme:   ThemeGithubDark,
				SortBy:  SortByPRs,
				Limit:   tt.limit,
			}

			svg, err := RenderSVG(stats, opts)
			if err != nil {
				t.Fatalf("RenderSVG() unexpected error: %v", err)
			}

			// Count occurrences of repo-name class
			count := strings.Count(svg, `class="repo-name"`)
			if count != tt.wantCount {
				t.Errorf("RenderSVG() got %d repos, want %d", count, tt.wantCount)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  string
	}{
		{
			name:  "zero",
			input: 0,
			want:  "0",
		},
		{
			name:  "small_number",
			input: 42,
			want:  "42",
		},
		{
			name:  "hundreds",
			input: 999,
			want:  "999",
		},
		{
			name:  "thousands",
			input: 1000,
			want:  "1.0K",
		},
		{
			name:  "thousands_mid",
			input: 1567,
			want:  "1.6K",
		},
		{
			name:  "tens_of_thousands",
			input: 42000,
			want:  "42.0K",
		},
		{
			name:  "hundreds_of_thousands",
			input: 125450,
			want:  "125.5K",
		},
		{
			name:  "millions",
			input: 1000000,
			want:  "1.0M",
		},
		{
			name:  "millions_mid",
			input: 2500000,
			want:  "2.5M",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatNumber(tt.input)
			if got != tt.want {
				t.Errorf("formatNumber(%d) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetTopContributions(t *testing.T) {
	stats := &ossstats.Stats{
		Username: "testuser",
		Contributions: []ossstats.Contribution{
			{RepoName: "repo-a", Stars: 1000, PRsMerged: 10, Commits: 50},
			{RepoName: "repo-b", Stars: 2000, PRsMerged: 5, Commits: 100},
			{RepoName: "repo-c", Stars: 500, PRsMerged: 20, Commits: 30},
		},
	}

	tests := []struct {
		name      string
		sortBy    SortBy
		limit     int
		wantFirst string
		wantLen   int
	}{
		{
			name:      "sort_by_prs_limit_2",
			sortBy:    SortByPRs,
			limit:     2,
			wantFirst: "repo-c",
			wantLen:   2,
		},
		{
			name:      "sort_by_stars_limit_1",
			sortBy:    SortByStars,
			limit:     1,
			wantFirst: "repo-b",
			wantLen:   1,
		},
		{
			name:      "sort_by_commits_all",
			sortBy:    SortByCommits,
			limit:     10,
			wantFirst: "repo-b",
			wantLen:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTopContributions(stats, tt.sortBy, tt.limit)

			if len(got) != tt.wantLen {
				t.Errorf("getTopContributions() len = %d, want %d", len(got), tt.wantLen)
			}

			if len(got) > 0 && got[0].RepoName != tt.wantFirst {
				t.Errorf("getTopContributions() first = %s, want %s", got[0].RepoName, tt.wantFirst)
			}
		})
	}
}

func TestGetThemeColors(t *testing.T) {
	tests := []struct {
		name           string
		theme          BadgeTheme
		wantBackground string
		wantAccent     string
	}{
		{
			name:           "dark_theme",
			theme:          ThemeGithubDark,
			wantBackground: "#0d1117",
			wantAccent:     "#58a6ff",
		},
		{
			name:           "light_theme",
			theme:          ThemeGithubLight,
			wantBackground: "#ffffff",
			wantAccent:     "#0969da",
		},
		{
			name:           "dracula_theme",
			theme:          ThemeDracula,
			wantBackground: "#282a36",
			wantAccent:     "#bd93f9",
		},
		{
			name:           "nord_theme",
			theme:          ThemeNord,
			wantBackground: "#2e3440",
			wantAccent:     "#88c0d0",
		},
		{
			name:           "unknown_theme_defaults_to_dark",
			theme:          BadgeTheme("unknown"),
			wantBackground: "#0d1117",
			wantAccent:     "#58a6ff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colors := GetThemeColors(tt.theme)

			if colors.Background != tt.wantBackground {
				t.Errorf("GetThemeColors().Background = %s, want %s", colors.Background, tt.wantBackground)
			}
			if colors.Accent != tt.wantAccent {
				t.Errorf("GetThemeColors().Accent = %s, want %s", colors.Accent, tt.wantAccent)
			}
		})
	}
}

func TestRenderSVG_CompactBadgeContent(t *testing.T) {
	stats := &ossstats.Stats{
		Username: "testuser",
		Summary: ossstats.Summary{
			TotalProjects:  42,
			TotalPRsMerged: 1567,
		},
	}

	opts := BadgeOptions{
		Style:   StyleCompact,
		Variant: VariantDefault,
		Theme:   ThemeGithubDark,
	}

	svg, err := RenderSVG(stats, opts)
	if err != nil {
		t.Fatalf("RenderSVG() unexpected error: %v", err)
	}

	// Check for expected content
	if !strings.Contains(svg, "42 projects") {
		t.Error("Compact badge missing '42 projects'")
	}
	if !strings.Contains(svg, "1.6K PRs") {
		t.Error("Compact badge missing '1.6K PRs'")
	}
}

func TestRenderSVG_MinimalBadgeContent(t *testing.T) {
	stats := &ossstats.Stats{
		Username: "testuser",
		Summary: ossstats.Summary{
			TotalProjects: 42,
		},
	}

	opts := BadgeOptions{
		Style:   StyleMinimal,
		Variant: VariantDefault,
		Theme:   ThemeGithubDark,
	}

	svg, err := RenderSVG(stats, opts)
	if err != nil {
		t.Fatalf("RenderSVG() unexpected error: %v", err)
	}

	// Check for expected content
	if !strings.Contains(svg, "42 Projects") {
		t.Error("Minimal badge missing '42 Projects'")
	}
}
