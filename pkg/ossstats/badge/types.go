package badge

// BadgeStyle represents the type of badge to generate
type BadgeStyle string

const (
	StyleSummary  BadgeStyle = "summary"  // 400x200 - Key metrics
	StyleCompact  BadgeStyle = "compact"  // 280x28 - Shields.io style
	StyleDetailed BadgeStyle = "detailed" // 400x320 - Full stats
	StyleMinimal  BadgeStyle = "minimal"  // 120x28 - Project count only
)

// BadgeTheme represents the color scheme for the badge
type BadgeTheme string

const (
	ThemeDark  BadgeTheme = "dark"  // Dark background, light text
	ThemeLight BadgeTheme = "light" // Light background, dark text
)

// ThemeColors holds the color palette for a theme
type ThemeColors struct {
	Background    string
	Text          string
	TextSecondary string
	Border        string
	Accent        string
}

// SortBy represents how contributions should be sorted in detailed view
type SortBy string

const (
	SortByPRs     SortBy = "prs"
	SortByStars   SortBy = "stars"
	SortByCommits SortBy = "commits"
)

// BadgeOptions contains all configuration for badge generation
type BadgeOptions struct {
	Style  BadgeStyle
	Theme  BadgeTheme
	SortBy SortBy // For detailed badge - how to sort contributions (default: prs)
	Limit  int    // For detailed badge - max contributions to show (default: 5)
}

// DefaultThemeDark returns the dark theme color palette
func DefaultThemeDark() ThemeColors {
	return ThemeColors{
		Background:    "#0d1117",
		Text:          "#c9d1d9",
		TextSecondary: "#8b949e",
		Border:        "#30363d",
		Accent:        "#58a6ff",
	}
}

// DefaultThemeLight returns the light theme color palette
func DefaultThemeLight() ThemeColors {
	return ThemeColors{
		Background:    "#ffffff",
		Text:          "#24292f",
		TextSecondary: "#57606a",
		Border:        "#d0d7de",
		Accent:        "#0969da",
	}
}

// GetThemeColors returns the color palette for a given theme
func GetThemeColors(theme BadgeTheme) ThemeColors {
	switch theme {
	case ThemeLight:
		return DefaultThemeLight()
	default:
		return DefaultThemeDark()
	}
}
