package analytics

import (
	"context"
	"os"

	"github.com/mixpanel/mixpanel-go"
)

type Analytics struct {
	Client *mixpanel.ApiClient
}

func CreateAnalytics() Analytics {
	mixpanelClient := mixpanel.NewApiClient(os.Getenv("MIXPANEL_PROJECT_TOKEN"))
	return Analytics{Client: mixpanelClient}
}

func (analytics Analytics) Track(name string, params map[string]any) error {
	ctx := context.Background()
	return analytics.Client.Track(ctx, []*mixpanel.Event{
		analytics.Client.NewEvent(name, "", params),
	})
}

func (analytics Analytics) TrackToolUsage(os string, version string, ci bool) error {
	params := map[string]any{
		"os":           os,
		"tool-version": version,
		"ci":           ci,
		"project":      "gh-oss-stats",
	}
	return analytics.Track("usage", params)
}
