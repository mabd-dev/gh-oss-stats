package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats"
	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats/badge"
)

const version = "1.0.0"

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gh-oss-stats [options]\n\n")
		fmt.Fprintf(os.Stderr, "Fetch GitHub OSS contribution statistics and optionally generate SVG badges.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Basic usage\n")
		fmt.Fprintf(os.Stderr, "  gh-oss-stats --user mabd-dev\n\n")
		fmt.Fprintf(os.Stderr, "  # Save JSON output\n")
		fmt.Fprintf(os.Stderr, "  gh-oss-stats --user mabd-dev -o stats.json\n\n")
		fmt.Fprintf(os.Stderr, "  # Generate summary badge\n")
		fmt.Fprintf(os.Stderr, "  gh-oss-stats --user mabd-dev --badge --badge-output badge.svg\n\n")
		fmt.Fprintf(os.Stderr, "  # Generate detailed dark theme badge with top 10 by stars\n")
		fmt.Fprintf(os.Stderr, "  gh-oss-stats --user mabd-dev --badge --badge-style detailed --badge-theme dark --badge-sort stars --badge-limit 10\n\n")
		fmt.Fprintf(os.Stderr, "  # Generate compact light theme badge\n")
		fmt.Fprintf(os.Stderr, "  gh-oss-stats --user mabd-dev --badge --badge-style compact --badge-theme light\n\n")
		fmt.Fprintf(os.Stderr, "Badge Styles:\n")
		fmt.Fprintf(os.Stderr, "  summary  - 400x200px with key metrics (default)\n")
		fmt.Fprintf(os.Stderr, "  compact  - 280x28px shields.io style\n")
		fmt.Fprintf(os.Stderr, "  detailed - 400x320px with top contributions\n")
		fmt.Fprintf(os.Stderr, "  minimal  - 120x28px project count only\n\n")

		fmt.Fprintf(os.Stderr, "Badge Variant (design approach):\n")
		fmt.Fprintf(os.Stderr, "  default     - Modern cards with gradients (all styles supported)\n")
		fmt.Fprintf(os.Stderr, "  text-based  - Clean typography focus (detailed style only)\n")

	}
}

func main() {
	var (
		username     = flag.String("user", "", "GitHub username (required)")
		userShort    = flag.String("u", "", "GitHub username (short)")
		token        = flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub token (default: $GITHUB_TOKEN)")
		tokenShort   = flag.String("t", "", "GitHub token (short)")
		includeLOC   = flag.Bool("include-loc", ossstats.DefaultIncludeLOC, "Include LOC metrics")
		includePRs   = flag.Bool("include-prs", ossstats.DefaultIncludePRDetails, "Include PR details")
		minStars     = flag.Int("min-stars", ossstats.DefaultMinStars, "Minimum repo stars")
		maxPRs       = flag.Int("max-prs", ossstats.DefaultMaxPRS, "Max PRs to fetch")
		excludeOrgs  = flag.String("exclude-orgs", "", "Comma-separated list of organizations to exclude")
		output       = flag.String("output", "", "Output file (default: stdout)")
		outputShort  = flag.String("o", "", "Output file (short)")
		verbose      = flag.Bool("verbose", false, "Verbose logging to stderr")
		verboseShort = flag.Bool("v", false, "Verbose logging (short)")
		timeoutSec   = flag.Int("timeout", int(ossstats.DefaultTimeout.Seconds()), "Timeout in seconds")
		showVersion  = flag.Bool("version", false, "Print version")
		debug        = flag.Bool("debug", false, "Uses fake data when true")

		// Badge generation flags
		generateBadge   = flag.Bool("badge", false, "Generate SVG badge")
		badgeStyleStr   = flag.String("badge-style", string(badge.DefaultBadgeStyle), "Badge style: summary, compact, detailed, minimal")
		badgeVariantStr = flag.String("badge-variant", string(badge.DefaultBadgeVariant), "Badge variants: default, text-based")
		badgeThemeStr   = flag.String("badge-theme", string(badge.DefaultBadgeTheme), "Badge theme: dark, light, nord, dracula, ...")
		badgeOutputStr  = flag.String("badge-output", "", "Badge output file (default: badge.svg)")
		badgeSortStr    = flag.String("badge-sort", string(badge.DefaultSortBy), "Sort contributions by: prs, stars, commits (for detailed badge)")
		badgeLimit      = flag.Int("badge-limit", badge.DefaultPRsLimit, "Number of contributions to show (for detailed badge)")
	)

	flag.Parse()

	if *showVersion {
		fmt.Printf("gh-oss-stats v%s\n", version)
		os.Exit(0)
	}

	// Merge short and long flags
	if *userShort != "" {
		username = userShort
	}
	if *tokenShort != "" {
		token = tokenShort
	}
	if *outputShort != "" {
		output = outputShort
	}
	if *verboseShort {
		*verbose = true
	}

	// Validate required flags
	if *username == "" {
		fmt.Fprintf(os.Stderr, "Error: --user is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Validate numerical flags
	if *minStars < 0 {
		fmt.Fprintf(os.Stderr, "Error: --min-stars must be >= 0 (got: %d)\n\n", *minStars)
		os.Exit(1)
	}
	if *maxPRs <= 0 {
		fmt.Fprintf(os.Stderr, "Error: --max-prs must be > 0 (got: %d)\n\n", *maxPRs)
		os.Exit(1)
	}
	if *badgeLimit <= 0 {
		fmt.Fprintf(os.Stderr, "Error: --badge-limit must be > 0 (got: %d)\n\n", *badgeLimit)
		os.Exit(1)
	}
	if *timeoutSec <= 0 {
		fmt.Fprintf(os.Stderr, "Error: --timeout must be > 0 seconds (got: %d)\n\n", *timeoutSec)
		os.Exit(1)
	}

	// Validate badge options
	badgeStyle, err := badge.BadgeStyleFromName(*badgeStyleStr)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	badgeVariant, err := badge.BadgeVariantFromName(*badgeVariantStr)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	badgeTheme, err := badge.BadgeThemeFromName(*badgeThemeStr)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	badgeSortBy, err := badge.SortByFromName(*badgeSortStr)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Warn if no token provided (not an error, but rate limits will be severe)
	if *token == "" {
		fmt.Fprintf(os.Stderr, "Warning: No GitHub token provided. You'll hit rate limits quickly (60 requests/hour).\n")
		fmt.Fprintf(os.Stderr, "Hint: Set GITHUB_TOKEN environment variable or use --token flag\n")
		fmt.Fprintf(os.Stderr, "      Create a token at: https://github.com/settings/tokens\n\n")
	}

	// Set up logger
	var logger ossstats.Logger
	if *verbose {
		logger = log.New(os.Stderr, "[gh-oss-stats] ", log.LstdFlags)
	}

	// Create client with options
	opts := []ossstats.Option{
		ossstats.WithLOC(*includeLOC),
		ossstats.WithPRDetails(*includePRs),
		ossstats.WithMinStars(*minStars),
		ossstats.WithMaxPRs(*maxPRs),
		ossstats.WithTimeout(time.Duration(*timeoutSec) * time.Second),
		ossstats.WithDebug(*debug),
	}

	if *token != "" {
		opts = append(opts, ossstats.WithToken(*token))
	}

	if *excludeOrgs != "" {
		orgs := strings.Split(*excludeOrgs, ",")
		// Trim whitespace from each org name
		for i, org := range orgs {
			orgs[i] = strings.TrimSpace(org)
		}
		opts = append(opts, ossstats.WithExcludeOrgs(orgs))
	}

	if logger != nil {
		opts = append(opts, ossstats.WithLogger(logger))
	}

	client := ossstats.New(opts...)

	// Fetch contributions
	ctx := context.Background()
	stats, err := client.GetContributions(ctx, *username)

	// Handle errors
	if err != nil {
		// Check for partial results
		if partialErr, ok := err.(*ossstats.ErrPartialResults); ok {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", partialErr)
			stats = partialErr.Stats
		} else if rateLimitErr, ok := err.(*ossstats.ErrRateLimited); ok {
			fmt.Fprintf(os.Stderr, "Error: %v\n", rateLimitErr)
			os.Exit(1)
		} else if authErr, ok := err.(*ossstats.ErrAuthentication); ok {
			fmt.Fprintf(os.Stderr, "Error: %v\n", authErr)
			fmt.Fprintf(os.Stderr, "Hint: Provide a token with --token or set GITHUB_TOKEN\n")
			os.Exit(1)
		} else if notFoundErr, ok := err.(*ossstats.ErrNotFound); ok {
			fmt.Fprintf(os.Stderr, "Error: %v\n", notFoundErr)
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	// Write JSON output
	writeStats(output, verbose, stats)

	// Generate badge if requested
	if *generateBadge {
		if err := writeBadge(badgeStyle, badgeVariant, badgeTheme, badgeSortBy, badgeOutputStr, badgeLimit, verbose, stats); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating badge: %v\n", err)
			os.Exit(1)
		}
	}
}

func writeStats(
	output *string,
	verbose *bool,
	stats *ossstats.Stats,
) {
	jsonData, encodeErr := json.MarshalIndent(stats, "", "  ")

	if encodeErr != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", encodeErr)
		os.Exit(1)
	}

	// Write output
	if *output != "" {
		if err := os.WriteFile(*output, jsonData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "Output written to %s\n", *output)
		}
	} else {
		fmt.Println(string(jsonData))
	}
}

func writeBadge(
	style badge.BadgeStyle,
	variant badge.BadgeVariant,
	theme badge.BadgeTheme,
	sortBy badge.SortBy,
	output *string,
	limit *int,
	verbose *bool,
	stats *ossstats.Stats,
) error {

	if style == badge.StyleMinimal {
		fmt.Fprintf(os.Stderr, "\033[33mWarning: 'minimal' badge style will be removed in 0.3.0\n\033[0m")
	}

	// Create badge options
	opts := badge.BadgeOptions{
		Style:   style,
		Variant: variant,
		Theme:   theme,
		SortBy:  sortBy,
		Limit:   *limit,
	}

	// Generate SVG
	svg, err := badge.RenderSVG(stats, opts)
	if err != nil {
		return fmt.Errorf("failed to render badge: %w", err)
	}

	// Determine output file
	outputFile := *output
	if outputFile == "" {
		outputFile = "badge.svg"
	}

	// Write badge to file
	if err := os.WriteFile(outputFile, []byte(svg), 0644); err != nil {
		return fmt.Errorf("failed to write badge: %w", err)
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Badge written to %s (%s/%s/%s)\n", outputFile, variant, style, theme)
	}

	return nil
}
