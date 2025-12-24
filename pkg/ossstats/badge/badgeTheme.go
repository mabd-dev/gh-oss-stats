package badge

import (
	"fmt"
	"strings"
)

var DefaultBadgeTheme = ThemeDark

// BadgeTheme represents the color scheme for the badge
type BadgeTheme string

const (
	ThemeDark        BadgeTheme = "dark"  // Dark background, light text
	ThemeLight       BadgeTheme = "light" // Light background, dark text
	ThemeAuroraGlass BadgeTheme = "aurora glass"
)

func BadgeThemeFromName(name string) (BadgeTheme, error) {
	switch strings.ToLower(name) {
	case "dark":
		return ThemeDark, nil
	case "light":
		return ThemeLight, nil
	case "aurora glass":
		return ThemeAuroraGlass, nil
	}
	err := fmt.Errorf("invalid badge theme: %s (must be: dark, light)", name)
	return DefaultBadgeTheme, err
}

// ThemeColors holds the color palette for a theme
type ThemeColors struct {
	Background    string
	Text          string
	TextSecondary string
	Border        string
	Accent        string
}

// DefaultThemeDark returns the dark theme color palette
var defaultThemeDark = ThemeColors{
	Background:    "#0d1117",
	Text:          "#c9d1d9",
	TextSecondary: "#8b949e",
	Border:        "#30363d",
	Accent:        "#58a6ff",
}

// GetThemeColors returns the color palette for a given theme
func GetThemeColors(theme BadgeTheme) ThemeColors {
	switch theme {
	case ThemeLight:
		return ThemeColors{
			Background:    "#ffffff",
			Text:          "#24292f",
			TextSecondary: "#57606a",
			Border:        "#d0d7de",
			Accent:        "#0969da",
		}
	case ThemeDark:
		return defaultThemeDark
	case ThemeAuroraGlass:
		return ThemeColors{
			Background:    "#0b1020", // deep blue-black
			Text:          "#e6edf3", // soft white
			TextSecondary: "#9aa4b2", // muted gray-blue
			Border:        "#1f2a44", // cool blue border
			Accent:        "#7dd3fc", // icy blue (primary accent)
		}
	default:
		return defaultThemeDark
	}
}
