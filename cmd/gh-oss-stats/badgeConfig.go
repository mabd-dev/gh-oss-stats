package main

import (
	"flag"

	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats/badge"
)

type BadgeConfig struct {
	style   string
	variant string
	theme   string
	output  string
	sort    string
	limit   int
}

// newBadgeConfig creates a new BadgeConfig with default values
func newBadgeConfig() *BadgeConfig {
	return &BadgeConfig{
		style:   string(badge.DefaultBadgeStyle),
		variant: string(badge.DefaultBadgeVariant),
		theme:   string(badge.DefaultBadgeTheme),
		output:  "",
		sort:    string(badge.DefaultSortBy),
		limit:   badge.DefaultPRsLimit,
	}
}

func (bf *BadgeConfig) registerBadgeFlags(fs *flag.FlagSet) {
	fs.StringVar(&bf.style, "badge-style", string(badge.DefaultBadgeStyle), "Badge style: summary, compact, detailed")
	fs.StringVar(&bf.variant, "badge-variant", string(badge.DefaultBadgeVariant), "Badge variants: default, text-based")
	fs.StringVar(&bf.theme, "badge-theme", string(badge.DefaultBadgeTheme), "Badge theme: dark, light, nord, dracula, ...")
	fs.StringVar(&bf.output, "badge-output", "", "Badge output file (default: badge.svg)")
	fs.StringVar(&bf.sort, "badge-sort", string(badge.DefaultSortBy), "Sort contributions by: prs, stars, commits")
	fs.IntVar(&bf.limit, "badge-limit", badge.DefaultPRsLimit, "Number of contributions to show")
}
