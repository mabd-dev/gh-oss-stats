package main

import (
	"fmt"
	"os"

	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats"
	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats/badge"
)

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

	return badge.BadgeOptions{
		Style:   badgeStyle,
		Variant: badgeVariant,
		Theme:   badgeTheme,
		SortBy:  badgeSortBy,
		Limit:   conf.limit,
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
