package gitlab

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sue445/gitpanda/util"
	"gitlab.com/gitlab-org/api/client-go"
	"golang.org/x/sync/errgroup"
)

type mergeRequestFetcher struct {
}

func (f *mergeRequestFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(reProjectName + `/merge_requests/(\d+)`).FindStringSubmatch(path)

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
		mrID, _ := strconv.ParseInt(matched[2], 10, 64)
		mr, err = util.WithDebugLogging("mergeRequestFetcher(GetMergeRequest)", isDebugLogging, func() (*gitlab.MergeRequest, error) {
			mr, _, err := client.MergeRequests.GetMergeRequest(projectName, mrID, nil)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return mr, nil
		})
		if err != nil {
			return errors.WithStack(err)
		}

		description = mr.Description
		authorName = mr.Author.Name
		authorAvatarURL = mr.Author.AvatarURL
		footerTime = mr.CreatedAt

		matched2 := regexp.MustCompile(`#note_(\d+)$`).FindStringSubmatch(path)

		if matched2 != nil {
			note, err := util.WithDebugLogging("mergeRequestFetcher(GetMergeRequestNote)", isDebugLogging, func() (*gitlab.Note, error) {
				noteID, _ := strconv.ParseInt(matched2[1], 10, 64)
				note, _, err := client.Notes.GetMergeRequestNote(projectName, mrID, noteID)
				if err != nil {
					return nil, errors.WithStack(err)
				}
				return note, nil
			})
			if err != nil {
				return errors.WithStack(err)
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
		project, err = util.WithDebugLogging("mergeRequestFetcher(GetProject)", isDebugLogging, func() (*gitlab.Project, error) {
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
