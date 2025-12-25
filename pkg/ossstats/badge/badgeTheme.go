package badge

import (
	"fmt"
	"strings"
)

var DefaultBadgeTheme = ThemeDark

// BadgeTheme represents the color scheme for the badge
type BadgeTheme string

const (
	ThemeDark    BadgeTheme = "dark"  // Github dark
	ThemeLight   BadgeTheme = "light" // Github light
	ThemeDracula BadgeTheme = "dracula"
	ThemeNord    BadgeTheme = "nord"
)

func BadgeThemeFromName(name string) (BadgeTheme, error) {
	switch strings.ToLower(name) {
	case "dark":
		return ThemeDark, nil
	case "light":
		return ThemeLight, nil
	case "dracula":
		return ThemeDracula, nil
	case "nord":
		return ThemeNord, nil
	}
	err := fmt.Errorf("invalid badge theme: %s (must be: dark, light, aurora glass)", name)
	return DefaultBadgeTheme, err
}

// ThemeColors holds the color palette for a theme
type ThemeColors struct {
	// Backgrounds
	Background    string // Main background
	BackgroundAlt string // Cards, boxes, secondary areas

	// Text
	Text          string // Primary text (headings, values)
	TextSecondary string // Labels, muted text

	// UI
	Border string // Borders, dividers
	Accent string // Primary brand color, highlights

	// Semantic (for stats)
	Positive string // Additions, success, growth
	Negative string // Deletions, errors
	Star     string // Star counts (optional, can default to Accent)
}

// GetThemeColors returns the color palette for a given theme
func GetThemeColors(theme BadgeTheme) ThemeColors {
	switch theme {
	case ThemeLight:
		return ThemeColors{
			Background:    "#ffffff",
			BackgroundAlt: "#f6f8fa",
			Text:          "#1f2328",
			TextSecondary: "#656d76",
			Border:        "#d0d7de",
			Accent:        "#0969da",
			Positive:      "#1a7f37",
			Negative:      "#cf222e",
			Star:          "#9a6700",
		}
	case ThemeDracula:
		return ThemeColors{
			Background:    "#282a36",
			BackgroundAlt: "#44475a",
			Text:          "#f8f8f2",
			TextSecondary: "#6272a4",
			Border:        "#44475a",
			Accent:        "#bd93f9",
			Positive:      "#50fa7b",
			Negative:      "#ff5555",
			Star:          "#f1fa8c",
		}
	case ThemeNord:
		return ThemeColors{
			Background:    "#2e3440",
			BackgroundAlt: "#3b4252",
			Text:          "#d8dee9",
			TextSecondary: "#",
			Border:        "#",
			Accent:        "#",
			Positive:      "#",
			Negative:      "#",
			Star:          "#",
		}
	default:
		// Github dark
		return ThemeColors{
			Background:    "#0d1117",
			BackgroundAlt: "#161b22",
			Text:          "#e6edf3",
			TextSecondary: "#8b949e",
			Border:        "#30363d",
			Accent:        "#58a6ff",
			Positive:      "#3fb950",
			Negative:      "#f85149",
			Star:          "#e3b341",
		}

	}
}
