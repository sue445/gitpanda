package webhook

import (
	"github.com/jarcoal/httpmock"
	"github.com/sue445/gitpanda/gitlab"
	"github.com/sue445/gitpanda/testutil"
	"testing"
)

func TestSlackWebhook_Request(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register stub
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site",
		httpmock.NewStringResponder(200, testutil.ReadTestData("../gitlab/testdata/project.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/merge_requests/1",
		httpmock.NewStringResponder(200, testutil.ReadTestData("../gitlab/testdata/merge_request.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/users?username=john_smith",
		httpmock.NewStringResponder(200, testutil.ReadTestData("../gitlab/testdata/users.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fmissing-repo",
		httpmock.NewStringResponder(404, testutil.ReadTestData("../gitlab/testdata/project_not_found.json")),
	)

	httpmock.RegisterResponder(
		"POST",
		"https://slack.com/api/chat.unfurl",
		httpmock.NewStringResponder(200, "{ \"ok\": true }"),
	)

	type args struct {
		body          string
		truncateLines int
	}

	s := NewSlackWebhook(
		"xoxp-0000000000-0000000000-000000000000-00000000000000000000000000000000",
		"",
		&gitlab.URLParserParams{
			APIEndpoint:    "http://example.com/api/v4",
			BaseURL:        "http://example.com",
			PrivateToken:   "xxxxxxxxxx",
			IsDebugLogging: true,
		},
	)
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Successful event_callback (valid only)",
			args: args{
				body:          testutil.ReadTestData("testdata/event_callback_link_shared.json"),
				truncateLines: 0,
			},
			want: "ok",
		},
		{
			name: "Successful event_callback (invalid only)",
			args: args{
				body:          testutil.ReadTestData("testdata/event_callback_link_shared_invalid.json"),
				truncateLines: 0,
			},
			want:    "Failed: FetchURL",
			wantErr: true,
		},
		{
			name: "Successful event_callback (valid and invalid)",
			args: args{
				body:          testutil.ReadTestData("testdata/event_callback_link_shared_valid_and_invalid.json"),
				truncateLines: 0,
			},
			want: "ok",
		},
		{
			name: "Successful url_verification",
			args: args{
				body:          testutil.ReadTestData("testdata/url_verification.json"),
				truncateLines: 0,
			},
			want: "3eZbrw1aBm2rZgRNFdxV2595E9CY3gmdALWMmHkvFXO7tYXAYM8P",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Request(tt.args.body, tt.args.truncateLines)
			if (err != nil) != tt.wantErr {
				t.Errorf("SlackWebhook.Request() error = %+v and got = %s, wantErr %+v", err, got, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SlackWebhook.Request() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
