package main

import (
	"flag"
	"fmt"
	"os"
)

// demoCmd flag set
var demoCmd = flag.NewFlagSet("demo", flag.ExitOnError)

func init() {
	demoCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gh-oss-stats demo [options]\n\n")
		fmt.Fprintf(os.Stderr, "Generate demo badge using sample data.\n\n")
		fmt.Fprintf(os.Stderr, "This command is useful for:\n")
		fmt.Fprintf(os.Stderr, "  - Testing badge styles and themes without fetching real data\n")
		fmt.Fprintf(os.Stderr, "  - Previewing how badges will look with realistic data\n")
		fmt.Fprintf(os.Stderr, "  - Generating example badges for documentation\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		demoCmd.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Generate demo badge with summary style\n")
		fmt.Fprintf(os.Stderr, "  gh-oss-stats demo --badge-style summary --badge-theme dark\n\n")
		fmt.Fprintf(os.Stderr, "  # Generate compact demo badge\n")
		fmt.Fprintf(os.Stderr, "  gh-oss-stats demo --badge-style compact --badge-output demo.svg\n\n")
	}
}

func runDemoCmd(args []string) {
	// Initialize local badge configuration
	bc := newBadgeConfig()
	bc.registerBadgeFlags(demoCmd)
	demoCmd.Parse(args)

	// TODO: Implement demo badge generation
	// 1. Create sample ossstats.Stats with realistic data
	//    - Include popular repos like kubernetes, react, linux
	//    - Use reasonable PR/commit/star counts
	// 2. Validate badge options (style, variant, theme, sort)
	// 3. Generate and save badge SVG using sample data
	// See mainCmd.go writeBadge() for reference implementation

	fmt.Fprintf(os.Stderr, "Error: 'demo' command not yet implemented\n")
	fmt.Fprintf(os.Stderr, "This will generate badges using sample data for testing/preview\n")
	fmt.Fprintf(os.Stderr, "Run 'gh-oss-stats demo --help' for planned usage\n")
	os.Exit(1)
}
