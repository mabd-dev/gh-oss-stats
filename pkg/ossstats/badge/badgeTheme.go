package badge

import (
	"fmt"
	"strings"
)

var DefaultBadgeTheme = ThemeGithubDark

// BadgeTheme represents the color scheme for the badge
type BadgeTheme string

const (
	ThemeGithubDark   BadgeTheme = "dark"
	ThemeGithubLight  BadgeTheme = "light"
	ThemeDracula      BadgeTheme = "dracula"
	ThemeNord         BadgeTheme = "nord"
	ThemeGruvboxDark  BadgeTheme = "gruvbox-dark"
	ThemeGruvboxLight BadgeTheme = "gruvbox-light"
)

func BadgeThemeFromName(name string) (BadgeTheme, error) {
	switch strings.ToLower(name) {
	case "dark":
		return ThemeGithubDark, nil
	case "light":
		return ThemeGithubLight, nil
	case "dracula":
		return ThemeDracula, nil
	case "nord":
		return ThemeNord, nil
	case "gruvbox-dark":
		return ThemeGruvboxDark, nil
	case "gruvbox-light":
		return ThemeGruvboxLight, nil
	}
	err := fmt.Errorf("invalid badge theme: %s (must be: dark, light, dracula, nord, gruvbox-dark, gruvbox-light)", name)
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
	case ThemeGithubLight:
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
			TextSecondary: "#81a1c1",
			Border:        "#4c566a",
			Accent:        "#88c0d0",
			Positive:      "#a3be8c",
			Negative:      "#bf616a",
			Star:          "#ebcb8b",
		}
	case ThemeGruvboxDark:
		return ThemeColors{
			Background:    "#282828",
			BackgroundAlt: "#3c3836",
			Text:          "#ebdbb2",
			TextSecondary: "#a89984",
			Border:        "#3c3836",
			Accent:        "#458588",
			Positive:      "#98971a",
			Negative:      "#cc241d",
			Star:          "#d79921",
		}
	case ThemeGruvboxLight:
		return ThemeColors{
			Background:    "#fbf1c7",
			BackgroundAlt: "#ebdbb2",
			Text:          "#3c3836",
			TextSecondary: "#7c6f64",
			Border:        "#ebdbb2",
			Accent:        "#458588",
			Positive:      "#98971a",
			Negative:      "#cc241d",
			Star:          "#cc241d",
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
