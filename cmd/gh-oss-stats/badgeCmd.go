package main

import (
	"flag"
	"fmt"
	"os"
)

// badgeCmd flag set
var badgeCmd = flag.NewFlagSet("badge", flag.ExitOnError)

// Badge command flags
var (
	badgeFromFile = badgeCmd.String("from-file", "", "Path to stats JSON file")
	badgeData     = badgeCmd.String("data", "", "Stats as JSON string")
)

func init() {
	badgeCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gh-oss-stats badge [options]\n\n")
		fmt.Fprintf(os.Stderr, "Generate badge from existing stats JSON.\n\n")
		fmt.Fprintf(os.Stderr, "This command allows you to generate badges without re-fetching data from GitHub,\n")
		fmt.Fprintf(os.Stderr, "which is useful for creating multiple badge variants from the same stats.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		badgeCmd.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Generate badge from file\n")
		fmt.Fprintf(os.Stderr, "  gh-oss-stats badge --from-file stats.json --badge-style summary\n\n")
		fmt.Fprintf(os.Stderr, "  # Generate badge from JSON string\n")
		fmt.Fprintf(os.Stderr, "  gh-oss-stats badge --data '{\"username\":\"...\",...}' --badge-style compact\n\n")
	}
}

func runBadgeCmd(args []string) {
	// Initialize local badge configuration
	bc := newBadgeConfig()
	bc.registerBadgeFlags(badgeCmd)
	badgeCmd.Parse(args)

	// TODO: Implement badge generation from JSON
	// 1. Validate that either --from-file or --data is provided
	// 2. Load and parse JSON into ossstats.Stats
	// 3. Validate badge options (style, variant, theme, sort)
	// 4. Generate and save badge SVG
	// See mainCmd.go writeBadge() for reference implementation

	fmt.Fprintf(os.Stderr, "Error: 'badge' command not yet implemented\n")
	fmt.Fprintf(os.Stderr, "This will allow generating badges from existing JSON data\n")
	fmt.Fprintf(os.Stderr, "Run 'gh-oss-stats badge --help' for planned usage\n")
	os.Exit(1)
}
