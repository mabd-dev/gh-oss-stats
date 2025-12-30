package main

import (
	"fmt"
	"os"
)

const version = "0.3.1"

func main() {
	// Get arguments after program name
	args := os.Args[1:]

	// If no arguments or first arg is not a sub-command, run main command
	if len(args) == 0 {
		runMainCmd(args)
		return
	}

	// Route to sub-commands
	switch args[0] {
	case "badge":
		runBadgeCmd(args[1:])
	case "demo":
		runDemoCmd(args[1:])
	case "version":
		fmt.Printf("gh-oss-stats v%s\n", version)
		os.Exit(0)
	default:
		// If first arg doesn't match a sub-command, treat as main command args
		runMainCmd(args)
	}
}
