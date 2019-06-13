package main

import (
	"io/ioutil"
	"path"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func readTestData(filename string) string {
	buf, err := ioutil.ReadFile(path.Join("test", filename))

	if err != nil {
		panic(err)
	}

	return string(buf)
}

func TestGitlabUrlParser_FetchURL(t *testing.T) {
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
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/issues/1",
		httpmock.NewStringResponder(200, readTestData("gitlab/issue.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/issues/1/notes/302",
		httpmock.NewStringResponder(200, readTestData("gitlab/issue_note.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/merge_requests/1",
		httpmock.NewStringResponder(200, readTestData("gitlab/merge_request.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/merge_requests/1/notes/301",
		httpmock.NewStringResponder(200, readTestData("gitlab/merge_request_note.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/users?username=john_smith",
		httpmock.NewStringResponder(200, readTestData("gitlab/users.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/repository/files/gitlabci-templates%2Fcontinuous_bundle_update.yml/raw?ref=master",
		httpmock.NewStringResponder(200, readTestData("gitlab/gitlabci-templates/continuous_bundle_update.yml")),
	)

	p, err := NewGitlabURLParser(&GitLabURLParserParams{
		APIEndpoint:  "http://example.com/api/v4",
		BaseURL:      "http://example.com",
		PrivateToken: "xxxxxxxxxx",
	})

	if err != nil {
		panic(err)
	}

	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    *GitLabPage
		wantErr bool
	}{
		{
			name: "Unknown URL",
			args: args{
				url: "http://foo.com/",
			},
			want: nil,
		},
		{
			name: "GitLab URL",
			args: args{
				url: "http://example.com/",
			},
			want: nil,
		},
		{
			name: "Project URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site",
			},
			want: &GitLabPage{
				Title:           "Diaspora / Diaspora Project Site · GitLab",
				Description:     "Diaspora Project",
				AuthorName:      "Diaspora",
				AuthorAvatarURL: "http://example.com/images/diaspora.png",
				AvatarURL:       "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				canTruncateDescription: true,
			},
		},
		{
			name: "Issue URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/issues/1",
			},
			want: &GitLabPage{
				Title:           "Ut commodi ullam eos dolores perferendis nihil sunt. · Issues · Diaspora / Diaspora Project Site · GitLab",
				Description:     "Omnis vero earum sunt corporis dolor et placeat.",
				AuthorName:      "Administrator",
				AuthorAvatarURL: "https://gitlab.example.com/images/root.png",
				AvatarURL:       "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				canTruncateDescription: true,
			},
		},
		{
			name: "Issue comment URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/issues/1#note_302",
			},
			want: &GitLabPage{
				Title:           "Ut commodi ullam eos dolores perferendis nihil sunt. · Issues · Diaspora / Diaspora Project Site · GitLab",
				Description:     "closed",
				AuthorName:      "Pip",
				AuthorAvatarURL: "http://localhost:3000/uploads/user/avatar/1/pipin.jpeg",
				AvatarURL:       "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				canTruncateDescription: true,
			},
		},
		{
			name: "MergeRequest URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/merge_requests/1",
			},
			want: &GitLabPage{
				Title:           "test1 · Merge Requests · Diaspora / Diaspora Project Site · GitLab",
				Description:     "fixed login page css paddings",
				AuthorName:      "Administrator",
				AuthorAvatarURL: "https://gitlab.example.com/images/admin.png",
				AvatarURL:       "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				canTruncateDescription: true,
			},
		},
		{
			name: "MergeRequest comment URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/merge_requests/1#note_301",
			},
			want: &GitLabPage{
				Title:           "test1 · Merge Requests · Diaspora / Diaspora Project Site · GitLab",
				Description:     "Comment for MR",
				AuthorName:      "Pip",
				AuthorAvatarURL: "http://localhost:3000/uploads/user/avatar/1/pipin.jpeg",
				AvatarURL:       "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				canTruncateDescription: true,
			},
		},
		{
			name: "User URL",
			args: args{
				url: "http://example.com/john_smith",
			},
			want: &GitLabPage{
				Title:           "John Smith · GitLab",
				Description:     "John Smith",
				AuthorName:      "John Smith",
				AuthorAvatarURL: "http://localhost:3000/uploads/user/avatar/1/cd8.jpeg",
				AvatarURL:       "http://localhost:3000/uploads/user/avatar/1/cd8.jpeg",
				canTruncateDescription: true,
			},
		},
		{
			name: "Blob URL (single line)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/blob/master/gitlabci-templates/continuous_bundle_update.yml#L4",
			},
			want: &GitLabPage{
				Title:           "gitlabci-templates/continuous_bundle_update.yml:4",
				Description:     "```\n  variables:\n```",
				AuthorName:      "",
				AuthorAvatarURL: "",
				AvatarURL:       "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				canTruncateDescription: false,
			},
		},
		{
			name: "Blob URL (multiple line)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/blob/master/gitlabci-templates/continuous_bundle_update.yml#L4-9",
			},
			want: &GitLabPage{
				Title:           "gitlabci-templates/continuous_bundle_update.yml:4-9",
				Description:     "```\n  variables:\n    CACHE_VERSION: \"v1\"\n    GIT_EMAIL:     \"gitlabci@example.com\"\n    GIT_USER:      \"GitLab CI\"\n    LABELS:        \"bundle update\"\n    OPTIONS:       \"\"\n```",
				AuthorName:      "",
				AuthorAvatarURL: "",
				AvatarURL:       "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				canTruncateDescription: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.FetchURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitlabURLParser.FetchURL() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitlabURLParser.FetchURL() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
