package gitlab

import (
	"fmt"
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
	re := regexp.MustCompile("^([^/]+)/([^/]+)/pipelines/(\\d+)")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := matched[1] + "/" + matched[2]

	var eg errgroup.Group

	var pipeline *gitlab.Pipeline
	eg.Go(func() error {
		var err error

		pipelineID, _ := strconv.Atoi(matched[3])

		start := time.Now()
		pipeline, _, err = client.Pipelines.GetPipeline(projectName, pipelineID)

		if err != nil {
			return err
		}

		if isDebugLogging {
			duration := time.Now().Sub(start)
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
			return err
		}

		if isDebugLogging {
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] pipelineFetcher (%s): project=%+v\n", duration, project)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	page := &Page{
		Title:                  strings.Join([]string{"Pipeline", project.NameWithNamespace, "GitLab"}, titleSeparator),
		Description:            "",
		AuthorName:             pipeline.User.Name,
		AuthorAvatarURL:        pipeline.User.AvatarURL,
		AvatarURL:              "",
		CanTruncateDescription: true,
		FooterTitle:            project.PathWithNamespace,
		FooterURL:              project.WebURL,
		FooterTime:             pipeline.CreatedAt,
	}

	return page, nil
}
