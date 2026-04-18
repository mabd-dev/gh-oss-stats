package main

import (
	"fmt"
	"os"

	"github.com/mabd-dev/gh-oss-stats/internal/telemetry"
)

const version = "0.3.5"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		telemetry.Send(version, "")
		runMainCmd(args)
		return
	}

	// Route to sub-commands, or fallback to main command
	switch args[0] {
	case "badge":
		telemetry.Send(version, "badge")
		runBadgeCmd(args[1:])
	case "demo":
		telemetry.Send(version, "demo")
		runDemoCmd(args[1:])
	case "version":
		fmt.Printf("gh-oss-stats v%s\n", version)
		os.Exit(0)
	default:
		telemetry.Send(version, "")
		runMainCmd(args)
	}
}
