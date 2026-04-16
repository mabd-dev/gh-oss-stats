package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"

	analytics "github.com/mabd-dev/gh-oss-stats/internal/analytic"
	"github.com/mabd-dev/gh-oss-stats/internal/utils"
)

var (
	toolName          = "gh-oss-stats"
	telemetryFileName = "telemetry.json"
)

var mixpanelToken = "dev-token"

type Telemetry struct {
	NoticeShown bool   `json:"noticeShown"`
	UserUUID    string `json:"userUUID"`
}

func Send(version string) {
	telemetry, err := readOrCreateTelemetry()
	if err != nil {
		return
	}

	telemetryDisabled := os.Getenv("GH_OSS_STATS_TELEMETRY_DISABLED")
	if telemetryDisabled == "1" {
		return
	}

	isCI := os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != ""
	if isCI {
		sendTrackUsageEvent(telemetry.UserUUID, version, true)
		return
	}

	if !telemetry.NoticeShown {
		printNotice()

		telemetry.NoticeShown = true
		storeTelemetry(*telemetry)
	}

	sendTrackUsageEvent(telemetry.UserUUID, version, false)
}

func readOrCreateTelemetry() (*Telemetry, error) {
	t, err := readTelemetry()
	if err != nil {
		return nil, err
	}

	if t != nil {
		return t, nil
	}

	// Create + save new telemetry file

	userUUID := uuid.New().String()
	telemetry := Telemetry{
		NoticeShown: false,
		UserUUID:    userUUID,
	}
	if err := storeTelemetry(telemetry); err != nil {
		return nil, err
	}
	return &telemetry, nil
}

func readTelemetry() (*Telemetry, error) {
	configDir, _ := os.UserConfigDir()
	telementryPath := filepath.Join(configDir, toolName, telemetryFileName)

	exists, err := utils.FileExists(telementryPath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	data, err := os.ReadFile(telementryPath)
	if err != nil {
		return nil, err
	}

	var telemetry Telemetry
	if err := json.Unmarshal(data, &telemetry); err != nil {
		return nil, err
	}

	return &telemetry, nil
}

func storeTelemetry(t Telemetry) error {
	jsonData, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return err
	}

	configDir, _ := os.UserConfigDir()
	dir := filepath.Join(configDir, toolName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	filePath := filepath.Join(dir, telemetryFileName)

	return utils.WriteToFile(jsonData, filePath)
}

func printNotice() {
	fmt.Println("gh-oss-stats collects anonymous usage telemetry to help improve the tool.")
	fmt.Println("No personal data or GitHub credentials are collected.")
	fmt.Println("To disable: export GH_OSS_STATS_TELEMETRY_DISABLED=1")
	fmt.Println("More info: https://github.com/mabd-dev/gh-oss-stats#telemetry")
}

func sendTrackUsageEvent(userUUID string, version string, isCI bool) error {
	analytics := analytics.CreateAnalytics(userUUID)
	return analytics.TrackToolUsage(runtime.GOOS, version, isCI)
}
