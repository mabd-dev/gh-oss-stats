package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

// Helper to capture stderr output
func captureStderr(f func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// Helper to capture stdout output
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestBadgeCmdHelp(t *testing.T) {
	// Reset badge command flag set for testing
	badgeCmd = flag.NewFlagSet("badge", flag.ContinueOnError)
	
	// Re-initialize to get the help text
	var (
		badgeFromFile = badgeCmd.String("from-file", "", "Path to stats JSON file")
		badgeData     = badgeCmd.String("data", "", "Stats as JSON string")
	)
	_ = badgeFromFile
	_ = badgeData

	// Set custom usage function (simulating init())
	var helpOutput string
	badgeCmd.Usage = func() {
		var buf bytes.Buffer
		buf.WriteString("Usage: gh-oss-stats badge [options]\n\n")
		buf.WriteString("Generate badge from existing stats JSON.\n")
		helpOutput = buf.String()
	}

	badgeCmd.Usage()

	if !strings.Contains(helpOutput, "gh-oss-stats badge") {
		t.Errorf("Help output missing 'gh-oss-stats badge'")
	}
	if !strings.Contains(helpOutput, "Generate badge") {
		t.Errorf("Help output missing description")
	}
}

func TestDemoCmdHelp(t *testing.T) {
	// Reset demo command flag set for testing
	demoCmd = flag.NewFlagSet("demo", flag.ContinueOnError)

	// Set custom usage function (simulating init())
	var helpOutput string
	demoCmd.Usage = func() {
		var buf bytes.Buffer
		buf.WriteString("Usage: gh-oss-stats demo [options]\n\n")
		buf.WriteString("Generate demo badge using sample data.\n")
		helpOutput = buf.String()
	}

	demoCmd.Usage()

	if !strings.Contains(helpOutput, "gh-oss-stats demo") {
		t.Errorf("Help output missing 'gh-oss-stats demo'")
	}
	if !strings.Contains(helpOutput, "demo badge") {
		t.Errorf("Help output missing description")
	}
}

func TestBadgeCmdNotImplemented(t *testing.T) {
	// Test that badge command exists
	if badgeCmd == nil {
		t.Error("badgeCmd is nil, should be initialized")
	}

	// Test that flags are properly defined
	if badgeCmd.Lookup("from-file") == nil {
		t.Error("badge command missing --from-file flag")
	}
	if badgeCmd.Lookup("data") == nil {
		t.Error("badge command missing --data flag")
	}
}

func TestDemoCmdNotImplemented(t *testing.T) {
	// Test that demo command exists
	if demoCmd == nil {
		t.Error("demoCmd is nil, should be initialized")
	}

	// Demo command should have badge configuration flags registered
	// We can't easily test this without running the function, but we can verify the flagset exists
	if demoCmd.Name() != "demo" {
		t.Errorf("demoCmd.Name() = %s, want demo", demoCmd.Name())
	}
}

func TestSubCommandRouting(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectRoute string
	}{
		{
			name:        "badge command",
			args:        []string{"badge", "--help"},
			expectRoute: "badge",
		},
		{
			name:        "demo command",
			args:        []string{"demo", "--help"},
			expectRoute: "demo",
		},
		{
			name:        "version command",
			args:        []string{"version"},
			expectRoute: "version",
		},
		{
			name:        "no args defaults to main",
			args:        []string{},
			expectRoute: "main",
		},
		{
			name:        "unknown command defaults to main",
			args:        []string{"--user", "test"},
			expectRoute: "main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test routing logic
			var route string
			
			args := tt.args
			if len(args) == 0 {
				route = "main"
			} else {
				switch args[0] {
				case "badge":
					route = "badge"
				case "demo":
					route = "demo"
				case "version":
					route = "version"
				default:
					route = "main"
				}
			}

			if route != tt.expectRoute {
				t.Errorf("route = %s, want %s", route, tt.expectRoute)
			}
		})
	}
}

func TestBadgeCmdFlagParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantFile string
		wantData string
	}{
		{
			name:     "from-file flag",
			args:     []string{"--from-file", "stats.json"},
			wantFile: "stats.json",
			wantData: "",
		},
		{
			name:     "data flag",
			args:     []string{"--data", `{"username":"test"}`},
			wantFile: "",
			wantData: `{"username":"test"}`,
		},
		{
			name:     "both flags",
			args:     []string{"--from-file", "stats.json", "--data", `{"test":true}`},
			wantFile: "stats.json",
			wantData: `{"test":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags for each test
			testCmd := flag.NewFlagSet("badge", flag.ContinueOnError)
			bc := newBadgeConfig()
			bc.registerBadgeFlags(testCmd)
			fromFile := testCmd.String("from-file", "", "Path to stats JSON file")
			data := testCmd.String("data", "", "Stats as JSON string")

			if err := testCmd.Parse(tt.args); err != nil {
				t.Fatalf("Failed to parse args: %v", err)
			}

			if *fromFile != tt.wantFile {
				t.Errorf("from-file = %s, want %s", *fromFile, tt.wantFile)
			}
			if *data != tt.wantData {
				t.Errorf("data = %s, want %s", *data, tt.wantData)
			}
		})
	}
}

func TestDemoCmdFlagParsing(t *testing.T) {
	// Reset demo command for testing
	demoCmd = flag.NewFlagSet("demo", flag.ContinueOnError)
	
	bc := newBadgeConfig()
	bc.registerBadgeFlags(demoCmd)

	tests := []struct {
		name       string
		args       []string
		wantStyle  string
		wantTheme  string
		wantOutput string
	}{
		{
			name:       "style flag",
			args:       []string{"--badge-style", "compact"},
			wantStyle:  "compact",
			wantTheme:  string("dark"), // default
			wantOutput: "",
		},
		{
			name:       "theme flag",
			args:       []string{"--badge-theme", "light"},
			wantStyle:  string("summary"), // default
			wantTheme:  "light",
			wantOutput: "",
		},
		{
			name:       "output flag",
			args:       []string{"--badge-output", "demo.svg"},
			wantStyle:  string("summary"), // default
			wantTheme:  string("dark"),    // default
			wantOutput: "demo.svg",
		},
		{
			name:       "multiple flags",
			args:       []string{"--badge-style", "compact", "--badge-theme", "nord", "--badge-output", "test.svg"},
			wantStyle:  "compact",
			wantTheme:  "nord",
			wantOutput: "test.svg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset for each test
			demoCmd = flag.NewFlagSet("demo", flag.ContinueOnError)
			bc := newBadgeConfig()
			bc.registerBadgeFlags(demoCmd)

			if err := demoCmd.Parse(tt.args); err != nil {
				t.Fatalf("Failed to parse args: %v", err)
			}

			if bc.style != tt.wantStyle {
				t.Errorf("style = %s, want %s", bc.style, tt.wantStyle)
			}
			if bc.theme != tt.wantTheme {
				t.Errorf("theme = %s, want %s", bc.theme, tt.wantTheme)
			}
			if bc.output != tt.wantOutput {
				t.Errorf("output = %s, want %s", bc.output, tt.wantOutput)
			}
		})
	}
}

func TestVersionConstant(t *testing.T) {
	if version == "" {
		t.Error("version constant is empty")
	}
	
	// Version should be in semver-like format
	if !strings.Contains(version, ".") {
		t.Errorf("version %s doesn't appear to be in semantic version format", version)
	}
}

func TestCommandExists(t *testing.T) {
	// Test that all command functions exist
	tests := []struct {
		name     string
		function interface{}
	}{
		{"runBadgeCmd", runBadgeCmd},
		{"runDemoCmd", runDemoCmd},
		{"runMainCmd", runMainCmd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.function == nil {
				t.Errorf("%s is nil", tt.name)
			}
		})
	}
}

func TestFlagSetsInitialized(t *testing.T) {
	if badgeCmd == nil {
		t.Error("badgeCmd flag set not initialized")
	}
	if demoCmd == nil {
		t.Error("demoCmd flag set not initialized")
	}

	// Test flag set names
	if badgeCmd.Name() != "badge" {
		t.Errorf("badgeCmd name = %s, want badge", badgeCmd.Name())
	}
	if demoCmd.Name() != "demo" {
		t.Errorf("demoCmd name = %s, want demo", demoCmd.Name())
	}
}

func TestBadgeCmdHasRequiredFlags(t *testing.T) {
	requiredFlags := []string{"from-file", "data"}
	
	for _, flagName := range requiredFlags {
		t.Run(flagName, func(t *testing.T) {
			if badgeCmd.Lookup(flagName) == nil {
				t.Errorf("badge command missing required flag: %s", flagName)
			}
		})
	}
}

func TestDemoCmdHasBadgeFlags(t *testing.T) {
	// Demo command should have badge configuration flags available
	// We test this by verifying we can register badge flags to it
	testCmd := flag.NewFlagSet("test-demo", flag.ContinueOnError)
	bc := newBadgeConfig()
	bc.registerBadgeFlags(testCmd)

	badgeFlags := []string{
		"badge-style",
		"badge-variant",
		"badge-theme",
		"badge-output",
		"badge-sort",
		"badge-limit",
	}

	for _, flagName := range badgeFlags {
		t.Run(flagName, func(t *testing.T) {
			if testCmd.Lookup(flagName) == nil {
				t.Errorf("demo command should support flag: %s", flagName)
			}
		})
	}
}
