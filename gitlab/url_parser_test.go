package gitlab

import (
	"github.com/stretchr/testify/assert"
	"github.com/sue445/gitpanda/testutil"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func tp(t time.Time) *time.Time {
	return &t
}

func TestGitlabUrlParser_FetchURL(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register stub
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/",
		func(_ *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(404, "{\"error\":\"404 Not Found\"}")
			resp.Header.Set("RateLimit-Limit", "600")
			return resp, nil
		},
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/project.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/gitlab-org%2Fdiaspora-project-site",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/project_without_owner.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/my-group%2Fdiaspora%2Fdiaspora-project-site",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/project_with_subgroup.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/issues/1",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/issue.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/issues/1/notes/302",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/issue_note.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/merge_requests/1",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/merge_request.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/merge_requests/1/notes/301",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/merge_request_note.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/users?username=john_smith",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/users.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/users?username=gitlab-org",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/users_by_group_name.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/groups/gitlab-org?with_projects=false",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/group.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/repository/files/gitlabci-templates%2Fcontinuous_bundle_update%2Eyml/raw?ref=master",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/gitlabci-templates/continuous_bundle_update.yml")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/repository/files/icon%2Epng/raw?ref=master",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/icon.png")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/jobs/8",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/job.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/jobs/8/trace",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/job_trace.txt")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/pipelines/46",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/pipeline.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/snippets/1",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/project_snippet.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/snippets/1/raw",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/snippet_code.rb")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/snippets/3",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/project_snippet.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/snippets/1/notes/400",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/project_snippet_note.json")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/snippets/3/raw",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/snippet_code.rb")),
	)
	httpmock.RegisterResponder(
		"GET",
		"http://example.com/api/v4/projects/diaspora%2Fdiaspora-project-site/repository/commits/c9ddb5f48418b4e2a9e41982b8177018114003d1",
		httpmock.NewStringResponder(200, testutil.ReadTestData("testdata/commit.json")),
	)

	p, err := NewGitlabURLParser(&URLParserParams{
		APIEndpoint:     "http://example.com/api/v4",
		BaseURL:         "http://example.com",
		PrivateToken:    "xxxxxxxxxx",
		GitPandaVersion: "v0.0.0",
		IsDebugLogging:  true,
		HTTPClient:      http.DefaultClient,
	})

	if err != nil {
		panic(err)
	}

	type args struct {
		url string
	}

	tests := []struct {
		name string
		args args
		want *Page
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
			want: &Page{
				Title:                  "Diaspora / Diaspora Project Site · GitLab",
				Description:            "Diaspora Project",
				AuthorName:             "Diaspora",
				AuthorAvatarURL:        "http://example.com/images/diaspora.png",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 9, 30, 13, 46, 2, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Project URL (ends with slash)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/",
			},
			want: &Page{
				Title:                  "Diaspora / Diaspora Project Site · GitLab",
				Description:            "Diaspora Project",
				AuthorName:             "Diaspora",
				AuthorAvatarURL:        "http://example.com/images/diaspora.png",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 9, 30, 13, 46, 2, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Project (without owner) URL",
			args: args{
				url: "http://example.com/gitlab-org/diaspora-project-site",
			},
			want: &Page{
				Title:                  "GitLab.org / Diaspora Project Site · GitLab",
				Description:            "Diaspora Project",
				AuthorName:             "",
				AuthorAvatarURL:        "",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "gitlab-org/diaspora-project-site",
				FooterURL:              "http://example.com/gitlab-org/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 9, 30, 13, 46, 2, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Project (with subgroup) URL",
			args: args{
				url: "http://example.com/my-group/diaspora/diaspora-project-site",
			},
			want: &Page{
				Title:                  "My Group / Diaspora / Diaspora Project Site · GitLab",
				Description:            "Diaspora Project",
				AuthorName:             "Diaspora",
				AuthorAvatarURL:        "http://example.com/images/diaspora.png",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "my-group/diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/my-group/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 9, 30, 13, 46, 2, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Project URL (with anchor)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site#readme",
			},
			want: &Page{
				Title:                  "Diaspora / Diaspora Project Site · GitLab",
				Description:            "Diaspora Project",
				AuthorName:             "Diaspora",
				AuthorAvatarURL:        "http://example.com/images/diaspora.png",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 9, 30, 13, 46, 2, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Issue URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/issues/1",
			},
			want: &Page{
				Title:                  "Ut commodi ullam eos dolores perferendis nihil sunt. · Issues · Diaspora / Diaspora Project Site · GitLab",
				Description:            "Omnis vero earum sunt corporis dolor et placeat.",
				AuthorName:             "Administrator",
				AuthorAvatarURL:        "https://gitlab.example.com/images/root.png",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2016, 1, 4, 15, 31, 46, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Issue URL (ends with slash)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/issues/1/",
			},
			want: &Page{
				Title:                  "Ut commodi ullam eos dolores perferendis nihil sunt. · Issues · Diaspora / Diaspora Project Site · GitLab",
				Description:            "Omnis vero earum sunt corporis dolor et placeat.",
				AuthorName:             "Administrator",
				AuthorAvatarURL:        "https://gitlab.example.com/images/root.png",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2016, 1, 4, 15, 31, 46, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Issue URL (including /-/)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/-/issues/1",
			},
			want: &Page{
				Title:                  "Ut commodi ullam eos dolores perferendis nihil sunt. · Issues · Diaspora / Diaspora Project Site · GitLab",
				Description:            "Omnis vero earum sunt corporis dolor et placeat.",
				AuthorName:             "Administrator",
				AuthorAvatarURL:        "https://gitlab.example.com/images/root.png",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2016, 1, 4, 15, 31, 46, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Issue comment URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/issues/1#note_302",
			},
			want: &Page{
				Title:                  "Ut commodi ullam eos dolores perferendis nihil sunt. · Issues · Diaspora / Diaspora Project Site · GitLab",
				Description:            "closed",
				AuthorName:             "Pip",
				AuthorAvatarURL:        "http://localhost:3000/uploads/user/avatar/1/pipin.jpeg",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 10, 2, 9, 22, 45, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Issue comment URL (ends with slash)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/issues/1/#note_302",
			},
			want: &Page{
				Title:                  "Ut commodi ullam eos dolores perferendis nihil sunt. · Issues · Diaspora / Diaspora Project Site · GitLab",
				Description:            "closed",
				AuthorName:             "Pip",
				AuthorAvatarURL:        "http://localhost:3000/uploads/user/avatar/1/pipin.jpeg",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 10, 2, 9, 22, 45, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Issue comment URL (including /-/)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/-/issues/1#note_302",
			},
			want: &Page{
				Title:                  "Ut commodi ullam eos dolores perferendis nihil sunt. · Issues · Diaspora / Diaspora Project Site · GitLab",
				Description:            "closed",
				AuthorName:             "Pip",
				AuthorAvatarURL:        "http://localhost:3000/uploads/user/avatar/1/pipin.jpeg",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 10, 2, 9, 22, 45, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "MergeRequest URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/merge_requests/1",
			},
			want: &Page{
				Title:                  "test1 · Merge Requests · Diaspora / Diaspora Project Site · GitLab",
				Description:            "fixed login page css paddings",
				AuthorName:             "Administrator",
				AuthorAvatarURL:        "https://gitlab.example.com/images/admin.png",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2017, 4, 29, 8, 46, 0, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "MergeRequest URL (ends with slash)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/merge_requests/1/",
			},
			want: &Page{
				Title:                  "test1 · Merge Requests · Diaspora / Diaspora Project Site · GitLab",
				Description:            "fixed login page css paddings",
				AuthorName:             "Administrator",
				AuthorAvatarURL:        "https://gitlab.example.com/images/admin.png",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2017, 4, 29, 8, 46, 0, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "MergeRequest URL (including /-/)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/-/merge_requests/1",
			},
			want: &Page{
				Title:                  "test1 · Merge Requests · Diaspora / Diaspora Project Site · GitLab",
				Description:            "fixed login page css paddings",
				AuthorName:             "Administrator",
				AuthorAvatarURL:        "https://gitlab.example.com/images/admin.png",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2017, 4, 29, 8, 46, 0, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "MergeRequest comment URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/merge_requests/1#note_301",
			},
			want: &Page{
				Title:                  "test1 · Merge Requests · Diaspora / Diaspora Project Site · GitLab",
				Description:            "Comment for MR",
				AuthorName:             "Pip",
				AuthorAvatarURL:        "http://localhost:3000/uploads/user/avatar/1/pipin.jpeg",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 10, 2, 8, 57, 14, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "MergeRequest comment URL (ends with slash)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/merge_requests/1/#note_301",
			},
			want: &Page{
				Title:                  "test1 · Merge Requests · Diaspora / Diaspora Project Site · GitLab",
				Description:            "Comment for MR",
				AuthorName:             "Pip",
				AuthorAvatarURL:        "http://localhost:3000/uploads/user/avatar/1/pipin.jpeg",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 10, 2, 8, 57, 14, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "MergeRequest comment URL (including /-/)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/-/merge_requests/1#note_301",
			},
			want: &Page{
				Title:                  "test1 · Merge Requests · Diaspora / Diaspora Project Site · GitLab",
				Description:            "Comment for MR",
				AuthorName:             "Pip",
				AuthorAvatarURL:        "http://localhost:3000/uploads/user/avatar/1/pipin.jpeg",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 10, 2, 8, 57, 14, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "User URL",
			args: args{
				url: "http://example.com/john_smith",
			},
			want: &Page{
				Title:                  "John Smith · GitLab",
				Description:            "John Smith",
				AuthorName:             "John Smith",
				AuthorAvatarURL:        "http://localhost:3000/uploads/user/avatar/1/cd8.jpeg",
				AvatarURL:              "http://localhost:3000/uploads/user/avatar/1/cd8.jpeg",
				CanTruncateDescription: true,
				FooterTitle:            "@john_smith",
				FooterURL:              "http://localhost:3000/john_smith",
				FooterTime:             tp(time.Date(2012, 5, 23, 8, 0, 58, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "User URL (ends with slash)",
			args: args{
				url: "http://example.com/john_smith/",
			},
			want: &Page{
				Title:                  "John Smith · GitLab",
				Description:            "John Smith",
				AuthorName:             "John Smith",
				AuthorAvatarURL:        "http://localhost:3000/uploads/user/avatar/1/cd8.jpeg",
				AvatarURL:              "http://localhost:3000/uploads/user/avatar/1/cd8.jpeg",
				CanTruncateDescription: true,
				FooterTitle:            "@john_smith",
				FooterURL:              "http://localhost:3000/john_smith",
				FooterTime:             tp(time.Date(2012, 5, 23, 8, 0, 58, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Group URL",
			args: args{
				url: "http://example.com/gitlab-org",
			},
			want: &Page{
				Title:                  "GitLab.org · GitLab",
				Description:            "Open source software to collaborate on code",
				AuthorName:             "",
				AuthorAvatarURL:        "",
				AvatarURL:              "https://assets.gitlab-static.net/uploads/-/system/group/avatar/9970/logo-extra-whitespace.png",
				CanTruncateDescription: true,
				FooterTitle:            "@gitlab-org",
				FooterURL:              "https://gitlab.com/groups/gitlab-org",
				FooterTime:             nil,
				Color:                  "",
			},
		},
		{
			name: "Group URL (ends with slash)",
			args: args{
				url: "http://example.com/gitlab-org/",
			},
			want: &Page{
				Title:                  "GitLab.org · GitLab",
				Description:            "Open source software to collaborate on code",
				AuthorName:             "",
				AuthorAvatarURL:        "",
				AvatarURL:              "https://assets.gitlab-static.net/uploads/-/system/group/avatar/9970/logo-extra-whitespace.png",
				CanTruncateDescription: true,
				FooterTitle:            "@gitlab-org",
				FooterURL:              "https://gitlab.com/groups/gitlab-org",
				FooterTime:             nil,
				Color:                  "",
			},
		},
		{
			name: "Blob URL (single line)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/blob/master/gitlabci-templates/continuous_bundle_update.yml#L4",
			},
			want: &Page{
				Title:                  "gitlabci-templates/continuous_bundle_update.yml:4",
				Description:            "```\n  variables:\n```",
				AuthorName:             "",
				AuthorAvatarURL:        "",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: false,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             nil,
				Color:                  "",
			},
		},
		{
			name: "Blob URL (single line & including /-/)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/-/blob/master/gitlabci-templates/continuous_bundle_update.yml#L4",
			},
			want: &Page{
				Title:                  "gitlabci-templates/continuous_bundle_update.yml:4",
				Description:            "```\n  variables:\n```",
				AuthorName:             "",
				AuthorAvatarURL:        "",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: false,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             nil,
				Color:                  "",
			},
		},
		{
			name: "Blob URL (single line & including /-/ & conain url_params)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/-/blob/master/gitlabci-templates/continuous_bundle_update.yml?ref_type=head#L4",
			},
			want: &Page{
				Title:                  "gitlabci-templates/continuous_bundle_update.yml:4",
				Description:            "```\n  variables:\n```",
				AuthorName:             "",
				AuthorAvatarURL:        "",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: false,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             nil,
				Color:                  "",
			},
		},
		{
			name: "Blob URL (multiple line)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/blob/master/gitlabci-templates/continuous_bundle_update.yml#L4-9",
			},
			want: &Page{
				Title:                  "gitlabci-templates/continuous_bundle_update.yml:4-9",
				Description:            "```\n  variables:\n    CACHE_VERSION: \"v1\"\n    GIT_EMAIL:     \"gitlabci@example.com\"\n    GIT_USER:      \"GitLab CI\"\n    LABELS:        \"bundle update\"\n    OPTIONS:       \"\"\n```",
				AuthorName:             "",
				AuthorAvatarURL:        "",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: false,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             nil,
				Color:                  "",
			},
		},
		{
			name: "Blob URL (without line hash)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/blob/master/gitlabci-templates/continuous_bundle_update.yml",
			},
			want: &Page{
				Title:                  "gitlabci-templates/continuous_bundle_update.yml",
				Description:            "```\ncontinuous_bundle_update:\n  image: ruby\n\n  variables:\n    CACHE_VERSION: \"v1\"\n    GIT_EMAIL:     \"gitlabci@example.com\"\n    GIT_USER:      \"GitLab CI\"\n    LABELS:        \"bundle update\"\n    OPTIONS:       \"\"\n\n  cache:\n    key: \"$CACHE_VERSION-$CI_BUILD_NAME\"\n    paths:\n      - vendor/bundle/\n\n  script:\n    - bundle install --path vendor/bundle --clean\n    - gem install --no-doc gitlabci-bundle-update-mr\n    - gitlabci-bundle-update-mr --user=\"$GIT_USER\" --email=\"$GIT_EMAIL\" --labels=\"$LABELS\" $OPTIONS\n\n  only:\n    - schedules\n\n```",
				AuthorName:             "",
				AuthorAvatarURL:        "",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: false,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             nil,
				Color:                  "",
			},
		},
		{
			name: "Blob URL (non UTF-8 file)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/blob/master/icon.png",
			},
			want: &Page{
				Title:                  "icon.png",
				Description:            "",
				AuthorName:             "",
				AuthorAvatarURL:        "",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: false,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             nil,
				Color:                  "",
			},
		},
		{
			name: "Job URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/jobs/8",
			},
			want: &Page{
				Title:                  "rubocop (#8) · Jobs · Diaspora / Diaspora Project Site · GitLab",
				Description:            "[failed] Job [#8](http://example.com/diaspora/diaspora-project-site/-/jobs/8) of branch master by root in 1s",
				AuthorName:             "Administrator",
				AuthorAvatarURL:        "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=80&d=identicon",
				AvatarURL:              "",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2015, 12, 24, 15, 51, 21, 0, time.UTC)),
				Color:                  "#db3b21",
			},
		},
		{
			name: "Job URL (including /-/)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/-/jobs/8",
			},
			want: &Page{
				Title:                  "rubocop (#8) · Jobs · Diaspora / Diaspora Project Site · GitLab",
				Description:            "[failed] Job [#8](http://example.com/diaspora/diaspora-project-site/-/jobs/8) of branch master by root in 1s",
				AuthorName:             "Administrator",
				AuthorAvatarURL:        "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=80&d=identicon",
				AvatarURL:              "",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2015, 12, 24, 15, 51, 21, 0, time.UTC)),
				Color:                  "#db3b21",
			},
		},
		{
			name: "Job URL (single line)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/jobs/8#L6",
			},
			want: &Page{
				Title:                  "rubocop (#8) · Jobs · Diaspora / Diaspora Project Site · GitLab",
				Description:            "[failed] Job [#8](http://example.com/diaspora/diaspora-project-site/-/jobs/8) of branch master by root in 1s\n```\nPreparing the \"docker+machine\" executor\n```",
				AuthorName:             "Administrator",
				AuthorAvatarURL:        "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=80&d=identicon",
				AvatarURL:              "",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2015, 12, 24, 15, 51, 21, 0, time.UTC)),
				Color:                  "#db3b21",
			},
		},
		{
			name: "Job URL (multiple line)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/jobs/8#L4-6",
			},
			want: &Page{
				Title:                  "rubocop (#8) · Jobs · Diaspora / Diaspora Project Site · GitLab",
				Description:            "[failed] Job [#8](http://example.com/diaspora/diaspora-project-site/-/jobs/8) of branch master by root in 1s\n```\nResolving secrets\n\nPreparing the \"docker+machine\" executor\n```",
				AuthorName:             "Administrator",
				AuthorAvatarURL:        "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=80&d=identicon",
				AvatarURL:              "",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2015, 12, 24, 15, 51, 21, 0, time.UTC)),
				Color:                  "#db3b21",
			},
		},
		{
			name: "Pipeline URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/pipelines/46",
			},
			want: &Page{
				Title:                  "Pipeline · Diaspora / Diaspora Project Site · GitLab",
				Description:            "[success] Pipeline [#46](https://example.com/foo/bar/pipelines/46) of branch master by root in 0s",
				AuthorName:             "Administrator",
				AuthorAvatarURL:        "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=80&d=identicon",
				AvatarURL:              "",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2016, 8, 11, 11, 28, 34, 0, time.UTC)),
				Color:                  "#1aaa55",
			},
		},
		{
			name: "Pipeline URL (including /-/)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/-/pipelines/46",
			},
			want: &Page{
				Title:                  "Pipeline · Diaspora / Diaspora Project Site · GitLab",
				Description:            "[success] Pipeline [#46](https://example.com/foo/bar/pipelines/46) of branch master by root in 0s",
				AuthorName:             "Administrator",
				AuthorAvatarURL:        "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=80&d=identicon",
				AvatarURL:              "",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2016, 8, 11, 11, 28, 34, 0, time.UTC)),
				Color:                  "#1aaa55",
			},
		},
		{
			name: "Project Snippet URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/snippets/1",
			},
			want: &Page{
				Title:                  "add.rb",
				Description:            "```\nputs 'Hello'\n```",
				AuthorName:             "John Smith",
				AuthorAvatarURL:        "",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: false,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2012, 6, 28, 10, 52, 4, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Project Snippet URL (including /-/)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/-/snippets/1",
			},
			want: &Page{
				Title:                  "add.rb",
				Description:            "```\nputs 'Hello'\n```",
				AuthorName:             "John Smith",
				AuthorAvatarURL:        "",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: false,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2012, 6, 28, 10, 52, 4, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Project Snippet comment URL",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/snippets/1#note_400",
			},
			want: &Page{
				Title:                  "add.rb",
				Description:            "comment",
				AuthorName:             "Pip",
				AuthorAvatarURL:        "http://localhost:3000/uploads/user/avatar/1/pipin.jpeg",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 10, 2, 9, 22, 45, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Project Snippet comment URL (including /-/)",
			args: args{
				url: "http://example.com/diaspora/diaspora-project-site/-/snippets/1#note_400",
			},
			want: &Page{
				Title:                  "add.rb",
				Description:            "comment",
				AuthorName:             "Pip",
				AuthorAvatarURL:        "http://localhost:3000/uploads/user/avatar/1/pipin.jpeg",
				AvatarURL:              "http://example.com/uploads/project/avatar/3/uploads/avatar.png",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2013, 10, 2, 9, 22, 45, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Snippet URL",
			args: args{
				url: "http://example.com/snippets/3",
			},
			want: &Page{
				Title:                  "add.rb",
				Description:            "```\nputs 'Hello'\n```",
				AuthorName:             "John Smith",
				AuthorAvatarURL:        "",
				AvatarURL:              "",
				CanTruncateDescription: false,
				FooterTitle:            "",
				FooterURL:              "",
				FooterTime:             tp(time.Date(2012, 6, 28, 10, 52, 4, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Snippet URL (including /-/)",
			args: args{
				url: "http://example.com/-/snippets/3",
			},
			want: &Page{
				Title:                  "add.rb",
				Description:            "```\nputs 'Hello'\n```",
				AuthorName:             "John Smith",
				AuthorAvatarURL:        "",
				AvatarURL:              "",
				CanTruncateDescription: false,
				FooterTitle:            "",
				FooterURL:              "",
				FooterTime:             tp(time.Date(2012, 6, 28, 10, 52, 4, 0, time.UTC)),
				Color:                  "",
			},
		},
		{
			name: "Commit URL",
			args: args{
				// ref. https://gitlab.com/gitlab-org/gitlab/-/commit/c9ddb5f48418b4e2a9e41982b8177018114003d1
				url: "http://example.com/diaspora/diaspora-project-site/-/commit/c9ddb5f48418b4e2a9e41982b8177018114003d1",
			},
			want: &Page{
				Title:                  "Merge branch 'renovate-workhorse/github.com-prometheus-client_golang-1.x' into 'master'",
				Description:            "Merge branch 'renovate-workhorse/github.com-prometheus-client_golang-1.x' into 'master'\n\nworkhorse: Update module github.com/prometheus/client_golang to v1.14.0\n\nSee merge request https://gitlab.com/gitlab-org/gitlab/-/merge_requests/103423\n\nMerged-by: Jacob Vosmaer <jacob@gitlab.com>\nApproved-by: Stan Hu <stanhu@gmail.com>\nCo-authored-by: GitLab Renovate Bot <gitlab-bot@gitlab.com>",
				AuthorName:             "Jacob Vosmaer",
				AuthorAvatarURL:        "",
				AvatarURL:              "",
				CanTruncateDescription: true,
				FooterTitle:            "diaspora/diaspora-project-site",
				FooterURL:              "http://example.com/diaspora/diaspora-project-site",
				FooterTime:             tp(time.Date(2022, 11, 9, 14, 3, 6, 0, time.UTC)),
				Color:                  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.FetchURL(tt.args.url)

			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
