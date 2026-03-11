package main

import (
	"testing"

	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats/badge"
)

func TestCreateBadgeOptionsNoCustomColors(t *testing.T) {
	conf := BadgeConfig{
		style:   string(badge.DefaultBadgeStyle),
		variant: string(badge.DefaultBadgeVariant),
		theme:   string(badge.DefaultBadgeTheme),
		sort:    string(badge.DefaultSortBy),
		limit:   badge.DefaultPRsLimit,
		// all color fields left as zero value ("")
	}

	opts, err := createBadgeOptions(conf)
	if err != nil {
		t.Fatalf("createBadgeOptions failed: %v", err)
	}

	if opts.CustomColors != nil {
		t.Errorf("CustomColors = %+v, want nil when no color flags are set", opts.CustomColors)
	}
}

func TestCreateBadgeOptionsAllCustomColors(t *testing.T) {
	conf := BadgeConfig{
		style:              string(badge.DefaultBadgeStyle),
		variant:            string(badge.DefaultBadgeVariant),
		theme:              string(badge.DefaultBadgeTheme),
		sort:               string(badge.DefaultSortBy),
		limit:              badge.DefaultPRsLimit,
		colorBackground:    "#0d1117",
		colorBackgroundAlt: "#161b22",
		colorText:          "#c9d1d9",
		colorTextSecondary: "#8b949e",
		colorBorder:        "#30363d",
		colorAccent:        "#58a6ff",
		colorPositive:      "#3fb950",
		colorNegative:      "#f85149",
		colorStar:          "#e3b341",
	}

	opts, err := createBadgeOptions(conf)
	if err != nil {
		t.Fatalf("createBadgeOptions failed: %v", err)
	}

	if opts.CustomColors == nil {
		t.Fatal("CustomColors is nil, want non-nil when color flags are set")
	}

	c := opts.CustomColors
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"Background", c.Background, "#0d1117"},
		{"BackgroundAlt", c.BackgroundAlt, "#161b22"},
		{"Text", c.Text, "#c9d1d9"},
		{"TextSecondary", c.TextSecondary, "#8b949e"},
		{"Border", c.Border, "#30363d"},
		{"Accent", c.Accent, "#58a6ff"},
		{"Positive", c.Positive, "#3fb950"},
		{"Negative", c.Negative, "#f85149"},
		{"Star", c.Star, "#e3b341"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("CustomColors.%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}

func TestCreateBadgeOptionsPartialCustomColors(t *testing.T) {
	// Only accent and star are set; all others should be empty strings in CustomColors.
	conf := BadgeConfig{
		style:       string(badge.DefaultBadgeStyle),
		variant:     string(badge.DefaultBadgeVariant),
		theme:       string(badge.DefaultBadgeTheme),
		sort:        string(badge.DefaultSortBy),
		limit:       badge.DefaultPRsLimit,
		colorAccent: "#ff6b6b",
		colorStar:   "#ffd700",
	}

	opts, err := createBadgeOptions(conf)
	if err != nil {
		t.Fatalf("createBadgeOptions failed: %v", err)
	}

	if opts.CustomColors == nil {
		t.Fatal("CustomColors is nil, want non-nil when at least one color flag is set")
	}

	if opts.CustomColors.Accent != "#ff6b6b" {
		t.Errorf("Accent = %q, want %q", opts.CustomColors.Accent, "#ff6b6b")
	}
	if opts.CustomColors.Star != "#ffd700" {
		t.Errorf("Star = %q, want %q", opts.CustomColors.Star, "#ffd700")
	}
	// Unset fields should be empty strings (theme override logic in RenderSVG will skip them).
	if opts.CustomColors.Background != "" {
		t.Errorf("Background = %q, want empty string for unset flag", opts.CustomColors.Background)
	}
}

func TestCreateBadgeOptionsSingleColorFlag(t *testing.T) {
	colorFields := []struct {
		name     string
		makeConf func() BadgeConfig
		check    func(*badge.ThemeColors) (string, string) // returns got, want
	}{
		{
			"background",
			func() BadgeConfig {
				c := baseConf()
				c.colorBackground = "#aabbcc"
				return c
			},
			func(c *badge.ThemeColors) (string, string) { return c.Background, "#aabbcc" },
		},
		{
			"background-alt",
			func() BadgeConfig {
				c := baseConf()
				c.colorBackgroundAlt = "#112233"
				return c
			},
			func(c *badge.ThemeColors) (string, string) { return c.BackgroundAlt, "#112233" },
		},
		{
			"text",
			func() BadgeConfig {
				c := baseConf()
				c.colorText = "#ffffff"
				return c
			},
			func(c *badge.ThemeColors) (string, string) { return c.Text, "#ffffff" },
		},
		{
			"text-secondary",
			func() BadgeConfig {
				c := baseConf()
				c.colorTextSecondary = "#aaaaaa"
				return c
			},
			func(c *badge.ThemeColors) (string, string) { return c.TextSecondary, "#aaaaaa" },
		},
		{
			"border",
			func() BadgeConfig {
				c := baseConf()
				c.colorBorder = "#333333"
				return c
			},
			func(c *badge.ThemeColors) (string, string) { return c.Border, "#333333" },
		},
		{
			"accent",
			func() BadgeConfig {
				c := baseConf()
				c.colorAccent = "#ff0000"
				return c
			},
			func(c *badge.ThemeColors) (string, string) { return c.Accent, "#ff0000" },
		},
		{
			"positive",
			func() BadgeConfig {
				c := baseConf()
				c.colorPositive = "#00ff00"
				return c
			},
			func(c *badge.ThemeColors) (string, string) { return c.Positive, "#00ff00" },
		},
		{
			"negative",
			func() BadgeConfig {
				c := baseConf()
				c.colorNegative = "#ff00ff"
				return c
			},
			func(c *badge.ThemeColors) (string, string) { return c.Negative, "#ff00ff" },
		},
		{
			"star",
			func() BadgeConfig {
				c := baseConf()
				c.colorStar = "#ffff00"
				return c
			},
			func(c *badge.ThemeColors) (string, string) { return c.Star, "#ffff00" },
		},
	}

	for _, tt := range colorFields {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := createBadgeOptions(tt.makeConf())
			if err != nil {
				t.Fatalf("createBadgeOptions failed: %v", err)
			}
			if opts.CustomColors == nil {
				t.Fatal("CustomColors is nil, want non-nil")
			}
			got, want := tt.check(opts.CustomColors)
			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}

// baseConf returns a minimal valid BadgeConfig with no color overrides.
func baseConf() BadgeConfig {
	return BadgeConfig{
		style:   string(badge.DefaultBadgeStyle),
		variant: string(badge.DefaultBadgeVariant),
		theme:   string(badge.DefaultBadgeTheme),
		sort:    string(badge.DefaultSortBy),
		limit:   badge.DefaultPRsLimit,
	}
}

func TestValidateHexColor(t *testing.T) {
	tests := []struct {
		value   string
		wantErr bool
	}{
		// Valid
		{"#abc", false},
		{"#ABC", false},
		{"#1a2b3c", false},
		{"#1A2B3C", false},
		{"#00000000", false}, // 8-digit with alpha
		// Invalid
		{"", true},           // empty (callers skip empty, but function itself rejects)
		{"abc", true},        // missing #
		{"#gg0000", true},    // invalid hex digits
		{"#12345", true},     // 5 digits
		{"#1234567", true},   // 7 digits
		{"#123456789", true}, // 9 digits
		{"red", true},        // named color
		{"rgb(0,0,0)", true}, // CSS function
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			err := validateHexColor("badge-color-accent", tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHexColor(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestValidateHexColorErrorMessage(t *testing.T) {
	err := validateHexColor("badge-color-accent", "notacolor")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	msg := err.Error()
	if !contains(msg, "--badge-color-accent") {
		t.Errorf("error message should mention the flag name, got: %q", msg)
	}
	if !contains(msg, "notacolor") {
		t.Errorf("error message should mention the invalid value, got: %q", msg)
	}
}

func TestCreateBadgeOptionsInvalidColor(t *testing.T) {
	colorFlagFields := []struct {
		flagName string
		setConf  func(*BadgeConfig)
	}{
		{"badge-color-background", func(c *BadgeConfig) { c.colorBackground = "red" }},
		{"badge-color-background-alt", func(c *BadgeConfig) { c.colorBackgroundAlt = "blue" }},
		{"badge-color-text", func(c *BadgeConfig) { c.colorText = "white" }},
		{"badge-color-text-secondary", func(c *BadgeConfig) { c.colorTextSecondary = "gray" }},
		{"badge-color-border", func(c *BadgeConfig) { c.colorBorder = "black" }},
		{"badge-color-accent", func(c *BadgeConfig) { c.colorAccent = "invalid" }},
		{"badge-color-positive", func(c *BadgeConfig) { c.colorPositive = "green" }},
		{"badge-color-negative", func(c *BadgeConfig) { c.colorNegative = "rgb(255,0,0)" }},
		{"badge-color-star", func(c *BadgeConfig) { c.colorStar = "gold" }},
	}

	for _, tt := range colorFlagFields {
		t.Run(tt.flagName, func(t *testing.T) {
			conf := baseConf()
			tt.setConf(&conf)

			_, err := createBadgeOptions(conf)
			if err == nil {
				t.Fatalf("expected error for invalid color on --%s, got nil", tt.flagName)
			}
			if !contains(err.Error(), "--"+tt.flagName) {
				t.Errorf("error should mention the flag name --%s, got: %q", tt.flagName, err.Error())
			}
		})
	}
}

func TestCreateBadgeOptionsValidColors(t *testing.T) {
	validColors := []string{"#fff", "#FFF", "#1a2b3c", "#1A2B3C", "#00000000"}

	for _, color := range validColors {
		t.Run(color, func(t *testing.T) {
			conf := baseConf()
			conf.colorAccent = color

			_, err := createBadgeOptions(conf)
			if err != nil {
				t.Errorf("unexpected error for valid color %q: %v", color, err)
			}
		})
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
