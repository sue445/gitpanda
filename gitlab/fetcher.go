package gitlab

import "github.com/xanzy/go-gitlab"

type fetcher interface {
	fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error)
}
