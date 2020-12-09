package gitlab

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
	"regexp"
	"strings"
	"time"
)

type projectFetcher struct {
}

func (f *projectFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	re := regexp.MustCompile(reProjectName + "/?$")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	start := time.Now()
	project, _, err := client.Projects.GetProject(matched[1], nil)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if isDebugLogging {
		duration := time.Now().Sub(start)
		fmt.Printf("[DEBUG] projectFetcher (%s): project=%+v\n", duration, project)
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
