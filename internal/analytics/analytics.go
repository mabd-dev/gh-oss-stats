package analytics

import (
	"context"

	"github.com/mixpanel/mixpanel-go"
)

var mixpanelToken = ""

type Analytics struct {
	userUUID string
	Client   *mixpanel.ApiClient
}

func CreateAnalytics(userUUID string) Analytics {
	mixpanelClient := mixpanel.NewApiClient(mixpanelToken)
	return Analytics{
		userUUID: userUUID,
		Client:   mixpanelClient,
	}
}

func (analytics Analytics) Track(name string, params map[string]any) error {
	ctx := context.Background()
	return analytics.Client.Track(ctx, []*mixpanel.Event{
		analytics.Client.NewEvent(name, analytics.userUUID, params),
	})
}

func (analytics Analytics) TrackToolUsage(
	os string,
	version string,
	ci bool,
	command string,
) error {
	params := map[string]any{
		"os":           os,
		"tool-version": version,
		"ci":           ci,
		"project":      "gh-oss-stats",
		"command":      command,
	}
	return analytics.Track("usage", params)
}
