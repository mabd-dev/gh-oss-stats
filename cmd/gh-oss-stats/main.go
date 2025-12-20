package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gh-oss-tools/gh-oss-stats/pkg/ossstats"
)

const version = "1.0.0"

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
		output       = flag.String("output", "", "Output file (default: stdout)")
		outputShort  = flag.String("o", "", "Output file (short)")
		verbose      = flag.Bool("verbose", false, "Verbose logging to stderr")
		verboseShort = flag.Bool("v", false, "Verbose logging (short)")
		timeoutSec   = flag.Int("timeout", int(ossstats.DefaultTimeout.Seconds()), "Timeout in seconds")
		showVersion  = flag.Bool("version", false, "Print version")
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
	}

	if *token != "" {
		opts = append(opts, ossstats.WithToken(*token))
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

	// Encode JSON
	writeStats(output, verbose, stats)
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
