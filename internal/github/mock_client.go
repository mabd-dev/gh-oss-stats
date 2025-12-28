package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// MockAPIClient is a mock GitHub API client that reads from local JSON files.
// Used for testing and development without hitting the real GitHub API.
type MockAPIClient struct {
	mockDataDir string
}

// NewMockAPIClient creates a new mock API client.
func NewMockAPIClient() *MockAPIClient {
	// Get the path to the mock data directory
	// Assuming we're running from the project root
	mockDataDir := filepath.Join("internal", "github", "mockResponses")

	return &MockAPIClient{
		mockDataDir: mockDataDir,
	}
}

// SearchIssues returns mock search results for merged PRs.
func (c *MockAPIClient) SearchIssues(ctx context.Context, query string, page, perPage int) (*SearchIssuesResponse, *http.Response, error) {
	var result SearchIssuesResponse

	if err := json.Unmarshal([]byte(mergedPrs), &result); err != nil {
		return nil, nil, fmt.Errorf("unmarshal pull request failed: %w", err)
	}

	// Create a mock response
	mockResp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
	}
	mockResp.Header.Set("X-RateLimit-Remaining", "5000")
	mockResp.Header.Set("X-RateLimit-Limit", "5000")

	return &result, mockResp, nil
}

// GetPullRequest returns mock PR details.
func (c *MockAPIClient) GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, *http.Response, error) {
	jsonData := `
{
  "number": 20,
  "state": "closed",
  "title": "Feature/share app",
  "user": {
    "login": "mabd-dev",
    "id": 133316956,
    "type": "User"
  },
  "created_at": "2025-12-12T11:05:16Z",
  "updated_at": "2025-12-17T06:42:03Z",
  "closed_at": "2025-12-17T05:14:39Z",
  "merged_at": "2025-12-17T05:14:39Z",
  "merged": true,
  "commits": 3,
  "additions": 150,
  "deletions": 20,
  "changed_files": 5,
  "html_url": "https://github.com/ibad-al-rahman/android-public/pull/20"
}
	`

	var result PullRequest
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		return nil, nil, fmt.Errorf("unmarshal pull request failed: %w", err)
	}

	// Create a mock response
	mockResp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
	}
	mockResp.Header.Set("X-RateLimit-Remaining", "5000")
	mockResp.Header.Set("X-RateLimit-Limit", "5000")

	return &result, mockResp, nil
}

// GetRepository returns mock repository information.
func (c *MockAPIClient) GetRepository(ctx context.Context, owner, repo string) (*Repository, *http.Response, error) {
	jsonData := `
{
  "name": "android-public",
  "full_name": "ibad-al-rahman/android-public",
  "owner": {
    "login": "ibad-al-rahman",
    "id": 123456,
    "type": "User"
  },
  "description": "Android app for Ibad Al-Rahman",
  "html_url": "https://github.com/ibad-al-rahman/android-public",
  "fork": false,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2025-12-17T06:42:03Z",
  "pushed_at": "2025-12-17T05:14:39Z",
  "stargazers_count": 15,
  "language": "Kotlin",
  "forks_count": 2,
  "open_issues_count": 3,
  "default_branch": "main"
}
	`
	var result Repository

	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		return nil, nil, fmt.Errorf("unmarshal pull request failed: %w", err)
	}

	// Create a mock response
	mockResp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
	}
	mockResp.Header.Set("X-RateLimit-Remaining", "5000")
	mockResp.Header.Set("X-RateLimit-Limit", "5000")

	return &result, mockResp, nil
}

// GetRateLimit returns mock rate limit information.
func (c *MockAPIClient) GetRateLimit(ctx context.Context) (*RateLimitResponse, error) {
	return &RateLimitResponse{
		Resources: RateLimitResources{
			Core: RateLimit{
				Limit:     5000,
				Remaining: 5000,
				Reset:     1735689600, // Some future timestamp
				Used:      0,
			},
			Search: RateLimit{
				Limit:     30,
				Remaining: 30,
				Reset:     1735689600,
				Used:      0,
			},
		},
	}, nil
}

var mergedPrs = `{
    "total_count": 17,
    "incomplete_results": false,
    "items": [
        {
            "url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/20",
            "repository_url": "https://api.github.com/repos/ibad-al-rahman/android-public",
            "labels_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/20/labels{/name}",
            "comments_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/20/comments",
            "events_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/20/events",
            "html_url": "https://github.com/ibad-al-rahman/android-public/pull/20",
            "id": 3723008628,
            "node_id": "PR_kwDONSGDVc64i7R5",
            "number": 20,
            "title": "Feature/share app",
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 0,
            "created_at": "2025-12-12T11:05:16Z",
            "updated_at": "2025-12-17T06:42:03Z",
            "closed_at": "2025-12-17T05:14:39Z",
            "author_association": "CONTRIBUTOR",
            "type": null,
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/ibad-al-rahman/android-public/pulls/20",
                "html_url": "https://github.com/ibad-al-rahman/android-public/pull/20",
                "diff_url": "https://github.com/ibad-al-rahman/android-public/pull/20.diff",
                "patch_url": "https://github.com/ibad-al-rahman/android-public/pull/20.patch",
                "merged_at": "2025-12-17T05:14:39Z"
            },
            "body": "#10 \r\n\r\n### What this feature does\r\nThis PR add a button in settings screen. When clicked open share menu with text containing android+iOS links in play+app stores\r\n\r\nLet me know if you want to change the share text message\r\n\r\n",
            "reactions": {
                "url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/20/reactions",
                "total_count": 0,
                "+1": 0,
                "-1": 0,
                "laugh": 0,
                "hooray": 0,
                "confused": 0,
                "heart": 0,
                "rocket": 0,
                "eyes": 0
            },
            "timeline_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/20/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/19",
            "repository_url": "https://api.github.com/repos/ibad-al-rahman/android-public",
            "labels_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/19/labels{/name}",
            "comments_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/19/comments",
            "events_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/19/events",
            "html_url": "https://github.com/ibad-al-rahman/android-public/pull/19",
            "id": 3672911025,
            "node_id": "PR_kwDONSGDVc617qVf",
            "number": 19,
            "title": "Fix vertical padding in weekly view",
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 0,
            "created_at": "2025-11-27T20:20:53Z",
            "updated_at": "2025-11-28T21:03:58Z",
            "closed_at": "2025-11-28T21:03:58Z",
            "author_association": "CONTRIBUTOR",
            "type": null,
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/ibad-al-rahman/android-public/pulls/19",
                "html_url": "https://github.com/ibad-al-rahman/android-public/pull/19",
                "diff_url": "https://github.com/ibad-al-rahman/android-public/pull/19.diff",
                "patch_url": "https://github.com/ibad-al-rahman/android-public/pull/19.patch",
                "merged_at": "2025-11-28T21:03:58Z"
            },
            "body": "Vertical padding in 'weekly' view is inconsistent with 'daily' view\r\n\r\nIn 'weekly':\r\n1. bottom padding is 0\r\n2. top padding is parent of parent column modifier, which causes a cut look when scrolling",
            "timeline_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/19/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/16",
            "repository_url": "https://api.github.com/repos/ibad-al-rahman/android-public",
            "labels_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/16/labels{/name}",
            "comments_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/16/comments",
            "events_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/16/events",
            "html_url": "https://github.com/ibad-al-rahman/android-public/pull/16",
            "id": 3653985904,
            "node_id": "PR_kwDONSGDVc608xz3",
            "number": 16,
            "title": "Fix sharedPerfs key when saving digest",
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 0,
            "created_at": "2025-11-22T06:03:23Z",
            "updated_at": "2025-11-27T20:12:07Z",
            "closed_at": "2025-11-22T06:31:05Z",
            "author_association": "CONTRIBUTOR",
            "type": null,
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/ibad-al-rahman/android-public/pulls/16",
                "html_url": "https://github.com/ibad-al-rahman/android-public/pull/16",
                "diff_url": "https://github.com/ibad-al-rahman/android-public/pull/16.diff",
                "patch_url": "https://github.com/ibad-al-rahman/android-public/pull/16.patch",
                "merged_at": "2025-11-22T06:31:05Z"
            },
            "body": null,
            "timeline_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/16/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/17",
            "repository_url": "https://api.github.com/repos/ibad-al-rahman/android-public",
            "labels_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/17/labels{/name}",
            "comments_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/17/comments",
            "events_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/17/events",
            "html_url": "https://github.com/ibad-al-rahman/android-public/pull/17",
            "id": 3655531202,
            "node_id": "PR_kwDONSGDVc61Br_N",
            "number": 17,
            "title": "Fix/settings screen",
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 0,
            "created_at": "2025-11-23T05:39:41Z",
            "updated_at": "2025-11-27T20:12:06Z",
            "closed_at": "2025-11-24T10:55:07Z",
            "author_association": "CONTRIBUTOR",
            "type": null,
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/ibad-al-rahman/android-public/pulls/17",
                "html_url": "https://github.com/ibad-al-rahman/android-public/pull/17",
                "diff_url": "https://github.com/ibad-al-rahman/android-public/pull/17.diff",
                "patch_url": "https://github.com/ibad-al-rahman/android-public/pull/17.patch",
                "merged_at": "2025-11-24T10:55:07Z"
            },
            "body": "### Issues\r\nThis PR fixes 3 things:\r\n- Settings top container was 'Column': for small screen this might be a problem if screen height is small, some elements won't be accessible. When this screen content increases, this will definitely become a problem\r\n- 'clickable modifier was placed before 'padding' which caused ripple effect (when user long press) to have a padding, which does not look nice\r\n- 'Clear Cache' button height was a tiny bit bigger than all other buttons\r\n\r\n### Solutions\r\n- Used 'LazyColumn' as container for the screen\r\n- Used 'ListItem' component to ensure same height and consistent look to all buttons\r\n\r\n\r\nBefore             |  After\r\n:-------------------------:|:-------------------------:\r\n![Screenshot_20251123_073430_Ibad Al-Rahman](https://github.com/user-attachments/assets/7fcb2f7e-55e4-4d03-9167-6db5df248550)  |  ![Screenshot_20251123_073514_Ibad Al-Rahman](https://github.com/user-attachments/assets/2a60f60f-a09e-4c1f-aa6b-74e81683024b)\r\n\r\n\r\n\r\n\r\n\r\n",
            "timeline_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/17/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/18",
            "repository_url": "https://api.github.com/repos/ibad-al-rahman/android-public",
            "labels_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/18/labels{/name}",
            "comments_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/18/comments",
            "events_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/18/events",
            "html_url": "https://github.com/ibad-al-rahman/android-public/pull/18",
            "id": 3655830448,
            "node_id": "PR_kwDONSGDVc61Cp7H",
            "number": 18,
            "title": "Theme Switcher",
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 1,
            "created_at": "2025-11-23T10:56:04Z",
            "updated_at": "2025-11-27T20:12:05Z",
            "closed_at": "2025-11-25T14:08:48Z",
            "author_association": "CONTRIBUTOR",
            "type": null,
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/ibad-al-rahman/android-public/pulls/18",
                "html_url": "https://github.com/ibad-al-rahman/android-public/pull/18",
                "diff_url": "https://github.com/ibad-al-rahman/android-public/pull/18.diff",
                "patch_url": "https://github.com/ibad-al-rahman/android-public/pull/18.patch",
                "merged_at": "2025-11-25T14:08:48Z"
            },
            "body": "- Created 'Settings-Repository' module \r\n- Created 'SettingsRepository' to manipulates settings data. For now it only has theme functionalities. We can add more as the app grows\r\n- Use new 'SettingsRepository' in 'Settings' screen\r\n- Create 'Theme' enum to hold all options, I think this is better than passing int's around the app. Enum makes it more readable and predictable\r\n- Added translation to theme options ",
            "timeline_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/18/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/nsh07/Tomato/issues/148",
            "repository_url": "https://api.github.com/repos/nsh07/Tomato",
            "labels_url": "https://api.github.com/repos/nsh07/Tomato/issues/148/labels{/name}",
            "comments_url": "https://api.github.com/repos/nsh07/Tomato/issues/148/comments",
            "events_url": "https://api.github.com/repos/nsh07/Tomato/issues/148/events",
            "html_url": "https://github.com/nsh07/Tomato/pull/148",
            "id": 3631969531,
            "node_id": "PR_kwDOPEHpd86zyTfP",
            "number": 148,
            "title": "Fix/info button",
            "labels": [],
            "state": "closed",
            "locked": false,
            "comments": 6,
            "created_at": "2025-11-17T06:44:50Z",
            "updated_at": "2025-11-22T17:13:01Z",
            "closed_at": "2025-11-19T12:06:16Z",
            "author_association": "CONTRIBUTOR",
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/nsh07/Tomato/pulls/148",
                "html_url": "https://github.com/nsh07/Tomato/pull/148",
                "diff_url": "https://github.com/nsh07/Tomato/pull/148.diff",
                "patch_url": "https://github.com/nsh07/Tomato/pull/148.patch",
                "merged_at": "2025-11-19T12:06:16Z"
            },
            "body": "fix for #146 issue\r\n\r\nApply 'Scaffold' content paddings to all 'Settings' screens",
            "timeline_url": "https://api.github.com/repos/nsh07/Tomato/issues/148/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/15",
            "repository_url": "https://api.github.com/repos/ibad-al-rahman/android-public",
            "labels_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/15/labels{/name}",
            "comments_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/15/comments",
            "events_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/15/events",
            "html_url": "https://github.com/ibad-al-rahman/android-public/pull/15",
            "id": 3650325835,
            "node_id": "PR_kwDONSGDVc60wV8I",
            "number": 15,
            "title": "Generate graph to visualize modules connections ",
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 0,
            "created_at": "2025-11-21T05:56:03Z",
            "updated_at": "2025-11-22T04:46:25Z",
            "closed_at": "2025-11-21T14:48:30Z",
            "author_association": "CONTRIBUTOR",
            "type": null,
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/ibad-al-rahman/android-public/pulls/15",
                "html_url": "https://github.com/ibad-al-rahman/android-public/pull/15",
                "diff_url": "https://github.com/ibad-al-rahman/android-public/pull/15.diff",
                "patch_url": "https://github.com/ibad-al-rahman/android-public/pull/15.patch",
                "merged_at": "2025-11-21T14:48:30Z"
            },
            "body": null,
            "timeline_url": "https://api.github.com/repos/ibad-al-rahman/android-public/issues/15/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/nsh07/Tomato/issues/138",
            "repository_url": "https://api.github.com/repos/nsh07/Tomato",
            "labels_url": "https://api.github.com/repos/nsh07/Tomato/issues/138/labels{/name}",
            "comments_url": "https://api.github.com/repos/nsh07/Tomato/issues/138/comments",
            "events_url": "https://api.github.com/repos/nsh07/Tomato/issues/138/events",
            "html_url": "https://github.com/nsh07/Tomato/pull/138",
            "id": 3606345525,
            "node_id": "PR_kwDOPEHpd86ycikr",
            "number": 138,
            "title": "Fix/align focus text",
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 11,
            "created_at": "2025-11-10T06:31:32Z",
            "updated_at": "2025-11-21T05:45:51Z",
            "closed_at": "2025-11-21T05:45:34Z",
            "author_association": "CONTRIBUTOR",
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/nsh07/Tomato/pulls/138",
                "html_url": "https://github.com/nsh07/Tomato/pull/138",
                "diff_url": "https://github.com/nsh07/Tomato/pull/138.diff",
                "patch_url": "https://github.com/nsh07/Tomato/pull/138.patch",
                "merged_at": "2025-11-21T05:45:34Z"
            },
            "body": "Based on #123 \r\n\r\nApparently 'TopAppBar' have slight padding on the right. 'Alignment.Center' is enough to center the text, no need for 'fillMaxWidth()'",
            "timeline_url": "https://api.github.com/repos/nsh07/Tomato/issues/138/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/issues/9",
            "repository_url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker",
            "labels_url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/issues/9/labels{/name}",
            "comments_url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/issues/9/comments",
            "events_url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/issues/9/events",
            "html_url": "https://github.com/qamarelsafadi/JetpackComposeTracker/pull/9",
            "id": 3204496021,
            "node_id": "PR_kwDONQBujs6diLmP",
            "number": 9,
            "title": "🔧 Refactor: Add Global Theme Support for UI Customization",
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 2,
            "created_at": "2025-07-05T07:18:08Z",
            "updated_at": "2025-07-22T14:26:12Z",
            "closed_at": "2025-07-21T21:39:53Z",
            "author_association": "CONTRIBUTOR",
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/pulls/9",
                "html_url": "https://github.com/qamarelsafadi/JetpackComposeTracker/pull/9",
                "diff_url": "https://github.com/qamarelsafadi/JetpackComposeTracker/pull/9.diff",
                "patch_url": "https://github.com/qamarelsafadi/JetpackComposeTracker/pull/9.patch",
                "merged_at": "2025-07-21T21:39:53Z"
            },
            "body": "### 🔧 Refactor: Add Global Theme Support for UI Customization\r\n✨ Summary\r\n\r\nThis PR refactors the trackRecompositions modifier to support global UI customization via a new RecompositionTrackerTheme object. Previously, all styles (border, text style, padding) were hardcoded, limiting flexibility. With this change, developers can now fully control the appearance of the recomposition tracker across their app.\r\n\r\n\r\n✅ What's Changed\r\n- Introduced RecompositionTrackerTheme as a singleton data object\r\n- Added individual configurable properties:\r\n  - enabled: Global switch to enable/disable tracking\r\n  - border: Controls border width, color, and shape\r\n  - textStyle: Customizable prefix, text style, and position (Offset)\r\n  - contentPadding: Padding applied to the Composable content\r\n- Removed hardcoded UI values from the trackRecompositions modifier\r\n- Updated documentation with a new \"UI Customization\" section for the README\r\n\r\n\r\n\r\n### 🧠 Motivation\r\n\r\nThis change allows users to:\r\n- Apply consistent styling across all tracked Composables\r\n- Adjust the visual appearance to match their design system\r\n- Enable or disable tracking globally with a single flag\r\n\r\n\r\nbased on #6 ",
            "timeline_url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/issues/9/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/android/cahier/issues/5",
            "repository_url": "https://api.github.com/repos/android/cahier",
            "labels_url": "https://api.github.com/repos/android/cahier/issues/5/labels{/name}",
            "comments_url": "https://api.github.com/repos/android/cahier/issues/5/comments",
            "events_url": "https://api.github.com/repos/android/cahier/issues/5/events",
            "html_url": "https://github.com/android/cahier/pull/5",
            "id": 3101526208,
            "node_id": "PR_kwDOOsxFSc6YL71W",
            "number": 5,
            "title": "show image and note in NoteItem",
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 2,
            "created_at": "2025-05-29T21:25:38Z",
            "updated_at": "2025-07-12T02:39:22Z",
            "closed_at": "2025-07-11T12:52:46Z",
            "author_association": "CONTRIBUTOR",
            "type": null,
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/android/cahier/pulls/5",
                "html_url": "https://github.com/android/cahier/pull/5",
                "diff_url": "https://github.com/android/cahier/pull/5.diff",
                "patch_url": "https://github.com/android/cahier/pull/5.patch",
                "merged_at": "2025-07-11T12:52:46Z"
            },
            "body": "Same database, but previously ui was not handled property to show note content (images and text)\r\n\r\n| Before | After |\r\n| --- | --- |\r\n| <img src=\"https://github.com/user-attachments/assets/7900b81c-50c8-4d3c-b40d-9e6ba643c76a\" height=\"480\" /> | <img src=\"https://github.com/user-attachments/assets/d71c44b1-e171-43aa-9e14-4794dbf550f6\" height=\"480\" /> | \r\n",
            "timeline_url": "https://api.github.com/repos/android/cahier/issues/5/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/issues/3",
            "repository_url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker",
            "labels_url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/issues/3/labels{/name}",
            "comments_url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/issues/3/comments",
            "events_url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/issues/3/events",
            "html_url": "https://github.com/qamarelsafadi/JetpackComposeTracker/pull/3",
            "id": 3144084558,
            "node_id": "PR_kwDONQBujs6aa06p",
            "number": 3,
            "title": "Suggest a way to safely keep trackRecompositions modifier",
            "user": {
                "login": "mabd-dev",
                "id": 133316956,
                "node_id": "U_kgDOB_JBXA",
                "avatar_url": "https://avatars.githubusercontent.com/u/133316956?v=4",
                "gravatar_id": "",
                "url": "https://api.github.com/users/mabd-dev",
                "html_url": "https://github.com/mabd-dev",
                "followers_url": "https://api.github.com/users/mabd-dev/followers",
                "following_url": "https://api.github.com/users/mabd-dev/following{/other_user}",
                "gists_url": "https://api.github.com/users/mabd-dev/gists{/gist_id}",
                "starred_url": "https://api.github.com/users/mabd-dev/starred{/owner}{/repo}",
                "subscriptions_url": "https://api.github.com/users/mabd-dev/subscriptions",
                "organizations_url": "https://api.github.com/users/mabd-dev/orgs",
                "repos_url": "https://api.github.com/users/mabd-dev/repos",
                "events_url": "https://api.github.com/users/mabd-dev/events{/privacy}",
                "received_events_url": "https://api.github.com/users/mabd-dev/received_events",
                "type": "User",
                "user_view_type": "public",
                "site_admin": false
            },
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 1,
            "created_at": "2025-06-13T16:19:01Z",
            "updated_at": "2025-06-14T20:55:24Z",
            "closed_at": "2025-06-14T20:55:24Z",
            "author_association": "CONTRIBUTOR",
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/pulls/3",
                "html_url": "https://github.com/qamarelsafadi/JetpackComposeTracker/pull/3",
                "diff_url": "https://github.com/qamarelsafadi/JetpackComposeTracker/pull/3.diff",
                "patch_url": "https://github.com/qamarelsafadi/JetpackComposeTracker/pull/3.patch",
                "merged_at": "2025-06-14T20:55:24Z"
            },
            "body": "### Why\r\nManually adding and removing trackRecompositions() during development is tedious and error-prone. This PR introduces a new extension function, trackRecompositionsIf(), which conditionally applies recomposition tracking based on a Boolean flag (e.g. BuildConfig.DEBUG, remote config, or any custom condition).\r\n\r\nThis approach allows developers to leave tracking code safely in production, while toggling it on/off via a single variable. It also makes it easier to test recomposition behavior in live environments without needing code changes.",
            "reactions": {
                "url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/issues/3/reactions",
                "total_count": 0,
                "+1": 0,
                "-1": 0,
                "laugh": 0,
                "hooray": 0,
                "confused": 0,
                "heart": 0,
                "rocket": 0,
                "eyes": 0
            },
            "timeline_url": "https://api.github.com/repos/qamarelsafadi/JetpackComposeTracker/issues/3/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/android/nav3-recipes/issues/10",
            "repository_url": "https://api.github.com/repos/android/nav3-recipes",
            "labels_url": "https://api.github.com/repos/android/nav3-recipes/issues/10/labels{/name}",
            "comments_url": "https://api.github.com/repos/android/nav3-recipes/issues/10/comments",
            "events_url": "https://api.github.com/repos/android/nav3-recipes/issues/10/events",
            "html_url": "https://github.com/android/nav3-recipes/pull/10",
            "id": 3089246340,
            "node_id": "PR_kwDOOpDy1s6Xia9t",
            "number": 10,
            "title": "Remove scaffold contentPadding consumption",
            "user": {
                "login": "mabd-dev",
                "id": 133316956,
                "node_id": "U_kgDOB_JBXA",
                "avatar_url": "https://avatars.githubusercontent.com/u/133316956?v=4",
                "gravatar_id": "",
                "url": "https://api.github.com/users/mabd-dev",
                "html_url": "https://github.com/mabd-dev",
                "followers_url": "https://api.github.com/users/mabd-dev/followers",
                "following_url": "https://api.github.com/users/mabd-dev/following{/other_user}",
                "gists_url": "https://api.github.com/users/mabd-dev/gists{/gist_id}",
                "starred_url": "https://api.github.com/users/mabd-dev/starred{/owner}{/repo}",
                "subscriptions_url": "https://api.github.com/users/mabd-dev/subscriptions",
                "organizations_url": "https://api.github.com/users/mabd-dev/orgs",
                "repos_url": "https://api.github.com/users/mabd-dev/repos",
                "events_url": "https://api.github.com/users/mabd-dev/events{/privacy}",
                "received_events_url": "https://api.github.com/users/mabd-dev/received_events",
                "type": "User",
                "user_view_type": "public",
                "site_admin": false
            },
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 1,
            "created_at": "2025-05-25T11:51:27Z",
            "updated_at": "2025-06-13T14:49:40Z",
            "closed_at": "2025-06-09T18:31:50Z",
            "author_association": "CONTRIBUTOR",
            "type": null,
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/android/nav3-recipes/pulls/10",
                "html_url": "https://github.com/android/nav3-recipes/pull/10",
                "diff_url": "https://github.com/android/nav3-recipes/pull/10.diff",
                "patch_url": "https://github.com/android/nav3-recipes/pull/10.patch",
                "merged_at": "2025-06-09T18:31:50Z"
            },
            "body": "| Before    | After |\r\n| -------- | ------- |\r\n| <img src=\"https://github.com/user-attachments/assets/011f5098-7e8a-45f6-8df7-6e5e0ab244fd\" alt=\"Sample Image\" height=\"480\">  | <img src=\"https://github.com/user-attachments/assets/d1fddca7-0ca1-4f0c-881a-748c6473b65e\" alt=\"Sample Image\" height=\"480\">    |\r\n\r\n\r\n#### The problem:\r\nsystem bar padding are not being applied\r\n\r\n#### Why:\r\nNavHost is consuming scaffold 'paddingValues' without applying them. \r\n\r\n\r\nContentRed, ContentGreeen etc… have 'safeDrawingPadding()' modifier so they would add system bar paddings but, their parent composable (NavHost) already used/consumed them\r\n",
            "reactions": {
                "url": "https://api.github.com/repos/android/nav3-recipes/issues/10/reactions",
                "total_count": 0,
                "+1": 0,
                "-1": 0,
                "laugh": 0,
                "hooray": 0,
                "confused": 0,
                "heart": 0,
                "rocket": 0,
                "eyes": 0
            },
            "timeline_url": "https://api.github.com/repos/android/nav3-recipes/issues/10/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/android/nav3-recipes/issues/11",
            "repository_url": "https://api.github.com/repos/android/nav3-recipes",
            "labels_url": "https://api.github.com/repos/android/nav3-recipes/issues/11/labels{/name}",
            "comments_url": "https://api.github.com/repos/android/nav3-recipes/issues/11/comments",
            "events_url": "https://api.github.com/repos/android/nav3-recipes/issues/11/events",
            "html_url": "https://github.com/android/nav3-recipes/pull/11",
            "id": 3089272777,
            "node_id": "PR_kwDOOpDy1s6XigHO",
            "number": 11,
            "title": "New navigation recipe picker activity",
            "user": {
                "login": "mabd-dev",
                "id": 133316956,
                "node_id": "U_kgDOB_JBXA",
                "avatar_url": "https://avatars.githubusercontent.com/u/133316956?v=4",
                "gravatar_id": "",
                "url": "https://api.github.com/users/mabd-dev",
                "html_url": "https://github.com/mabd-dev",
                "followers_url": "https://api.github.com/users/mabd-dev/followers",
                "following_url": "https://api.github.com/users/mabd-dev/following{/other_user}",
                "gists_url": "https://api.github.com/users/mabd-dev/gists{/gist_id}",
                "starred_url": "https://api.github.com/users/mabd-dev/starred{/owner}{/repo}",
                "subscriptions_url": "https://api.github.com/users/mabd-dev/subscriptions",
                "organizations_url": "https://api.github.com/users/mabd-dev/orgs",
                "repos_url": "https://api.github.com/users/mabd-dev/repos",
                "events_url": "https://api.github.com/users/mabd-dev/events{/privacy}",
                "received_events_url": "https://api.github.com/users/mabd-dev/received_events",
                "type": "User",
                "user_view_type": "public",
                "site_admin": false
            },
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 1,
            "created_at": "2025-05-25T12:37:32Z",
            "updated_at": "2025-06-13T14:49:39Z",
            "closed_at": "2025-06-09T18:24:56Z",
            "author_association": "CONTRIBUTOR",
            "type": null,
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/android/nav3-recipes/pulls/11",
                "html_url": "https://github.com/android/nav3-recipes/pull/11",
                "diff_url": "https://github.com/android/nav3-recipes/pull/11.diff",
                "patch_url": "https://github.com/android/nav3-recipes/pull/11.patch",
                "merged_at": "2025-06-09T18:24:56Z"
            },
            "body": "<img src=\"https://github.com/user-attachments/assets/1ba332f1-ccf3-474b-8a97-5dc490d39fb8\" height=\"480\" />\r\n\r\n### Why\r\nFaster to try out different recipes without having to re-run the app\r\n\r\n",
            "reactions": {
                "url": "https://api.github.com/repos/android/nav3-recipes/issues/11/reactions",
                "total_count": 2,
                "+1": 2,
                "-1": 0,
                "laugh": 0,
                "hooray": 0,
                "confused": 0,
                "heart": 0,
                "rocket": 0,
                "eyes": 0
            },
            "timeline_url": "https://api.github.com/repos/android/nav3-recipes/issues/11/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/android/cahier/issues/4",
            "repository_url": "https://api.github.com/repos/android/cahier",
            "labels_url": "https://api.github.com/repos/android/cahier/issues/4/labels{/name}",
            "comments_url": "https://api.github.com/repos/android/cahier/issues/4/comments",
            "events_url": "https://api.github.com/repos/android/cahier/issues/4/events",
            "html_url": "https://github.com/android/cahier/pull/4",
            "id": 3101454989,
            "node_id": "PR_kwDOOsxFSc6YLskN",
            "number": 4,
            "title": "Convert DrawingToolbox to LazyRow",
            "user": {
                "login": "mabd-dev",
                "id": 133316956,
                "node_id": "U_kgDOB_JBXA",
                "avatar_url": "https://avatars.githubusercontent.com/u/133316956?v=4",
                "gravatar_id": "",
                "url": "https://api.github.com/users/mabd-dev",
                "html_url": "https://github.com/mabd-dev",
                "followers_url": "https://api.github.com/users/mabd-dev/followers",
                "following_url": "https://api.github.com/users/mabd-dev/following{/other_user}",
                "gists_url": "https://api.github.com/users/mabd-dev/gists{/gist_id}",
                "starred_url": "https://api.github.com/users/mabd-dev/starred{/owner}{/repo}",
                "subscriptions_url": "https://api.github.com/users/mabd-dev/subscriptions",
                "organizations_url": "https://api.github.com/users/mabd-dev/orgs",
                "repos_url": "https://api.github.com/users/mabd-dev/repos",
                "events_url": "https://api.github.com/users/mabd-dev/events{/privacy}",
                "received_events_url": "https://api.github.com/users/mabd-dev/received_events",
                "type": "User",
                "user_view_type": "public",
                "site_admin": false
            },
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 0,
            "created_at": "2025-05-29T20:50:07Z",
            "updated_at": "2025-06-03T14:08:20Z",
            "closed_at": "2025-06-03T14:08:20Z",
            "author_association": "CONTRIBUTOR",
            "type": null,
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/android/cahier/pulls/4",
                "html_url": "https://github.com/android/cahier/pull/4",
                "diff_url": "https://github.com/android/cahier/pull/4.diff",
                "patch_url": "https://github.com/android/cahier/pull/4.patch",
                "merged_at": "2025-06-03T14:08:20Z"
            },
            "body": "To be able to see all buttons on portrait mode\r\n\r\n\r\n### Before\r\n<img src=\"https://github.com/user-attachments/assets/10e1711d-4ff4-48be-b688-c35eca73c5b7\" height=\"480\" />\r\n\r\n### After\r\nhttps://github.com/user-attachments/assets/5e84bfce-5f7d-4b00-8795-bb3e24b9e1e5\r\n\r\n\r\n\r\n",
            "reactions": {
                "url": "https://api.github.com/repos/android/cahier/issues/4/reactions",
                "total_count": 0,
                "+1": 0,
                "-1": 0,
                "laugh": 0,
                "hooray": 0,
                "confused": 0,
                "heart": 0,
                "rocket": 0,
                "eyes": 0
            },
            "timeline_url": "https://api.github.com/repos/android/cahier/issues/4/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/issues/6",
            "repository_url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number",
            "labels_url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/issues/6/labels{/name}",
            "comments_url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/issues/6/comments",
            "events_url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/issues/6/events",
            "html_url": "https://github.com/esatgozcu/Compose-Rolling-Number/pull/6",
            "id": 2934000014,
            "node_id": "PR_kwDOKBJ8r86PZDY3",
            "number": 6,
            "title": "Bumped up gradle, agp, kotlin and all dependencies to latest versions",
            "user": {
                "login": "mabd-dev",
                "id": 133316956,
                "node_id": "U_kgDOB_JBXA",
                "avatar_url": "https://avatars.githubusercontent.com/u/133316956?v=4",
                "gravatar_id": "",
                "url": "https://api.github.com/users/mabd-dev",
                "html_url": "https://github.com/mabd-dev",
                "followers_url": "https://api.github.com/users/mabd-dev/followers",
                "following_url": "https://api.github.com/users/mabd-dev/following{/other_user}",
                "gists_url": "https://api.github.com/users/mabd-dev/gists{/gist_id}",
                "starred_url": "https://api.github.com/users/mabd-dev/starred{/owner}{/repo}",
                "subscriptions_url": "https://api.github.com/users/mabd-dev/subscriptions",
                "organizations_url": "https://api.github.com/users/mabd-dev/orgs",
                "repos_url": "https://api.github.com/users/mabd-dev/repos",
                "events_url": "https://api.github.com/users/mabd-dev/events{/privacy}",
                "received_events_url": "https://api.github.com/users/mabd-dev/received_events",
                "type": "User",
                "user_view_type": "public",
                "site_admin": false
            },
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 0,
            "created_at": "2025-03-20T04:23:44Z",
            "updated_at": "2025-03-26T21:33:08Z",
            "closed_at": "2025-03-26T21:33:08Z",
            "author_association": "CONTRIBUTOR",
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/pulls/6",
                "html_url": "https://github.com/esatgozcu/Compose-Rolling-Number/pull/6",
                "diff_url": "https://github.com/esatgozcu/Compose-Rolling-Number/pull/6.diff",
                "patch_url": "https://github.com/esatgozcu/Compose-Rolling-Number/pull/6.patch",
                "merged_at": "2025-03-26T21:33:08Z"
            },
            "body": "Gradle bumped to 8.11.1\r\nAGP bumped to 8.7.3\r\nKotlin bumped to 2.1.10\r\nUsed compose plugin version=2.1.10\r\nUsed compose BOM version=2025.03.00\r\nBumped up other dependencies to latest versions\r\n\r\n\r\nNOTE: Not entirely sure about maven publications changes. I changed based on google [documentation ](https://developer.android.com/build/publish-library/upload-library), I don't have enough experience in it and not sure if i can generate outputs myself. \r\nI made sure all warning (caused by agp upgrade to 8.x) are fixed\r\n\r\n",
            "reactions": {
                "url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/issues/6/reactions",
                "total_count": 0,
                "+1": 0,
                "-1": 0,
                "laugh": 0,
                "hooray": 0,
                "confused": 0,
                "heart": 0,
                "rocket": 0,
                "eyes": 0
            },
            "timeline_url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/issues/6/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/issues/4",
            "repository_url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number",
            "labels_url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/issues/4/labels{/name}",
            "comments_url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/issues/4/comments",
            "events_url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/issues/4/events",
            "html_url": "https://github.com/esatgozcu/Compose-Rolling-Number/pull/4",
            "id": 2855304819,
            "node_id": "PR_kwDOKBJ8r86LVNHP",
            "number": 4,
            "title": "Delete RollingNumberVm, use it's parameters in RollingNumberView",
            "user": {
                "login": "mabd-dev",
                "id": 133316956,
                "node_id": "U_kgDOB_JBXA",
                "avatar_url": "https://avatars.githubusercontent.com/u/133316956?v=4",
                "gravatar_id": "",
                "url": "https://api.github.com/users/mabd-dev",
                "html_url": "https://github.com/mabd-dev",
                "followers_url": "https://api.github.com/users/mabd-dev/followers",
                "following_url": "https://api.github.com/users/mabd-dev/following{/other_user}",
                "gists_url": "https://api.github.com/users/mabd-dev/gists{/gist_id}",
                "starred_url": "https://api.github.com/users/mabd-dev/starred{/owner}{/repo}",
                "subscriptions_url": "https://api.github.com/users/mabd-dev/subscriptions",
                "organizations_url": "https://api.github.com/users/mabd-dev/orgs",
                "repos_url": "https://api.github.com/users/mabd-dev/repos",
                "events_url": "https://api.github.com/users/mabd-dev/events{/privacy}",
                "received_events_url": "https://api.github.com/users/mabd-dev/received_events",
                "type": "User",
                "user_view_type": "public",
                "site_admin": false
            },
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 0,
            "created_at": "2025-02-15T07:29:42Z",
            "updated_at": "2025-02-17T15:46:51Z",
            "closed_at": "2025-02-17T15:46:51Z",
            "author_association": "CONTRIBUTOR",
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/pulls/4",
                "html_url": "https://github.com/esatgozcu/Compose-Rolling-Number/pull/4",
                "diff_url": "https://github.com/esatgozcu/Compose-Rolling-Number/pull/4.diff",
                "patch_url": "https://github.com/esatgozcu/Compose-Rolling-Number/pull/4.patch",
                "merged_at": "2025-02-17T15:46:51Z"
            },
            "body": "based on this issue: [https://github.com/esatgozcu/Compose-Rolling-Number/issues/3](https://github.com/esatgozcu/Compose-Rolling-Number/issues/3)",
            "reactions": {
                "url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/issues/4/reactions",
                "total_count": 0,
                "+1": 0,
                "-1": 0,
                "laugh": 0,
                "hooray": 0,
                "confused": 0,
                "heart": 0,
                "rocket": 0,
                "eyes": 0
            },
            "timeline_url": "https://api.github.com/repos/esatgozcu/Compose-Rolling-Number/issues/4/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        },
        {
            "url": "https://api.github.com/repos/zuzmuz/nvimawscli/issues/12",
            "repository_url": "https://api.github.com/repos/zuzmuz/nvimawscli",
            "labels_url": "https://api.github.com/repos/zuzmuz/nvimawscli/issues/12/labels{/name}",
            "comments_url": "https://api.github.com/repos/zuzmuz/nvimawscli/issues/12/comments",
            "events_url": "https://api.github.com/repos/zuzmuz/nvimawscli/issues/12/events",
            "html_url": "https://github.com/zuzmuz/nvimawscli/pull/12",
            "id": 2279611190,
            "node_id": "PR_kwDOLySutM5uki-W",
            "number": 12,
            "title": "Feature/target groups",
            "labels": [],
            "state": "closed",
            "locked": false,
            "assignee": null,
            "assignees": [],
            "milestone": null,
            "comments": 0,
            "created_at": "2024-05-05T16:27:40Z",
            "updated_at": "2024-05-06T20:37:14Z",
            "closed_at": "2024-05-06T20:37:13Z",
            "author_association": "CONTRIBUTOR",
            "active_lock_reason": null,
            "draft": false,
            "pull_request": {
                "url": "https://api.github.com/repos/zuzmuz/nvimawscli/pulls/12",
                "html_url": "https://github.com/zuzmuz/nvimawscli/pull/12",
                "diff_url": "https://github.com/zuzmuz/nvimawscli/pull/12.diff",
                "patch_url": "https://github.com/zuzmuz/nvimawscli/pull/12.patch",
                "merged_at": "2024-05-06T20:37:13Z"
            },
            "body": " ### What is new\r\n- Fetch target groups and display on new buffer\r\n- Sort functionality to target groups table columns\r\n\r\n\r\n### Fixes\r\n- 'tables.lua' rendering function fails on line:88 when concatenating boolean\r\n\r\n### Notes\r\n- When table has boolean, vertical lines are not placed correctly. To replicate add 'HealthCheckEnabled' to 'preferred_target_groups_attributes' in config.lua",
            "timeline_url": "https://api.github.com/repos/zuzmuz/nvimawscli/issues/12/timeline",
            "performed_via_github_app": null,
            "state_reason": null,
            "score": 1.0
        }
    ]
}`
