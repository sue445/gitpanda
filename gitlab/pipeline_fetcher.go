package gitlab

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type pipelineFetcher struct {
}

func (f *pipelineFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(reProjectName + "/pipelines/(\\d+)").FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := sanitizeProjectName(matched[1])

	var eg errgroup.Group

	var pipeline *gitlab.Pipeline
	eg.Go(func() error {
		var err error

		pipelineID, _ := strconv.Atoi(matched[2])

		start := time.Now()
		pipeline, _, err = client.Pipelines.GetPipeline(projectName, pipelineID)

		if err != nil {
			return errors.WithStack(err)
		}

		if isDebugLogging {
			duration := time.Since(start)
			fmt.Printf("[DEBUG] pipelineFetcher (%s): pipeline=%+v\n", duration, pipeline)
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
			duration := time.Since(start)
			fmt.Printf("[DEBUG] pipelineFetcher (%s): project=%+v\n", duration, project)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	duration := time.Duration(pipeline.Duration) * time.Second
	page := &Page{
		Title:                  strings.Join([]string{"Pipeline", project.NameWithNamespace, "GitLab"}, titleSeparator),
		Description:            fmt.Sprintf("[%s] Pipeline [#%d](%s) of branch %s by %s in %s", pipeline.Status, pipeline.ID, pipeline.WebURL, pipeline.Ref, pipeline.User.Username, duration),
		AuthorName:             pipeline.User.Name,
		AuthorAvatarURL:        pipeline.User.AvatarURL,
		AvatarURL:              "",
		CanTruncateDescription: true,
		FooterTitle:            project.PathWithNamespace,
		FooterURL:              project.WebURL,
		FooterTime:             pipeline.CreatedAt,
		Color:                  ciStatusColor(pipeline.Status),
	}

	return page, nil
}
