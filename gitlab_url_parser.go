package main

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"regexp"
	"strconv"
	"strings"
)

const titleSeparator = " · "

// GitlabURLParser represents GitLab URL parser
type GitlabURLParser struct {
	baseURL string
	client  *gitlab.Client
}

// GitLabURLParserParams represents parameters of NewGitlabURLParser
type GitLabURLParserParams struct {
	APIEndpoint  string
	BaseURL      string
	PrivateToken string
}

// NewGitlabURLParser create new GitlabURLParser instance
func NewGitlabURLParser(params *GitLabURLParserParams) (*GitlabURLParser, error) {
	p := &GitlabURLParser{
		baseURL: params.BaseURL,
	}

	p.client = gitlab.NewClient(nil, params.PrivateToken)
	err := p.client.SetBaseURL(params.APIEndpoint)

	if err != nil {
		return nil, err
	}

	return p, nil
}

// FetchURL fetch GitLab url
func (p *GitlabURLParser) FetchURL(url string) (*GitLabPage, error) {
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

func (p *GitlabURLParser) fetchIssueURL(path string) (*GitLabPage, error) {
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

	description := issue.Description
	authorName := issue.Author.Name
	authorAvatarURL := issue.Author.AvatarURL

	re2 := regexp.MustCompile("#note_(\\d+)$")
	matched2 := re2.FindStringSubmatch(path)

	if matched2 != nil {
		noteID, _ := strconv.Atoi(matched2[1])
		note, _, err := p.client.Notes.GetIssueNote(projectName, issueID, noteID)

		if err != nil {
			return nil, err
		}

		description = note.Body
		authorName = note.Author.Name
		authorAvatarURL = note.Author.AvatarURL
	}

	page := &GitLabPage{
		Title:           strings.Join([]string{issue.Title, "Issues", project.NameWithNamespace, "GitLab"}, titleSeparator),
		Description:     description,
		AuthorName:      authorName,
		AuthorAvatarURL: authorAvatarURL,
		AvatarURL:       project.AvatarURL,
	}

	return page, nil
}

func (p *GitlabURLParser) fetchMergeRequestURL(path string) (*GitLabPage, error) {
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

	description := mr.Description
	authorName := mr.Author.Name
	authorAvatarURL := mr.Author.AvatarURL

	re2 := regexp.MustCompile("#note_(\\d+)$")
	matched2 := re2.FindStringSubmatch(path)

	if matched2 != nil {
		noteID, _ := strconv.Atoi(matched2[1])
		note, _, err := p.client.Notes.GetMergeRequestNote(projectName, mrID, noteID)

		if err != nil {
			return nil, err
		}

		description = note.Body
		authorName = note.Author.Name
		authorAvatarURL = note.Author.AvatarURL
	}

	page := &GitLabPage{
		Title:           strings.Join([]string{mr.Title, "Merge Requests", project.NameWithNamespace, "GitLab"}, titleSeparator),
		Description:     description,
		AuthorName:      authorName,
		AuthorAvatarURL: authorAvatarURL,
		AvatarURL:       project.AvatarURL,
	}

	return page, nil
}

func (p *GitlabURLParser) fetchProjectURL(path string) (*GitLabPage, error) {
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
		Title:           strings.Join([]string{project.NameWithNamespace, "GitLab"}, titleSeparator),
		Description:     project.Description,
		AuthorName:      project.Owner.Name,
		AuthorAvatarURL: project.Owner.AvatarURL,
		AvatarURL:       project.AvatarURL,
	}

	return page, nil
}

func (p *GitlabURLParser) fetchUserURL(path string) (*GitLabPage, error) {
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
		Title:           strings.Join([]string{user.Name, "GitLab"}, titleSeparator),
		Description:     user.Name,
		AuthorName:      user.Name,
		AuthorAvatarURL: user.AvatarURL,
		AvatarURL:       user.AvatarURL,
	}

	return page, nil
}
