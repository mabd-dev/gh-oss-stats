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
)

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

	generateBadge = flag.Bool("badge", false, "Generate SVG badge")

	debug = flag.Bool("debug", false, "Uses fake data when true")
)

func runMainCmd(args []string) {
	// Initialize local badge configuration
	badgeConfig := newBadgeConfig()
	badgeConfig.registerBadgeFlags(flag.CommandLine)
	flag.Parse()

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
	if badgeConfig.limit <= 0 {
		fmt.Fprintf(os.Stderr, "Error: --badge-limit must be > 0 (got: %d)\n\n", badgeConfig.limit)
		os.Exit(1)
	}
	if *timeoutSec <= 0 {
		fmt.Fprintf(os.Stderr, "Error: --timeout must be > 0 seconds (got: %d)\n\n", *timeoutSec)
		os.Exit(1)
	}

	badgeOption, err := createBadgeOptions(*badgeConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
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

	if strings.TrimSpace(*output) != "" {
		writeStatsToFile(output, stats)
		if *verbose {
			fmt.Fprintf(os.Stderr, "Output written to %s\n", *output)
		}
	} else if *generateBadge {
		if err := writeBadge(badgeOption, badgeConfig.output, verbose, stats); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating badge: %v\n", err)
			os.Exit(1)
		}
	} else { // Write stats to stdout
		jsonData := formatStats(*stats)
		fmt.Println(string(jsonData))
	}
}

func writeStatsToFile(output *string, stats *ossstats.Stats) {
	jsonData := formatStats(*stats)

	if err := os.WriteFile(*output, jsonData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		os.Exit(1)
	}
}

func formatStats(stats ossstats.Stats) []byte {
	jsonData, encodeErr := json.MarshalIndent(stats, "", "  ")

	if encodeErr != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", encodeErr)
		os.Exit(1)
	}
	return jsonData
}
