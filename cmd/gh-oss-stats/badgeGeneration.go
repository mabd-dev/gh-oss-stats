package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats"
	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats/badge"
)

var hexColorRE = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)

func validateHexColor(flag, value string) error {
	if !hexColorRE.MatchString(value) {
		return fmt.Errorf("invalid color for --%s: %q (expected hex format, e.g. #rgb, #rrggbb, #rrggbbaa)", flag, value)
	}
	return nil
}

func createBadgeOptions(conf BadgeConfig) (badge.BadgeOptions, error) {
	badgeStyle, err := badge.BadgeStyleFromName(conf.style)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	badgeVariant, err := badge.BadgeVariantFromName(conf.variant)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	badgeTheme, err := badge.BadgeThemeFromName(conf.theme)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	badgeSortBy, err := badge.SortByFromName(conf.sort)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	colorFlags := []struct {
		flag  string
		value string
	}{
		{"badge-color-background", conf.colorBackground},
		{"badge-color-background-alt", conf.colorBackgroundAlt},
		{"badge-color-text", conf.colorText},
		{"badge-color-text-secondary", conf.colorTextSecondary},
		{"badge-color-border", conf.colorBorder},
		{"badge-color-accent", conf.colorAccent},
		{"badge-color-positive", conf.colorPositive},
		{"badge-color-negative", conf.colorNegative},
		{"badge-color-star", conf.colorStar},
	}
	for _, cf := range colorFlags {
		if cf.value != "" {
			if err := validateHexColor(cf.flag, cf.value); err != nil {
				return badge.BadgeOptions{}, err
			}
		}
	}

	var customColors *badge.ThemeColors
	if conf.colorBackground != "" || conf.colorBackgroundAlt != "" ||
		conf.colorText != "" || conf.colorTextSecondary != "" ||
		conf.colorBorder != "" || conf.colorAccent != "" ||
		conf.colorPositive != "" || conf.colorNegative != "" || conf.colorStar != "" {
		customColors = &badge.ThemeColors{
			Background:    conf.colorBackground,
			BackgroundAlt: conf.colorBackgroundAlt,
			Text:          conf.colorText,
			TextSecondary: conf.colorTextSecondary,
			Border:        conf.colorBorder,
			Accent:        conf.colorAccent,
			Positive:      conf.colorPositive,
			Negative:      conf.colorNegative,
			Star:          conf.colorStar,
		}
	}

	return badge.BadgeOptions{
		Style:        badgeStyle,
		Variant:      badgeVariant,
		Theme:        badgeTheme,
		SortBy:       badgeSortBy,
		Limit:        conf.limit,
		CustomColors: customColors,
	}, nil
}

func writeBadge(
	opts badge.BadgeOptions,
	output string,
	verbose *bool,
	stats *ossstats.Stats,
) error {
	// Generate SVG
	svg, err := badge.RenderSVG(stats, opts)
	if err != nil {
		return fmt.Errorf("failed to render badge: %w", err)
	}

	// Determine output file
	outputFile := output
	if outputFile == "" {
		outputFile = "badge.svg"
	}

	// Write badge to file
	if err := os.WriteFile(outputFile, []byte(svg), 0644); err != nil {
		return fmt.Errorf("failed to write badge: %w", err)
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Badge written to %s (%s/%s/%s)\n", outputFile, opts.Variant, opts.Style, opts.Theme)
	}

	return nil
}
