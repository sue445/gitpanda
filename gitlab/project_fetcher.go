package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"regexp"
	"strings"
)

type projectFetcher struct {
}

func (f *projectFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	re := regexp.MustCompile("^([^/]+)/([^/]+)/?$")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	project, _, err := client.Projects.GetProject(matched[1]+"/"+matched[2], nil)

	if err != nil {
		return nil, err
	}

	if isDebugLogging {
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
