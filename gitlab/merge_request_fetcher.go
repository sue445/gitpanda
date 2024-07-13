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

type mergeRequestFetcher struct {
}

func (f *mergeRequestFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(reProjectName + "/merge_requests/(\\d+)").FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := sanitizeProjectName(matched[1])

	var eg errgroup.Group

	var mr *gitlab.MergeRequest
	description := ""
	authorName := ""
	authorAvatarURL := ""
	var footerTime *time.Time
	eg.Go(func() error {
		var err error
		mrID, _ := strconv.Atoi(matched[2])
		start := time.Now()
		mr, _, err = client.MergeRequests.GetMergeRequest(projectName, mrID, nil)

		if err != nil {
			return errors.WithStack(err)
		}

		if isDebugLogging {
			duration := time.Since(start)
			fmt.Printf("[DEBUG] mergeRequestFetcher (%s): mr=%+v\n", duration, mr)
		}

		description = mr.Description
		authorName = mr.Author.Name
		authorAvatarURL = mr.Author.AvatarURL
		footerTime = mr.CreatedAt

		matched2 := regexp.MustCompile(`#note_(\d+)$`).FindStringSubmatch(path)

		if matched2 != nil {
			noteID, _ := strconv.Atoi(matched2[1])
			start := time.Now()
			note, _, err := client.Notes.GetMergeRequestNote(projectName, mrID, noteID)

			if err != nil {
				return errors.WithStack(err)
			}

			if isDebugLogging {
				duration := time.Since(start)
				fmt.Printf("[DEBUG] mergeRequestFetcher (%s): note=%+v\n", duration, note)
			}

			description = note.Body
			authorName = note.Author.Name
			authorAvatarURL = note.Author.AvatarURL
			footerTime = note.CreatedAt
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
			fmt.Printf("[DEBUG] mergeRequestFetcher (%s): project=%+v\n", duration, project)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	page := &Page{
		Title:                  strings.Join([]string{mr.Title, "Merge Requests", project.NameWithNamespace, "GitLab"}, titleSeparator),
		Description:            description,
		AuthorName:             authorName,
		AuthorAvatarURL:        authorAvatarURL,
		AvatarURL:              project.AvatarURL,
		CanTruncateDescription: true,
		FooterTitle:            project.PathWithNamespace,
		FooterURL:              project.WebURL,
		FooterTime:             footerTime,
	}

	return page, nil
}
