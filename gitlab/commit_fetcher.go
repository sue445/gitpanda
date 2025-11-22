package gitlab

import (
	"regexp"

	"github.com/cockroachdb/errors"
	"github.com/sue445/gitpanda/util"
	"gitlab.com/gitlab-org/api/client-go"
	"golang.org/x/sync/errgroup"
)

type commitFetcher struct {
}

func (f *commitFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(reProjectName + `/commit/(.+)`).FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := sanitizeProjectName(matched[1])

	sha := matched[2]

	var eg errgroup.Group

	var commit *gitlab.Commit
	eg.Go(func() error {
		var err error
		commit, err = util.WithDebugLogging("commitFetcher(GetCommit)", isDebugLogging, func() (*gitlab.Commit, error) {
			commit, _, err := client.Commits.GetCommit(projectName, sha, &gitlab.GetCommitOptions{})
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return commit, nil
		})
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})

	var project *gitlab.Project
	eg.Go(func() error {
		var err error
		project, err = util.WithDebugLogging("commitFetcher(GetProject)", isDebugLogging, func() (*gitlab.Project, error) {
			project, _, err := client.Projects.GetProject(projectName, nil)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return project, nil
		})
		if err != nil {
			return errors.WithStack(err)
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
