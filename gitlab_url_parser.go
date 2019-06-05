package main

import (
	"fmt"

	//"github.com/xanzy/go-gitlab"
	"github.com/sue445/go-gitlab"
	"regexp"
	"strconv"
	"strings"
)

type GitlabUrlParser struct {
	baseURL string
	client  *gitlab.Client
}

type GitLabUrlParserParams struct {
	ApiEndpoint  string
	BaseUrl      string
	PrivateToken string
}

func NewGitlabUrlParser(params GitLabUrlParserParams) (*GitlabUrlParser, error) {
	p := &GitlabUrlParser{
		baseURL: params.BaseUrl,
	}

	p.client = gitlab.NewClient(nil, params.PrivateToken)
	err := p.client.SetBaseURL(params.ApiEndpoint)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *GitlabUrlParser) FetchURL(url string) (*GitLabPage, error) {
	if !strings.HasPrefix(url, p.baseURL) {
		return nil, nil
	}

	pos := len(p.baseURL)
	if !strings.HasSuffix(url, "/") {
		pos++
	}
	path := url[pos:]

	// Issue URL
	page, err := p.fetchIssueURL(path)

	if err != nil {
		return nil, err
	}

	if page != nil {
		return page, nil
	}

	// MergeRequest URL
	page, err = p.fetchMergeRequestURL(path)

	if err != nil {
		return nil, err
	}

	if page != nil {
		return page, nil
	}

	// Project URL
	page, err = p.fetchProjectURL(path)

	if err != nil {
		return nil, err
	}

	if page != nil {
		return page, nil
	}

	// User URL
	page, err = p.fetchUserURL(path)

	if err != nil {
		return nil, err
	}

	if page != nil {
		return page, nil
	}

	return nil, nil
}

func (p *GitlabUrlParser) fetchIssueURL(path string) (*GitLabPage, error) {
	re := regexp.MustCompile("^([^/]+)/([^/]+)/issues/(\\d+)")
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
		Title:           issue.Title,
		Description:     issue.Description,
		AuthorName:      issue.Author.Name,
		AuthorAvatarURL: issue.Author.AvatarURL,
		AvatarURL:       project.AvatarURL,
	}

	return page, nil
}

func (p *GitlabUrlParser) fetchMergeRequestURL(path string) (*GitLabPage, error) {
	re := regexp.MustCompile("^([^/]+)/([^/]+)/merge_requests/(\\d+)")
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
		Title:           mr.Title,
		Description:     mr.Description,
		AuthorName:      mr.Author.Name,
		AuthorAvatarURL: mr.Author.AvatarURL,
		AvatarURL:       project.AvatarURL,
	}

	return page, nil
}

func (p *GitlabUrlParser) fetchProjectURL(path string) (*GitLabPage, error) {
	re := regexp.MustCompile("^([^/]+)/([^/]+)/?$")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	project, _, err := p.client.Projects.GetProject(matched[1]+"/"+matched[2], nil)

	if err != nil {
		return nil, err
	}

	page := &GitLabPage{
		Title:           project.NameWithNamespace,
		Description:     project.Description,
		AuthorName:      project.Owner.Name,
		AuthorAvatarURL: project.Owner.AvatarURL,
		AvatarURL:       project.AvatarURL,
	}

	return page, nil
}

func (p *GitlabUrlParser) fetchUserURL(path string) (*GitLabPage, error) {
	re := regexp.MustCompile("^([^/]+)/?$")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	username := matched[1]
	users, _, err := p.client.Users.ListUsers(&gitlab.ListUsersOptions{Username: &username})

	if err != nil {
		return nil, err
	}

	if len(users) < 1 {
		return nil, fmt.Errorf("NotFound user: %s", username)
	}

	user := users[0]

	page := &GitLabPage{
		Title:           user.Name,
		Description:     user.Name,
		AuthorName:      user.Name,
		AuthorAvatarURL: user.AvatarURL,
		AvatarURL:       user.AvatarURL,
	}

	return page, nil
}
