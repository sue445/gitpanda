package gitlab

import "github.com/xanzy/go-gitlab"

// reProjectName represents regular expression matching to a project name
const reProjectName = "^([^/]+(?:/[^/]+)+)"

type fetcher interface {
	fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error)
}
