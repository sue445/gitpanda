package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
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
	if !strings.HasSuffix(p.baseURL, "/") {
		pos++
	}
	path := url[pos:]

	fetchers := []fetcher{
		&issueFetcher{},
		&mergeRequestFetcher{},
		&blobFetcher{},
		&projectFetcher{},
		&userFetcher{baseURL: p.baseURL},
	}

	for _, fetcher := range fetchers {
		page, err := fetcher.fetchPath(path, p.client, p.isDebugLogging)

		if err != nil {
			return nil, err
		}

		if page != nil {
			return page, nil
		}
	}

	return nil, nil
}
