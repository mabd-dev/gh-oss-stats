package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
	fmt.Printf("isCI=%v\n", isCI)

	if isCI {
		sendTrackUsageEvent(true)
		return
	}

	telemetryDisabled := os.Getenv("GH_OSS_STATS_TELEMETRY_DISABLED")
	if strings.ToLower(telemetryDisabled) == "true" {
		fmt.Println("telemetry disabled")
		return
	}

	firstRun := IsFirstRun()

	if firstRun {
		PrintNotice()
		if err := MarkFirstRunDone(); err != nil {
			println(err)
		}
	}

	sendTrackUsageEvent(false)
}

func PrintNotice() {
	println("We send analytics data!")
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
