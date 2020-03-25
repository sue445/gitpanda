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

type jobFetcher struct {
}

func (f *jobFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	re := regexp.MustCompile(reProjectName + "/jobs/(\\d+)")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := sanitizeProjectName(matched[1])

	var eg errgroup.Group

	var job *gitlab.Job
	eg.Go(func() error {
		var err error

		jobID, _ := strconv.Atoi(matched[2])

		start := time.Now()
		job, _, err = client.Jobs.GetJob(projectName, jobID)

		if err != nil {
			return err
		}

		if isDebugLogging {
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] jobFetcher (%s): job=%+v\n", duration, job)
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
			fmt.Printf("[DEBUG] jobFetcher (%s): project=%+v\n", duration, project)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	duration := time.Duration(job.Duration) * time.Second
	page := &Page{
		Title:                  strings.Join([]string{fmt.Sprintf("%s (#%d)", job.Name, job.ID), "Jobs", project.NameWithNamespace, "GitLab"}, titleSeparator),
		Description:            fmt.Sprintf("[%s] Job [#%d](%s) of branch %s by %s in %s", job.Status, job.ID, job.WebURL, job.Ref, job.User.Username, duration),
		AuthorName:             job.User.Name,
		AuthorAvatarURL:        job.User.AvatarURL,
		AvatarURL:              "",
		CanTruncateDescription: true,
		FooterTitle:            project.PathWithNamespace,
		FooterURL:              project.WebURL,
		FooterTime:             job.CreatedAt,
		Color:                  ciStatusColor(job.Status),
	}

	return page, nil
}
