package badge

import (
	"fmt"
	"strings"
)

var DefaultBadgeTheme = ThemeGithubDark

// BadgeTheme represents the color scheme for the badge
type BadgeTheme string

const (
	ThemeGithubDark          BadgeTheme = "dark"
	ThemeGithubLight         BadgeTheme = "light"
	ThemeDracula             BadgeTheme = "dracula"
	ThemeNord                BadgeTheme = "nord"
	ThemeGruvboxDark         BadgeTheme = "gruvbox-dark"
	ThemeGruvboxLight        BadgeTheme = "gruvbox-light"
	ThemeMonokai             BadgeTheme = "monokai"
	ThemeSolarizedDark       BadgeTheme = "solarized-dark"
	ThemeSolarizedLight      BadgeTheme = "solarized-light"
	ThemeTokyoNight          BadgeTheme = "tokyo-night"
	ThemeTokyoNightStorm     BadgeTheme = "tokyo-night-storm"
	ThemeTokyoNightLight     BadgeTheme = "tokyo-night-light"
	ThemeOneDark             BadgeTheme = "one-dark"
	ThemeOneDarkVivid        BadgeTheme = "one-dark-vivid"
	ThemeCatppuccinMocha     BadgeTheme = "catppuccin-mocha"
	ThemeCatppuccinMacchiato BadgeTheme = "catppuccin-macchiato"
	ThemeCatppuccinFrappe    BadgeTheme = "catppuccin-frappe"
	ThemeCatppuccinLatte     BadgeTheme = "catppuccin-latte"
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
	case "monokai":
		return ThemeMonokai, nil
	case "solarized-dark":
		return ThemeSolarizedDark, nil
	case "solarized-light":
		return ThemeSolarizedLight, nil
	case "tokyo-night":
		return ThemeTokyoNight, nil
	case "tokyo-night-storm":
		return ThemeTokyoNightStorm, nil
	case "tokyo-night-light":
		return ThemeTokyoNightLight, nil
	case "one-dark":
		return ThemeOneDark, nil
	case "one-dark-vivid":
		return ThemeOneDarkVivid, nil
	case "catppuccin-mocha":
		return ThemeCatppuccinMocha, nil
	case "catppuccin-macchiato":
		return ThemeCatppuccinMacchiato, nil
	case "catppuccin-frappe":
		return ThemeCatppuccinFrappe, nil
	case "catppuccin-latte":
		return ThemeCatppuccinLatte, nil
	}

	err := fmt.Errorf("invalid badge theme: %s (check https://github.com/mabd-dev/gh-oss-stats/blob/main/docs/badges/BADGE_THEMES.md for all possible values)", name)
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
	case ThemeMonokai:
		return ThemeColors{
			Background:    "#272822",
			BackgroundAlt: "#3E3D32",
			Text:          "#F8F8F2",
			TextSecondary: "#75715E",
			Border:        "#49483E",
			Accent:        "#66D9EF",
			Positive:      "#A6E22E",
			Negative:      "#F92672",
			Star:          "#E6DB74",
		}
	case ThemeSolarizedDark:
		return ThemeColors{
			Background:    "#002B36",
			BackgroundAlt: "#073642",
			Text:          "#839496",
			TextSecondary: "#586E75",
			Border:        "#073642",
			Accent:        "#268BD2",
			Positive:      "#859900",
			Negative:      "#DC322F",
			Star:          "#B58900",
		}
	case ThemeSolarizedLight:
		return ThemeColors{
			Background:    "#FDF6E3",
			BackgroundAlt: "#EEE8D5",
			Text:          "#657B83",
			TextSecondary: "#93A1A1",
			Border:        "#EEE8D5",
			Accent:        "#268BD2",
			Positive:      "#859900",
			Negative:      "#DC322F",
			Star:          "#B58900",
		}
	case ThemeTokyoNight:
		return ThemeColors{
			Background:    "#1A1B26",
			BackgroundAlt: "#24283B",
			Text:          "#C0CAF5",
			TextSecondary: "#565F89",
			Border:        "#414868",
			Accent:        "#7AA2F7",
			Positive:      "#9ECE6A",
			Negative:      "#F7768E",
			Star:          "#E0AF68",
		}
	case ThemeTokyoNightStorm:
		return ThemeColors{
			Background:    "#24283B",
			BackgroundAlt: "#1F2335",
			Text:          "#C0CAF5",
			TextSecondary: "#565F89",
			Border:        "#414868",
			Accent:        "#7AA2F7",
			Positive:      "#9ECE6A",
			Negative:      "#F7768E",
			Star:          "#E0AF68",
		}
	case ThemeTokyoNightLight:
		return ThemeColors{
			Background:    "#D5D6DB",
			BackgroundAlt: "#CBCCD1",
			Text:          "#343B58",
			TextSecondary: "#6172AF",
			Border:        "#C4C8DA",
			Accent:        "#34548A",
			Positive:      "#485E30",
			Negative:      "#8C4351",
			Star:          "#8F5E15",
		}
	case ThemeOneDark:
		return ThemeColors{
			Background:    "#282C34",
			BackgroundAlt: "#21252B",
			Text:          "#ABB2BF",
			TextSecondary: "#5C6370",
			Border:        "#3E4451",
			Accent:        "#61AFEF",
			Positive:      "#98C379",
			Negative:      "#E06C75",
			Star:          "#E5C07B",
		}
	case ThemeOneDarkVivid:
		return ThemeColors{
			Background:    "#282C34",
			BackgroundAlt: "#21252B",
			Text:          "#ABB2BF",
			TextSecondary: "#5C6370",
			Border:        "#3E4451",
			Accent:        "#528BFF",
			Positive:      "#98C379",
			Negative:      "#EF596F",
			Star:          "#D19A66",
		}
	case ThemeCatppuccinMocha:
		return ThemeColors{
			Background:    "#1E1E2E",
			BackgroundAlt: "#313244",
			Text:          "#CDD6F4",
			TextSecondary: "#A6ADC8",
			Border:        "#45475A",
			Accent:        "#89B4FA",
			Positive:      "#A6E3A1",
			Negative:      "#F38BA8",
			Star:          "#F9E2AF",
		}
	case ThemeCatppuccinMacchiato:
		return ThemeColors{
			Background:    "#24273A",
			BackgroundAlt: "#363A4F",
			Text:          "#CAD3F5",
			TextSecondary: "#A5ADCB",
			Border:        "#494D64",
			Accent:        "#8AADF4",
			Positive:      "#A6DA95",
			Negative:      "#ED8796",
			Star:          "#EED49F",
		}
	case ThemeCatppuccinFrappe:
		return ThemeColors{
			Background:    "#303446",
			BackgroundAlt: "#414559",
			Text:          "#C6D0F5",
			TextSecondary: "#A5ADCE",
			Border:        "#51576D",
			Accent:        "#8CAAEE",
			Positive:      "#A6D189",
			Negative:      "#E78284",
			Star:          "#E5C890",
		}
	case ThemeCatppuccinLatte:
		return ThemeColors{
			Background:    "#EFF1F5",
			BackgroundAlt: "#CCD0DA",
			Text:          "#4C4F69",
			TextSecondary: "#6C6F85",
			Border:        "#BCC0CC",
			Accent:        "#1E66F5",
			Positive:      "#40A02B",
			Negative:      "#D20F39",
			Star:          "#DF8E1D",
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
