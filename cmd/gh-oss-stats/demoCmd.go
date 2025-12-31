package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats"
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
	badgeConfig := newBadgeConfig()
	badgeConfig.registerBadgeFlags(demoCmd)
	demoCmd.Parse(args)

	// TODO: Implement demo badge generation
	// 1. Create sample ossstats.Stats with realistic data
	//    - Include popular repos like kubernetes, react, linux
	//    - Use reasonable PR/commit/star counts
	// 2. Validate badge options (style, variant, theme, sort)
	// 3. Generate and save badge SVG using sample data
	// See mainCmd.go writeBadge() for reference implementation

	var stats = ossstats.Stats{
		Username:    "mabd-dev",
		GeneratedAt: time.Date(2025, 12, 31, 5, 33, 15, 31869_000, time.UTC),
		Summary: ossstats.Summary{
			TotalProjects:  7,
			TotalPRsMerged: 17,
			TotalCommits:   17,
			TotalAdditions: 0,
			TotalDeletions: 0,
		},
		Contributions: []ossstats.Contribution{
			{
				Repo:              "ibad-al-rahman/android-public",
				Owner:             "ibad-al-rahman",
				RepoName:          "android-public",
				Description:       "Android app for Ibad Al-Rahman",
				RepoURL:           "https://github.com/ibad-al-rahman/android-public",
				Stars:             15,
				PRsMerged:         6,
				Commits:           6,
				Additions:         0,
				Deletions:         0,
				FirstContribution: time.Date(2025, 11, 21, 14, 48, 30, 0, time.UTC),
				LastContribution:  time.Date(2025, 12, 17, 5, 14, 39, 0, time.UTC),
			},
			{
				Repo:              "nsh07/Tomato",
				Owner:             "nsh07",
				RepoName:          "Tomato",
				Description:       "Android app for Ibad Al-Rahman",
				RepoURL:           "https://github.com/ibad-al-rahman/android-public",
				Stars:             15,
				PRsMerged:         2,
				Commits:           2,
				Additions:         0,
				Deletions:         0,
				FirstContribution: time.Date(2025, 11, 19, 12, 6, 16, 0, time.UTC),
				LastContribution:  time.Date(2025, 11, 21, 5, 45, 34, 0, time.UTC),
			},
			{
				Repo:              "qamarelsafadi/JetpackComposeTracker",
				Owner:             "qamarelsafadi",
				RepoName:          "JetpackComposeTracker",
				Description:       "Android app for Ibad Al-Rahman",
				RepoURL:           "https://github.com/ibad-al-rahman/android-public",
				Stars:             15,
				PRsMerged:         2,
				Commits:           2,
				Additions:         0,
				Deletions:         0,
				FirstContribution: time.Date(2025, 6, 14, 20, 55, 24, 0, time.UTC),
				LastContribution:  time.Date(2025, 7, 21, 21, 39, 53, 0, time.UTC),
			},
			{
				Repo:              "android/nav3-recipes",
				Owner:             "android",
				RepoName:          "nav3-recipes",
				Description:       "Android app for Ibad Al-Rahman",
				RepoURL:           "https://github.com/ibad-al-rahman/android-public",
				Stars:             15,
				PRsMerged:         2,
				Commits:           2,
				Additions:         0,
				Deletions:         0,
				FirstContribution: time.Date(2025, 6, 9, 18, 24, 56, 0, time.UTC),
				LastContribution:  time.Date(2025, 6, 9, 18, 31, 50, 0, time.UTC),
			},
			{
				Repo:              "android/cahier",
				Owner:             "android",
				RepoName:          "cahier",
				Description:       "Android app for Ibad Al-Rahman",
				RepoURL:           "https://github.com/ibad-al-rahman/android-public",
				Stars:             15,
				PRsMerged:         2,
				Commits:           2,
				Additions:         0,
				Deletions:         0,
				FirstContribution: time.Date(2025, 6, 3, 14, 8, 20, 0, time.UTC),
				LastContribution:  time.Date(2025, 7, 11, 12, 52, 46, 0, time.UTC),
			},
			{
				Repo:              "esatgozcu/Compose-Rolling-Number",
				Owner:             "esatgozcu",
				RepoName:          "Compose-Rolling-Number",
				Description:       "Android app for Ibad Al-Rahman",
				RepoURL:           "https://github.com/ibad-al-rahman/android-public",
				Stars:             15,
				PRsMerged:         2,
				Commits:           2,
				Additions:         0,
				Deletions:         0,
				FirstContribution: time.Date(2025, 2, 17, 15, 46, 51, 0, time.UTC),
				LastContribution:  time.Date(2025, 3, 26, 21, 33, 8, 0, time.UTC),
			},
			{
				Repo:              "zuzmuz/nvimawscli",
				Owner:             "zuzmuz",
				RepoName:          "nvimawscli",
				Description:       "Android app for Ibad Al-Rahman",
				RepoURL:           "https://github.com/ibad-al-rahman/android-public",
				Stars:             15,
				PRsMerged:         1,
				Commits:           1,
				Additions:         0,
				Deletions:         0,
				FirstContribution: time.Date(2024, 5, 6, 20, 37, 13, 0, time.UTC),
				LastContribution:  time.Date(2024, 5, 6, 20, 37, 13, 0, time.UTC),
			},
		},
	}

	badgeOption, err := createBadgeOptions(*badgeConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := writeBadge(
		badgeOption,
		badgeConfig.output,
		verbose,
		&stats,
	); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating badge: %v\n", err)
		os.Exit(1)
	}

	// fmt.Fprintf(os.Stderr, "Error: 'demo' command not yet implemented\n")
	// fmt.Fprintf(os.Stderr, "This will generate badges using sample data for testing/preview\n")
	// fmt.Fprintf(os.Stderr, "Run 'gh-oss-stats demo --help' for planned usage\n")
	// os.Exit(1)
}
