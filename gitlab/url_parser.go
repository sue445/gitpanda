package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"net/http"
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
	HTTPClient      *http.Client
}

// NewGitlabURLParser create new URLParser instance
func NewGitlabURLParser(params *URLParserParams) (*URLParser, error) {
	p := &URLParser{
		baseURL:        params.BaseURL,
		isDebugLogging: params.IsDebugLogging,
	}

	options := []gitlab.ClientOptionFunc{gitlab.WithBaseURL(params.APIEndpoint)}
	if params.HTTPClient != nil {
		options = append(options, gitlab.WithHTTPClient(params.HTTPClient))
	}
	client, err := gitlab.NewClient(params.PrivateToken, options...)

	if err != nil {
		return nil, err
	}

	p.client = client
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
		&snippetFetcher{},
		&issueFetcher{},
		&mergeRequestFetcher{},
		&jobFetcher{},
		&pipelineFetcher{},
		&blobFetcher{},
		&projectSnippetFetcher{},
		&projectFetcher{},
		&userOrGroupFetcher{baseURL: p.baseURL},
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
