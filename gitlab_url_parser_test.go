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

func TestGitlabUrlParser_FetchUrl(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register stub
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site",
		httpmock.NewStringResponder(200, readTestData("project.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/issues/1",
		httpmock.NewStringResponder(200, readTestData("issue.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/merge_requests/1",
		httpmock.NewStringResponder(200, readTestData("merge_request.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/users?username=john_smith",
		httpmock.NewStringResponder(200, readTestData("users.json")),
	)

	p := NewGitlabUrlParser(GitLabUrlParserParams{
		ApiEndpoint:  "http://example.com/api/v4",
		BaseUrl:      "http://example.com",
		PrivateToken: "xxxxxxxxxx",
	})
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
				Title:           "Diaspora / Diaspora Project Site",
				Description:     "Diaspora Project",
				AuthorName:      "Diaspora",
				AuthorAvatarURL: "http://example.com/images/diaspora.png",
				AvatarURL:       "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
			},
		},
		{
			name: "Issue URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/issues/1",
			},
			want: &GitLabPage{
				Title:           "Ut commodi ullam eos dolores perferendis nihil sunt.",
				Description:     "Omnis vero earum sunt corporis dolor et placeat.",
				AuthorName:      "Administrator",
				AuthorAvatarURL: "https://gitlab.example.com/images/root.png",
				AvatarURL:       "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
			},
		},
		{
			name: "MergeRequest URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/merge_requests/1",
			},
			want: &GitLabPage{
				Title:           "test1",
				Description:     "fixed login page css paddings",
				AuthorName:      "Administrator",
				AuthorAvatarURL: "https://gitlab.example.com/images/admin.png",
				AvatarURL:       "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
			},
		},
		{
			name: "User URL",
			args: args{
				url: "http://example.com/john_smith",
			},
			want: &GitLabPage{
				Title:           "John Smith",
				Description:     "John Smith",
				AuthorName:      "John Smith",
				AuthorAvatarURL: "http://localhost:3000/uploads/user/avatar/1/cd8.jpeg",
				AvatarURL:       "http://localhost:3000/uploads/user/avatar/1/cd8.jpeg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.FetchUrl(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitlabUrlParser.FetchUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitlabUrlParser.FetchUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
