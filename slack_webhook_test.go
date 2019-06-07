package main

import (
	"github.com/jarcoal/httpmock"
	"testing"
)

func TestSlackWebhook_Request(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register stub
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site",
		httpmock.NewStringResponder(200, readTestData("gitlab/project.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/merge_requests/1",
		httpmock.NewStringResponder(200, readTestData("gitlab/merge_request.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/users?username=john_smith",
		httpmock.NewStringResponder(200, readTestData("gitlab/users.json")),
	)

	httpmock.RegisterResponder(
		"POST",
		"https://slack.com/api/chat.unfurl",
		httpmock.NewStringResponder(200, "{ \"ok\": true }"),
	)

	type args struct {
		body        string
		verifyToken bool
	}

	s := NewSlackWebhook(
		"xxxxxxx",
		&GitLabURLParserParams{
			APIEndpoint:  "http://example.com/api/v4",
			BaseURL:      "http://example.com",
			PrivateToken: "xxxxxxxxxx",
		},
	)
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Successful",
			args: args{
				body:        readTestData("slack/event_callback_link_shared.json"),
				verifyToken: false,
			},
			want: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Request(tt.args.body, tt.args.verifyToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("SlackWebhook.Request() error = %v and got = %s, wantErr %v", err, got, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SlackWebhook.Request() = %v, want %v", got, tt.want)
			}
		})
	}
}
