package gitlab

import "github.com/sue445/go-gitlab"

type fetcher interface {
	fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error)
}
