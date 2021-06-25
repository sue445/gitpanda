package gitlab

import (
	"bytes"
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/pkg/errors"
	"github.com/sue445/gitpanda/util"
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

	jobID, _ := strconv.Atoi(matched[2])

	var eg errgroup.Group

	var job *gitlab.Job
	eg.Go(func() error {
		var err error

		start := time.Now()
		job, _, err = client.Jobs.GetJob(projectName, jobID)

		if err != nil {
			return errors.WithStack(err)
		}

		if isDebugLogging {
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] jobFetcher (%s): job=%+v\n", duration, job)
		}

		return nil
	})

	lineRe := regexp.MustCompile(".*#L([0-9-]+)$")
	lineMatched := lineRe.FindStringSubmatch(path)

	selectedLine := ""
	if lineMatched != nil {
		eg.Go(func() error {
			var err error

			start := time.Now()
			reader, _, err := client.Jobs.GetTraceFile(projectName, jobID)

			if err != nil {
				return errors.WithStack(err)
			}

			if isDebugLogging {
				duration := time.Now().Sub(start)
				fmt.Printf("[DEBUG] jobFetcher (%s)\n", duration)
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(reader)
			traceBody := normalizeJobTrace(buf.String())

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
		start := time.Now()
		project, _, err = client.Projects.GetProject(projectName, nil)

		if err != nil {
			return errors.WithStack(err)
		}

		if isDebugLogging {
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] jobFetcher (%s): project=%+v\n", duration, project)
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
