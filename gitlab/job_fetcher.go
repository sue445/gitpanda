package gitlab

import (
	"bytes"
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/cockroachdb/errors"
	"github.com/sue445/gitpanda/util"
	"gitlab.com/gitlab-org/api/client-go"
	"golang.org/x/sync/errgroup"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type jobFetcher struct {
}

func (f *jobFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(reProjectName + "/jobs/(\\d+)").FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := sanitizeProjectName(matched[1])

	jobID, _ := strconv.Atoi(matched[2])

	var eg errgroup.Group

	var job *gitlab.Job
	eg.Go(func() error {
		var err error
		job, err = util.WithDebugLogging("jobFetcher(GetJob)", isDebugLogging, func() (*gitlab.Job, error) {
			job, _, err := client.Jobs.GetJob(projectName, jobID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return job, nil
		})
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})

	lineMatched := regexp.MustCompile(".*#L([0-9-]+)$").FindStringSubmatch(path)

	selectedLine := ""
	if lineMatched != nil {
		eg.Go(func() error {
			var err error
			body, err := util.WithDebugLogging("jobFetcher(GetTraceFile)", isDebugLogging, func() (*string, error) {
				reader, _, err := client.Jobs.GetTraceFile(projectName, jobID)
				if err != nil {
					return nil, errors.WithStack(err)
				}

				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(reader)
				if err != nil {
					return nil, errors.WithStack(err)
				}
				body := buf.String()
				return &body, nil
			})
			if err != nil {
				return errors.WithStack(err)
			}

			traceBody := normalizeJobTrace(*body)

			lineHash := lineMatched[1]
			lines := strings.Split(lineHash, "-")

			switch len(lines) {
			case 1:
				line, _ := strconv.Atoi(lines[0])
				selectedLine = util.SelectLine(traceBody, line)
			case 2:
				startLine, _ := strconv.Atoi(lines[0])
				endLine, _ := strconv.Atoi(lines[1])
				selectedLine = util.SelectLines(traceBody, startLine, endLine)
			default:
				return fmt.Errorf("invalid line: L%s", lineHash)
			}
			return nil
		})
	}

	var project *gitlab.Project
	eg.Go(func() error {
		var err error
		project, err = util.WithDebugLogging("jobFetcher(GetProject)", isDebugLogging, func() (*gitlab.Project, error) {
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

	duration := time.Duration(job.Duration) * time.Second
	description := fmt.Sprintf("[%s] Job [#%d](%s) of branch %s by %s in %s", job.Status, job.ID, job.WebURL, job.Ref, job.User.Username, duration)

	if selectedLine != "" {
		description += fmt.Sprintf("\n```\n%s\n```", selectedLine)
	}

	page := &Page{
		Title:                  strings.Join([]string{fmt.Sprintf("%s (#%d)", job.Name, job.ID), "Jobs", project.NameWithNamespace, "GitLab"}, titleSeparator),
		Description:            description,
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

func normalizeJobTrace(raw string) string {
	normalized := stripansi.Strip(raw)

	normalized = regexp.MustCompile(` +\n`).ReplaceAllString(normalized, "\n")
	normalized = regexp.MustCompile(`\n+`).ReplaceAllString(normalized, "\n")
	normalized = regexp.MustCompile(`section_start:.+?\r`).ReplaceAllString(normalized, "")
	normalized = regexp.MustCompile(`section_end:.+?\r`).ReplaceAllString(normalized, "\n")

	return normalized
}
