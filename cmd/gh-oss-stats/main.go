package main

import (
	"fmt"
	"os"
	"runtime"

	analytics "github.com/mabd-dev/gh-oss-stats/internal/analytic"
)

const version = "0.3.4"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		sendAnalyticsUsage()
		runMainCmd(args)
		return
	}

	// Route to sub-commands, or fallback to main command
	switch args[0] {
	case "badge":
		sendAnalyticsUsage()
		runBadgeCmd(args[1:])
	case "demo":
		sendAnalyticsUsage()
		runDemoCmd(args[1:])
	case "version":
		fmt.Printf("gh-oss-stats v%s\n", version)
		os.Exit(0)
	default:
		sendAnalyticsUsage()
		runMainCmd(args)
	}
}

func sendAnalyticsUsage() {
	// TODO: if first time, and no on CI tell the user we are collecting
	//   + save that user has been told to not tell him again

	isCI := os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != ""
	firstTime := true

	if firstTime && !isCI {
		println("TODO: show data collection message")
	}

	analytics := analytics.CreateAnalytics()

	err := analytics.TrackToolUsage(runtime.GOOS, version, isCI)
	if err != nil {
		print("error sending analytics, err=")
		println(err)
	} else {
		println("analytics sent")
	}
}
