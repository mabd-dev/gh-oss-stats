package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	analytics "github.com/mabd-dev/gh-oss-stats/internal/analytic"
)

var (
	toolName          = "gh-oss-stats"
	telemetryFileName = "telemetry.json"
)

type Telemetry struct {
	NoticeShown bool `json:"noticeShown"`
}

func Send(version string) {
	sendTrackUsageEvent := func(isCI bool) error {
		analytics := analytics.CreateAnalytics()
		return analytics.TrackToolUsage(runtime.GOOS, version, isCI)
	}

	isCI := os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != ""
	if isCI {
		go sendTrackUsageEvent(true)
		return
	}

	telemetryDisabled := os.Getenv("GH_OSS_STATS_TELEMETRY_DISABLED")
	if telemetryDisabled == "1" {
		return
	}

	firstRun := IsFirstRun()

	if firstRun {
		PrintNotice()
		if err := MarkFirstRunDone(); err != nil {
			println(err)
		}
	}

	go sendTrackUsageEvent(false)
}

func PrintNotice() {
	fmt.Println("gh-oss-stats collects anonymous usage telemetry to help improve the tool.")
	fmt.Println("No personal data or GitHub credentials are collected.")
	fmt.Println("To disable: export GH_OSS_STATS_TELEMETRY_DISABLED=1")
	fmt.Println("More info: https://github.com/mabd-dev/gh-oss-stats#telemetry")
}

func IsFirstRun() bool {
	configDir, _ := os.UserConfigDir()
	markerPath := filepath.Join(configDir, toolName, telemetryFileName)
	_, err := os.Stat(markerPath)
	return os.IsNotExist(err)
}

func MarkFirstRunDone() error {
	telemetry := Telemetry{
		NoticeShown: true,
	}
	jsonData, err := json.Marshal(telemetry)
	if err != nil {
		return err
	}

	configDir, _ := os.UserConfigDir()
	dir := filepath.Join(configDir, toolName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(dir, telemetryFileName)

	return os.WriteFile(filePath, jsonData, 0644)
}
