package gitlab

import "github.com/xanzy/go-gitlab"

// ReProjectName represents regular expression matching to a project name
const ReProjectName = "^([^/]+(?:/[^/]+)+)"

type fetcher interface {
	fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error)
}
