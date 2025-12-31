package main

import (
	"fmt"
	"os"
)

const version = "0.3.2"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		runMainCmd(args)
		return
	}

	// Route to sub-commands, or fallback to main command
	switch args[0] {
	case "badge":
		runBadgeCmd(args[1:])
	case "demo":
		runDemoCmd(args[1:])
	case "version":
		fmt.Printf("gh-oss-stats v%s\n", version)
		os.Exit(0)
	default:
		runMainCmd(args)
	}
}
