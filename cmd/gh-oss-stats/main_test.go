package main

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats"
)

// resetFlags resets the flag.CommandLine for testing
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

// parseFlags is a helper that extracts the flag parsing logic for testing
func parseFlags(args []string) (*flagConfig, error) {
	resetFlags()

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
		pretty       = flag.Bool("pretty", true, "Pretty-print JSON")
		verbose      = flag.Bool("verbose", false, "Verbose logging to stderr")
		verboseShort = flag.Bool("v", false, "Verbose logging (short)")
		timeoutSec   = flag.Int("timeout", int(ossstats.DefaultTimeout.Seconds()), "Timeout in seconds")
		showVersion  = flag.Bool("version", false, "Print version")
	)

	if err := flag.CommandLine.Parse(args); err != nil {
		return nil, err
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

	return &flagConfig{
		username:    *username,
		token:       *token,
		includeLOC:  *includeLOC,
		includePRs:  *includePRs,
		minStars:    *minStars,
		maxPRs:      *maxPRs,
		output:      *output,
		pretty:      *pretty,
		verbose:     *verbose,
		timeoutSec:  *timeoutSec,
		showVersion: *showVersion,
	}, nil
}

// flagConfig holds parsed flag values for testing
type flagConfig struct {
	username    string
	token       string
	includeLOC  bool
	includePRs  bool
	minStars    int
	maxPRs      int
	output      string
	pretty      bool
	verbose     bool
	timeoutSec  int
	showVersion bool
}

func TestFlagDefaults(t *testing.T) {
	// Clear GITHUB_TOKEN for this test
	oldToken := os.Getenv("GITHUB_TOKEN")
	os.Unsetenv("GITHUB_TOKEN")
	defer func() {
		if oldToken != "" {
			os.Setenv("GITHUB_TOKEN", oldToken)
		}
	}()

	cfg, err := parseFlags([]string{"--user", "testuser"})
	if err != nil {
		t.Fatalf("parseFlags failed: %v", err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"username", cfg.username, "testuser"},
		{"token", cfg.token, ""},
		{"includeLOC", cfg.includeLOC, ossstats.DefaultIncludeLOC},
		{"includePRs", cfg.includePRs, ossstats.DefaultIncludePRDetails},
		{"minStars", cfg.minStars, ossstats.DefaultMinStars},
		{"maxPRs", cfg.maxPRs, ossstats.DefaultMaxPRS},
		{"output", cfg.output, ""},
		{"pretty", cfg.pretty, true},
		{"verbose", cfg.verbose, false},
		{"timeoutSec", cfg.timeoutSec, int(ossstats.DefaultTimeout.Seconds())},
		{"showVersion", cfg.showVersion, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestShortFlagUser(t *testing.T) {
	cfg, err := parseFlags([]string{"-u", "shortuser"})
	if err != nil {
		t.Fatalf("parseFlags failed: %v", err)
	}

	if cfg.username != "shortuser" {
		t.Errorf("username = %q, want %q", cfg.username, "shortuser")
	}
}

func TestShortFlagToken(t *testing.T) {
	cfg, err := parseFlags([]string{"--user", "testuser", "-t", "ghp_shorttoken"})
	if err != nil {
		t.Fatalf("parseFlags failed: %v", err)
	}

	if cfg.token != "ghp_shorttoken" {
		t.Errorf("token = %q, want %q", cfg.token, "ghp_shorttoken")
	}
}

func TestShortFlagOutput(t *testing.T) {
	cfg, err := parseFlags([]string{"--user", "testuser", "-o", "output.json"})
	if err != nil {
		t.Fatalf("parseFlags failed: %v", err)
	}

	if cfg.output != "output.json" {
		t.Errorf("output = %q, want %q", cfg.output, "output.json")
	}
}

func TestShortFlagVerbose(t *testing.T) {
	cfg, err := parseFlags([]string{"--user", "testuser", "-v"})
	if err != nil {
		t.Fatalf("parseFlags failed: %v", err)
	}

	if !cfg.verbose {
		t.Errorf("verbose = %v, want %v", cfg.verbose, true)
	}
}

func TestShortFlagOverridesLongFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "short user overrides long user",
			args: []string{"--user", "longuser", "-u", "shortuser"},
			want: "shortuser",
		},
		{
			name: "short token overrides long token",
			args: []string{"--user", "test", "--token", "longtoken", "-t", "shorttoken"},
			want: "shorttoken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseFlags(tt.args)
			if err != nil {
				t.Fatalf("parseFlags failed: %v", err)
			}

			var got string
			if tt.name == "short user overrides long user" {
				got = cfg.username
			} else {
				got = cfg.token
			}

			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTokenFromEnvironment(t *testing.T) {
	// Set GITHUB_TOKEN for this test
	oldToken := os.Getenv("GITHUB_TOKEN")
	os.Setenv("GITHUB_TOKEN", "ghp_env_token_123")
	defer func() {
		if oldToken != "" {
			os.Setenv("GITHUB_TOKEN", oldToken)
		} else {
			os.Unsetenv("GITHUB_TOKEN")
		}
	}()

	cfg, err := parseFlags([]string{"--user", "testuser"})
	if err != nil {
		t.Fatalf("parseFlags failed: %v", err)
	}

	if cfg.token != "ghp_env_token_123" {
		t.Errorf("token = %q, want %q", cfg.token, "ghp_env_token_123")
	}
}

func TestTokenFlagOverridesEnvironment(t *testing.T) {
	// Set GITHUB_TOKEN for this test
	oldToken := os.Getenv("GITHUB_TOKEN")
	os.Setenv("GITHUB_TOKEN", "ghp_env_token")
	defer func() {
		if oldToken != "" {
			os.Setenv("GITHUB_TOKEN", oldToken)
		} else {
			os.Unsetenv("GITHUB_TOKEN")
		}
	}()

	cfg, err := parseFlags([]string{"--user", "testuser", "--token", "ghp_flag_token"})
	if err != nil {
		t.Fatalf("parseFlags failed: %v", err)
	}

	if cfg.token != "ghp_flag_token" {
		t.Errorf("token = %q, want %q (flag should override env)", cfg.token, "ghp_flag_token")
	}
}

func TestVersionFlag(t *testing.T) {
	cfg, err := parseFlags([]string{"--version"})
	if err != nil {
		t.Fatalf("parseFlags failed: %v", err)
	}

	if !cfg.showVersion {
		t.Errorf("showVersion = %v, want %v", cfg.showVersion, true)
	}
}

func TestBooleanFlags(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		want  bool
		field func(*flagConfig) bool
	}{
		{
			name:  "include-loc true",
			args:  []string{"--user", "test", "--include-loc=true"},
			want:  true,
			field: func(c *flagConfig) bool { return c.includeLOC },
		},
		{
			name:  "include-loc false",
			args:  []string{"--user", "test", "--include-loc=false"},
			want:  false,
			field: func(c *flagConfig) bool { return c.includeLOC },
		},
		{
			name:  "include-prs true",
			args:  []string{"--user", "test", "--include-prs"},
			want:  true,
			field: func(c *flagConfig) bool { return c.includePRs },
		},
		{
			name:  "pretty false",
			args:  []string{"--user", "test", "--pretty=false"},
			want:  false,
			field: func(c *flagConfig) bool { return c.pretty },
		},
		{
			name:  "verbose true",
			args:  []string{"--user", "test", "--verbose"},
			want:  true,
			field: func(c *flagConfig) bool { return c.verbose },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseFlags(tt.args)
			if err != nil {
				t.Fatalf("parseFlags failed: %v", err)
			}

			got := tt.field(cfg)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegerFlags(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		want  int
		field func(*flagConfig) int
	}{
		{
			name:  "min-stars 100",
			args:  []string{"--user", "test", "--min-stars", "100"},
			want:  100,
			field: func(c *flagConfig) int { return c.minStars },
		},
		{
			name:  "min-stars 0",
			args:  []string{"--user", "test", "--min-stars", "0"},
			want:  0,
			field: func(c *flagConfig) int { return c.minStars },
		},
		{
			name:  "max-prs 1000",
			args:  []string{"--user", "test", "--max-prs", "1000"},
			want:  1000,
			field: func(c *flagConfig) int { return c.maxPRs },
		},
		{
			name:  "timeout 600",
			args:  []string{"--user", "test", "--timeout", "600"},
			want:  600,
			field: func(c *flagConfig) int { return c.timeoutSec },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseFlags(tt.args)
			if err != nil {
				t.Fatalf("parseFlags failed: %v", err)
			}

			got := tt.field(cfg)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComplexFlagCombination(t *testing.T) {
	cfg, err := parseFlags([]string{
		"-u", "complexuser",
		"-t", "ghp_token123",
		"--include-loc=true",
		"--include-prs",
		"--min-stars", "50",
		"--max-prs", "200",
		"-o", "results.json",
		"--pretty=false",
		"-v",
		"--timeout", "120",
	})
	if err != nil {
		t.Fatalf("parseFlags failed: %v", err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"username", cfg.username, "complexuser"},
		{"token", cfg.token, "ghp_token123"},
		{"includeLOC", cfg.includeLOC, true},
		{"includePRs", cfg.includePRs, true},
		{"minStars", cfg.minStars, 50},
		{"maxPRs", cfg.maxPRs, 200},
		{"output", cfg.output, "results.json"},
		{"pretty", cfg.pretty, false},
		{"verbose", cfg.verbose, true},
		{"timeoutSec", cfg.timeoutSec, 120},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestTimeoutConversion(t *testing.T) {
	tests := []struct {
		name        string
		timeoutSec  int
		wantSeconds int
	}{
		{"default timeout", int(ossstats.DefaultTimeout.Seconds()), 300},
		{"custom 60 seconds", 60, 60},
		{"custom 600 seconds", 600, 600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := time.Duration(tt.timeoutSec) * time.Second
			got := int(duration.Seconds())
			if got != tt.wantSeconds {
				t.Errorf("got %d seconds, want %d seconds", got, tt.wantSeconds)
			}
		})
	}
}

func TestOutputShortAndLongFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "long output flag",
			args: []string{"--user", "test", "--output", "long.json"},
			want: "long.json",
		},
		{
			name: "short output flag",
			args: []string{"--user", "test", "-o", "short.json"},
			want: "short.json",
		},
		{
			name: "short overrides long output",
			args: []string{"--user", "test", "--output", "long.json", "-o", "short.json"},
			want: "short.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseFlags(tt.args)
			if err != nil {
				t.Fatalf("parseFlags failed: %v", err)
			}

			if cfg.output != tt.want {
				t.Errorf("output = %q, want %q", cfg.output, tt.want)
			}
		})
	}
}

func TestVerboseShortAndLongFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{
			name: "long verbose flag",
			args: []string{"--user", "test", "--verbose"},
			want: true,
		},
		{
			name: "short verbose flag",
			args: []string{"--user", "test", "-v"},
			want: true,
		},
		{
			name: "both verbose flags",
			args: []string{"--user", "test", "--verbose", "-v"},
			want: true,
		},
		{
			name: "no verbose flags",
			args: []string{"--user", "test"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseFlags(tt.args)
			if err != nil {
				t.Fatalf("parseFlags failed: %v", err)
			}

			if cfg.verbose != tt.want {
				t.Errorf("verbose = %v, want %v", cfg.verbose, tt.want)
			}
		})
	}
}

func TestAllFlagsWithLongForm(t *testing.T) {
	cfg, err := parseFlags([]string{
		"--user", "longformuser",
		"--token", "ghp_longtoken",
		"--include-loc=true",
		"--include-prs=true",
		"--min-stars", "75",
		"--max-prs", "300",
		"--output", "output.json",
		"--pretty=false",
		"--verbose=true",
		"--timeout", "180",
	})
	if err != nil {
		t.Fatalf("parseFlags failed: %v", err)
	}

	expected := flagConfig{
		username:    "longformuser",
		token:       "ghp_longtoken",
		includeLOC:  true,
		includePRs:  true,
		minStars:    75,
		maxPRs:      300,
		output:      "output.json",
		pretty:      false,
		verbose:     true,
		timeoutSec:  180,
		showVersion: false,
	}

	if cfg.username != expected.username {
		t.Errorf("username = %q, want %q", cfg.username, expected.username)
	}
	if cfg.token != expected.token {
		t.Errorf("token = %q, want %q", cfg.token, expected.token)
	}
	if cfg.includeLOC != expected.includeLOC {
		t.Errorf("includeLOC = %v, want %v", cfg.includeLOC, expected.includeLOC)
	}
	if cfg.includePRs != expected.includePRs {
		t.Errorf("includePRs = %v, want %v", cfg.includePRs, expected.includePRs)
	}
	if cfg.minStars != expected.minStars {
		t.Errorf("minStars = %d, want %d", cfg.minStars, expected.minStars)
	}
	if cfg.maxPRs != expected.maxPRs {
		t.Errorf("maxPRs = %d, want %d", cfg.maxPRs, expected.maxPRs)
	}
	if cfg.output != expected.output {
		t.Errorf("output = %q, want %q", cfg.output, expected.output)
	}
	if cfg.pretty != expected.pretty {
		t.Errorf("pretty = %v, want %v", cfg.pretty, expected.pretty)
	}
	if cfg.verbose != expected.verbose {
		t.Errorf("verbose = %v, want %v", cfg.verbose, expected.verbose)
	}
	if cfg.timeoutSec != expected.timeoutSec {
		t.Errorf("timeoutSec = %d, want %d", cfg.timeoutSec, expected.timeoutSec)
	}
}

func TestNegativeIntegerValues(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "negative min-stars",
			args:      []string{"--user", "test", "--min-stars", "-10"},
			wantError: false, // flag parsing allows it, validation would catch it
		},
		{
			name:      "negative max-prs",
			args:      []string{"--user", "test", "--max-prs", "-5"},
			wantError: false,
		},
		{
			name:      "negative timeout",
			args:      []string{"--user", "test", "--timeout", "-60"},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseFlags(tt.args)
			hasError := err != nil
			if hasError != tt.wantError {
				t.Errorf("parseFlags error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestEmptyStringValues(t *testing.T) {
	cfg, err := parseFlags([]string{
		"--user", "",
		"--token", "",
		"--output", "",
	})
	if err != nil {
		t.Fatalf("parseFlags failed: %v", err)
	}

	// Empty strings should be preserved
	if cfg.username != "" {
		t.Errorf("username should be empty string, got %q", cfg.username)
	}
	if cfg.token != "" {
		t.Errorf("token should be empty string, got %q", cfg.token)
	}
	if cfg.output != "" {
		t.Errorf("output should be empty string, got %q", cfg.output)
	}
}

func TestInvalidFlagFormat(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "invalid boolean value",
			args: []string{"--user", "test", "--pretty=invalid"},
		},
		{
			name: "invalid integer value",
			args: []string{"--user", "test", "--min-stars", "notanumber"},
		},
		{
			name: "unknown flag",
			args: []string{"--user", "test", "--unknown-flag", "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseFlags(tt.args)
			if err == nil {
				t.Errorf("expected error for invalid flag format, got nil")
			}
		})
	}
}
