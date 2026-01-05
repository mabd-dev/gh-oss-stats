package main

import (
	"flag"
	"testing"

	"github.com/mabd-dev/gh-oss-stats/pkg/ossstats/badge"
)

func TestNewBadgeConfig(t *testing.T) {
	bc := newBadgeConfig()

	if bc == nil {
		t.Fatal("newBadgeConfig() returned nil")
	}

	// Test defaults by registering flags and checking their default values
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	bc.registerBadgeFlags(fs)

	tests := []struct {
		name    string
		gotFunc func() interface{}
		want    interface{}
	}{
		{"style default", func() interface{} { return fs.Lookup("badge-style").DefValue }, string(badge.DefaultBadgeStyle)},
		{"variant default", func() interface{} { return fs.Lookup("badge-variant").DefValue }, string(badge.DefaultBadgeVariant)},
		{"theme default", func() interface{} { return fs.Lookup("badge-theme").DefValue }, string(badge.DefaultBadgeTheme)},
		{"output default", func() interface{} { return fs.Lookup("badge-output").DefValue }, ""},
		{"sort default", func() interface{} { return fs.Lookup("badge-sort").DefValue }, string(badge.DefaultSortBy)},
		{"limit default", func() interface{} { return fs.Lookup("badge-limit").DefValue }, "5"}, // Default as string
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.gotFunc()
			if got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestBadgeConfigRegisterFlags(t *testing.T) {
	bc := newBadgeConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	bc.registerBadgeFlags(fs)

	// Test that flags were registered
	tests := []struct {
		name         string
		flagName     string
		wantNotNil   bool
	}{
		{"badge-style flag", "badge-style", true},
		{"badge-variant flag", "badge-variant", true},
		{"badge-theme flag", "badge-theme", true},
		{"badge-output flag", "badge-output", true},
		{"badge-sort flag", "badge-sort", true},
		{"badge-limit flag", "badge-limit", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fs.Lookup(tt.flagName)
			if (f != nil) != tt.wantNotNil {
				t.Errorf("flag %s registered = %v, want %v", tt.flagName, f != nil, tt.wantNotNil)
			}
		})
	}
}

func TestBadgeConfigFlagValues(t *testing.T) {
	bc := newBadgeConfig()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	bc.registerBadgeFlags(fs)

	// Parse some test flags
	args := []string{
		"--badge-style", "compact",
		"--badge-variant", "text-based",
		"--badge-theme", "light",
		"--badge-output", "test-badge.svg",
		"--badge-sort", "stars",
		"--badge-limit", "10",
	}

	if err := fs.Parse(args); err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Verify flags were parsed correctly by checking flag values
	tests := []struct {
		name    string
		flagName string
		want    string
	}{
		{"style", "badge-style", "compact"},
		{"variant", "badge-variant", "text-based"},
		{"theme", "badge-theme", "light"},
		{"output", "badge-output", "test-badge.svg"},
		{"sort", "badge-sort", "stars"},
		{"limit", "badge-limit", "10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fs.Lookup(tt.flagName)
			if f == nil {
				t.Fatalf("flag %s not found", tt.flagName)
			}
			got := f.Value.String()
			if got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestBadgeConfigMultipleFlagSets(t *testing.T) {
	// Test that the same BadgeConfig can be used with different flag sets
	bc := newBadgeConfig()

	fs1 := flag.NewFlagSet("test1", flag.ContinueOnError)
	bc.registerBadgeFlags(fs1)

	fs2 := flag.NewFlagSet("test2", flag.ContinueOnError)
	bc.registerBadgeFlags(fs2)

	// Parse different values in different flag sets
	if err := fs1.Parse([]string{"--badge-style", "summary"}); err != nil {
		t.Fatalf("Failed to parse fs1: %v", err)
	}

	if bc.style != "summary" {
		t.Errorf("After fs1.Parse, style = %s, want summary", bc.style)
	}

	if err := fs2.Parse([]string{"--badge-style", "compact"}); err != nil {
		t.Fatalf("Failed to parse fs2: %v", err)
	}

	// The same BadgeConfig should reflect the last parsed value
	if bc.style != "compact" {
		t.Errorf("After fs2.Parse, style = %s, want compact", bc.style)
	}
}

func TestBadgeConfigDefaultValues(t *testing.T) {
	// Create multiple instances to ensure defaults are consistent
	bc1 := newBadgeConfig()
	bc2 := newBadgeConfig()
	bc3 := newBadgeConfig()

	configs := []*BadgeConfig{bc1, bc2, bc3}

	for i, bc := range configs {
		if bc == nil {
			t.Errorf("config %d is nil", i)
			continue
		}

		// Test defaults by checking registered flag default values
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		bc.registerBadgeFlags(fs)

		if fs.Lookup("badge-style").DefValue != string(badge.DefaultBadgeStyle) {
			t.Errorf("config %d: style default = %s, want %s", i, fs.Lookup("badge-style").DefValue, badge.DefaultBadgeStyle)
		}
		if fs.Lookup("badge-variant").DefValue != string(badge.DefaultBadgeVariant) {
			t.Errorf("config %d: variant default = %s, want %s", i, fs.Lookup("badge-variant").DefValue, badge.DefaultBadgeVariant)
		}
		if fs.Lookup("badge-theme").DefValue != string(badge.DefaultBadgeTheme) {
			t.Errorf("config %d: theme default = %s, want %s", i, fs.Lookup("badge-theme").DefValue, badge.DefaultBadgeTheme)
		}
		if fs.Lookup("badge-limit").DefValue != "5" {
			t.Errorf("config %d: limit default = %s, want 5", i, fs.Lookup("badge-limit").DefValue)
		}
	}
}

func TestBadgeConfigPointerReturn(t *testing.T) {
	// Verify that newBadgeConfig returns a pointer
	bc1 := newBadgeConfig()
	bc2 := newBadgeConfig()

	// Modifying bc1 should not affect bc2
	bc1.style = "modified"
	
	if bc2.style == "modified" {
		t.Error("Modifying bc1 affected bc2; they should be independent instances")
	}
}
