package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats"
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
	badgeConfig := newBadgeConfig()
	badgeConfig.registerBadgeFlags(badgeCmd)
	badgeCmd.Parse(args)

	*badgeFromFile = strings.TrimSpace(*badgeFromFile)
	*badgeData = strings.TrimSpace(*badgeData)

	if *badgeFromFile == "" && *badgeData == "" {
		fmt.Fprintln(os.Stderr, "Error: badgeFromFile or data has to be provided")
		os.Exit(1)
	}

	if *badgeFromFile != "" {
		content, err := os.ReadFile(*badgeFromFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		err = generateBadgeFromJSONString(string(content), *badgeConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to parse json data, error=%v\n", err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *badgeData != "" {
		err := generateBadgeFromJSONString(*badgeData, *badgeConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to parse json data, error=%v\n", err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

}

func generateBadgeFromJSONString(statsJSON string, badgeConfig BadgeConfig) error {
	var stats ossstats.Stats
	err := json.Unmarshal([]byte(statsJSON), &stats)
	if err != nil {
		return err
	}

	badgeOption, err := createBadgeOptions(badgeConfig)
	if err != nil {
		return err
	}

	return writeBadge(badgeOption, badgeConfig.output, verbose, &stats)
}
