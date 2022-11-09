package gitlab

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
	"regexp"
	"time"
)

type commitFetcher struct {
}

func (f *commitFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(reProjectName + "/commit/(.+)").FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := sanitizeProjectName(matched[1])

	sha := matched[2]

	var eg errgroup.Group

	var commit *gitlab.Commit
	eg.Go(func() error {
		var err error

		start := time.Now()
		commit, _, err = client.Commits.GetCommit(projectName, sha)

		if err != nil {
			return errors.WithStack(err)
		}

		if isDebugLogging {
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] commitFetcher (%s): commit=%+v\n", duration, commit)
		}

		return nil
	})

	var project *gitlab.Project
	eg.Go(func() error {
		var err error
		start := time.Now()
		project, _, err = client.Projects.GetProject(projectName, nil)

		if err != nil {
			return errors.WithStack(err)
		}

		if isDebugLogging {
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] commitFetcher (%s): project=%+v\n", duration, project)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	page := &Page{
		Title:                  commit.Title,
		Description:            commit.Message,
		AuthorName:             commit.AuthorName,
		AuthorAvatarURL:        "",
		AvatarURL:              "",
		CanTruncateDescription: true,
		FooterTitle:            project.PathWithNamespace,
		FooterURL:              project.WebURL,
		FooterTime:             commit.CreatedAt,
		Color:                  "",
	}

	return page, nil
}
