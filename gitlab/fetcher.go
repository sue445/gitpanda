package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"regexp"
)

// reProjectName represents regular expression matching to a project name
const reProjectName = "^([^/]+(?:/[^/]+)+)"

type fetcher interface {
	fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error)
}

func sanitizeProjectName(projectName string) string {
	re := regexp.MustCompile(`(/-)$`)
	return re.ReplaceAllString(projectName, "")
}
