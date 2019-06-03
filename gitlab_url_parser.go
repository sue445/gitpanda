package main

import (
	"github.com/xanzy/go-gitlab"
	"regexp"
	"strconv"
	"strings"
)

type GitlabUrlParser struct {
	baseUrl string
	client  *gitlab.Client
}

type GitLabUrlParserParams struct {
	ApiEndpoint  string
	BaseUrl      string
	PrivateToken string
}

func NewGitlabUrlParser(params GitLabUrlParserParams) *GitlabUrlParser {
	p := &GitlabUrlParser{
		baseUrl: params.BaseUrl,
	}

	p.client = gitlab.NewClient(nil, params.PrivateToken)
	p.client.SetBaseURL(params.ApiEndpoint)

	return p
}

func (p *GitlabUrlParser) FetchUrl(url string) (*GitLabPage, error) {
	if !strings.HasPrefix(url, p.baseUrl) {
		return nil, nil
	}

	pos := len(p.baseUrl)
	if !strings.HasSuffix(url, "/") {
		pos++
	}
	path := url[pos:]

	// Issue URL
	page, err := p.fetchIssueUrl(path)

	if err != nil {
		return nil, err
	}

	if page != nil {
		return page, nil
	}

	// MergeRequest URL
	page, err = p.fetchMergeRequestUrl(path)

	if err != nil {
		return nil, err
	}

	if page != nil {
		return page, nil
	}

	// Project URL
	page, err = p.fetchProjectUrl(path)

	if err != nil {
		return nil, err
	}

	if page != nil {
		return page, nil
	}

	return nil, nil
}

func (p *GitlabUrlParser) fetchIssueUrl(path string) (*GitLabPage, error) {
	re := regexp.MustCompile("^(.+)/(.+)/issues/(\\d+)")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := matched[1] + "/" + matched[2]
	issueID, _ := strconv.Atoi(matched[3])
	issue, _, err := p.client.Issues.GetIssue(projectName, issueID)

	if err != nil {
		return nil, err
	}

	project, _, err := p.client.Projects.GetProject(projectName, nil)

	if err != nil {
		return nil, err
	}

	page := &GitLabPage{
		Title:            issue.Title,
		Description:      issue.Description,
		AuthorName:       issue.Author.Name,
		AuthorAvatarURL:  issue.Author.AvatarURL,
		ProjectAvatarURL: project.AvatarURL,
	}

	return page, nil
}

func (p *GitlabUrlParser) fetchMergeRequestUrl(path string) (*GitLabPage, error) {
	re := regexp.MustCompile("^(.+)/(.+)/merge_requests/(\\d+)")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := matched[1] + "/" + matched[2]
	mrID, _ := strconv.Atoi(matched[3])
	mr, _, err := p.client.MergeRequests.GetMergeRequest(projectName, mrID, nil)

	if err != nil {
		return nil, err
	}

	project, _, err := p.client.Projects.GetProject(projectName, nil)

	if err != nil {
		return nil, err
	}

	page := &GitLabPage{
		Title:            mr.Title,
		Description:      mr.Description,
		AuthorName:       mr.Author.Name,
		// AuthorAvatarURL:  mr.Author.AvatarURL, TODO: fix after
		AuthorAvatarURL: "",
		ProjectAvatarURL: project.AvatarURL,
	}

	return page, nil
}

func (p *GitlabUrlParser) fetchProjectUrl(path string) (*GitLabPage, error) {
	re := regexp.MustCompile("^(.+)/(.+)/?$")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	project, _, err := p.client.Projects.GetProject(matched[1]+"/"+matched[2], nil)

	if err != nil {
		return nil, err
	}

	page := &GitLabPage{
		Title:            project.NameWithNamespace,
		Description:      project.Description,
		AuthorName:       project.Owner.Name,
		AuthorAvatarURL:  project.Owner.AvatarURL,
		ProjectAvatarURL: project.AvatarURL,
	}

	return page, nil
}
