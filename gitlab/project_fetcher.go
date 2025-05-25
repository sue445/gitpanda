package gitlab

import (
	"github.com/cockroachdb/errors"
	"github.com/sue445/gitpanda/util"
	"gitlab.com/gitlab-org/api/client-go"
	"regexp"
	"strings"
)

type projectFetcher struct {
}

func (f *projectFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	// Remove anchor(#) in path (e.g. gitlab-org/gitlab#gitlab -> gitlab-org/gitlab)
	path = regexp.MustCompile("#.*$").ReplaceAllString(path, "")

	matched := regexp.MustCompile(reProjectName + "/?$").FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	project, err := util.WithDebugLogging("projectFetcher(GetProject)", isDebugLogging, func() (*gitlab.Project, error) {
		project, _, err := client.Projects.GetProject(matched[1], nil)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return project, nil
	})

	if err != nil {
		return nil, errors.WithStack(err)
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
