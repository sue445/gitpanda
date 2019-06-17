package gitlab

import (
	"fmt"
	"github.com/sue445/gitpanda/util"
	"github.com/xanzy/go-gitlab"
	"regexp"
	"strconv"
	"strings"
)

const titleSeparator = " · "

// URLParser represents GitLab URL parser
type URLParser struct {
	baseURL        string
	client         *gitlab.Client
	isDebugLogging bool
}

// URLParserParams represents parameters of NewGitlabURLParser
type URLParserParams struct {
	APIEndpoint     string
	BaseURL         string
	PrivateToken    string
	GitPandaVersion string
	IsDebugLogging  bool
}

// NewGitlabURLParser create new URLParser instance
func NewGitlabURLParser(params *URLParserParams) (*URLParser, error) {
	p := &URLParser{
		baseURL:        params.BaseURL,
		isDebugLogging: params.IsDebugLogging,
	}

	p.client = gitlab.NewClient(nil, params.PrivateToken)
	err := p.client.SetBaseURL(params.APIEndpoint)

	if err != nil {
		return nil, err
	}

	p.client.UserAgent = fmt.Sprintf("gitpanda/%s (+https://github.com/sue445/gitpanda)", params.GitPandaVersion)

	return p, nil
}

// FetchURL fetch GitLab url
func (p *URLParser) FetchURL(url string) (*Page, error) {
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

	// Blob URL
	page, err = p.fetchBlobURL(path)

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

func (p *URLParser) fetchIssueURL(path string) (*Page, error) {
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

	if p.isDebugLogging {
		fmt.Printf("[DEBUG] fetchIssueURL: issue=%+v\n", issue)
	}

	project, _, err := p.client.Projects.GetProject(projectName, nil)

	if err != nil {
		return nil, err
	}

	if p.isDebugLogging {
		fmt.Printf("[DEBUG] fetchIssueURL: project=%+v\n", project)
	}

	description := issue.Description
	authorName := issue.Author.Name
	authorAvatarURL := issue.Author.AvatarURL
	footerTime := issue.CreatedAt

	re2 := regexp.MustCompile("#note_(\\d+)$")
	matched2 := re2.FindStringSubmatch(path)

	if matched2 != nil {
		noteID, _ := strconv.Atoi(matched2[1])
		note, _, err := p.client.Notes.GetIssueNote(projectName, issueID, noteID)

		if err != nil {
			return nil, err
		}

		if p.isDebugLogging {
			fmt.Printf("[DEBUG] fetchIssueURL: note=%+v\n", note)
		}

		description = note.Body
		authorName = note.Author.Name
		authorAvatarURL = note.Author.AvatarURL
		footerTime = note.CreatedAt
	}

	page := &Page{
		Title:                  strings.Join([]string{issue.Title, "Issues", project.NameWithNamespace, "GitLab"}, titleSeparator),
		Description:            description,
		AuthorName:             authorName,
		AuthorAvatarURL:        authorAvatarURL,
		AvatarURL:              project.AvatarURL,
		CanTruncateDescription: true,
		FooterTitle:            project.PathWithNamespace,
		FooterURL:              project.WebURL,
		FooterTime:             footerTime,
	}

	return page, nil
}

func (p *URLParser) fetchMergeRequestURL(path string) (*Page, error) {
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

	if p.isDebugLogging {
		fmt.Printf("[DEBUG] fetchMergeRequestURL: mr=%+v\n", mr)
	}

	project, _, err := p.client.Projects.GetProject(projectName, nil)

	if err != nil {
		return nil, err
	}

	if p.isDebugLogging {
		fmt.Printf("[DEBUG] fetchMergeRequestURL: project=%+v\n", project)
	}

	description := mr.Description
	authorName := mr.Author.Name
	authorAvatarURL := mr.Author.AvatarURL
	footerTime := mr.CreatedAt

	re2 := regexp.MustCompile("#note_(\\d+)$")
	matched2 := re2.FindStringSubmatch(path)

	if matched2 != nil {
		noteID, _ := strconv.Atoi(matched2[1])
		note, _, err := p.client.Notes.GetMergeRequestNote(projectName, mrID, noteID)

		if err != nil {
			return nil, err
		}

		if p.isDebugLogging {
			fmt.Printf("[DEBUG] fetchMergeRequestURL: note=%+v\n", note)
		}

		description = note.Body
		authorName = note.Author.Name
		authorAvatarURL = note.Author.AvatarURL
		footerTime = note.CreatedAt
	}

	page := &Page{
		Title:                  strings.Join([]string{mr.Title, "Merge Requests", project.NameWithNamespace, "GitLab"}, titleSeparator),
		Description:            description,
		AuthorName:             authorName,
		AuthorAvatarURL:        authorAvatarURL,
		AvatarURL:              project.AvatarURL,
		CanTruncateDescription: true,
		FooterTitle:            project.PathWithNamespace,
		FooterURL:              project.WebURL,
		FooterTime:             footerTime,
	}

	return page, nil
}

func (p *URLParser) fetchProjectURL(path string) (*Page, error) {
	re := regexp.MustCompile("^([^/]+)/([^/]+)/?$")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	project, _, err := p.client.Projects.GetProject(matched[1]+"/"+matched[2], nil)

	if err != nil {
		return nil, err
	}

	if p.isDebugLogging {
		fmt.Printf("[DEBUG] fetchProjectURL: project=%+v\n", project)
	}

	page := &Page{
		Title:                  strings.Join([]string{project.NameWithNamespace, "GitLab"}, titleSeparator),
		Description:            project.Description,
		AvatarURL:              project.AvatarURL,
		CanTruncateDescription: true,
		FooterTitle:            project.PathWithNamespace,
		FooterURL:              project.WebURL,
		FooterTime:             project.CreatedAt,
	}

	if project.Owner != nil {
		page.AuthorName = project.Owner.Name
		page.AuthorAvatarURL = project.Owner.AvatarURL
	}

	return page, nil
}

func (p *URLParser) fetchUserURL(path string) (*Page, error) {
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

	if p.isDebugLogging {
		fmt.Printf("[DEBUG] fetchUserURL: users=%+v\n", users)
	}

	if len(users) < 1 {
		return nil, fmt.Errorf("NotFound user: %s", username)
	}

	user := users[0]

	page := &Page{
		Title:                  strings.Join([]string{user.Name, "GitLab"}, titleSeparator),
		Description:            user.Name,
		AuthorName:             user.Name,
		AuthorAvatarURL:        user.AvatarURL,
		AvatarURL:              user.AvatarURL,
		CanTruncateDescription: true,
		FooterTitle:            fmt.Sprintf("@%s", user.Username),
		FooterURL:              fmt.Sprintf("%s/%s", p.baseURL, user.Username),
		FooterTime:             user.CreatedAt,
	}

	return page, nil
}

func (p *URLParser) fetchBlobURL(path string) (*Page, error) {
	re := regexp.MustCompile("^([^/]+)/([^/]+)/blob/([^/]+)/(.+)#L([0-9-]+)$")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := matched[1] + "/" + matched[2]
	sha1 := matched[3]
	fileName := matched[4]
	rawFile, _, err := p.client.RepositoryFiles.GetRawFile(projectName, fileName, &gitlab.GetRawFileOptions{Ref: &sha1})

	if err != nil {
		return nil, err
	}

	fileBody := string(rawFile)

	if p.isDebugLogging {
		fmt.Printf("[DEBUG] fetchBlobURL: fileBody=%s\n", fileBody)
	}

	lineHash := matched[5]
	lines := strings.Split(lineHash, "-")

	selectedFile := ""
	lineRange := ""
	switch len(lines) {
	case 1:
		line, _ := strconv.Atoi(lines[0])
		lineRange = lines[0]
		selectedFile = util.SelectLine(fileBody, line)
	case 2:
		startLine, _ := strconv.Atoi(lines[0])
		endLine, _ := strconv.Atoi(lines[1])
		lineRange = fmt.Sprintf("%s-%s", lines[0], lines[1])
		selectedFile = util.SelectLines(fileBody, startLine, endLine)
	default:
		return nil, fmt.Errorf("Invalid line: L%s", lineHash)
	}

	project, _, err := p.client.Projects.GetProject(projectName, nil)

	if err != nil {
		return nil, err
	}

	page := &Page{
		Title:                  fmt.Sprintf("%s:%s", fileName, lineRange),
		Description:            fmt.Sprintf("```\n%s\n```", selectedFile),
		AuthorName:             "",
		AuthorAvatarURL:        "",
		AvatarURL:              project.AvatarURL,
		CanTruncateDescription: false,
		FooterTitle:            project.PathWithNamespace,
		FooterURL:              project.WebURL,
		FooterTime:             nil,
	}

	return page, nil
}
