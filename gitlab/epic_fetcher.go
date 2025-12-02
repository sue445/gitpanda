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

type epicFetcher struct {
}

func (f *epicFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(`^groups/(.+?)/epics/(\d+)`).FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	groupName := sanitizeProjectName(matched[1])

	var eg errgroup.Group

	var epic *gitlab.Epic
	description := ""
	authorName := ""
	authorAvatarURL := ""
	var footerTime *time.Time

	eg.Go(func() error {
		var err error
		epic, err = util.WithDebugLogging("epicFetcher(GetEpic)", isDebugLogging, func() (*gitlab.Epic, error) {
			epicIID, _ := strconv.ParseInt(matched[2], 10, 64)
			epic, _, err := client.Epics.GetEpic(groupName, epicIID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return epic, nil
		})
		if err != nil {
			return errors.WithStack(err)
		}

		description = epic.Description
		authorName = epic.Author.Name
		authorAvatarURL = epic.Author.AvatarURL
		footerTime = epic.CreatedAt

		matched2 := regexp.MustCompile(`#note_(\d+)$`).FindStringSubmatch(path)

		if matched2 != nil {
			noteID, _ := strconv.ParseInt(matched2[1], 10, 64)

			note, err := util.WithDebugLogging("noteFetcher(GetEpicNote)", isDebugLogging, func() (*gitlab.Note, error) {
				note, _, err := client.Notes.GetEpicNote(groupName, epic.ID, noteID)
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

	var group *gitlab.Group
	eg.Go(func() error {
		var err error
		group, err = util.WithDebugLogging("groupFetcher(GetGroup)", isDebugLogging, func() (*gitlab.Group, error) {
			// FIXME: `WithProjects` is deprecated and will be removed since API v5.
			//        However `with_projects` is default to `true`.
			//        c.f. https://docs.gitlab.com/api/groups/#get-a-single-group
			//        So `with_projects=false` is needed to get the group without including Projects.
			//        `WithProjects` can be deleted after API v5.
			group, _, err := client.Groups.GetGroup(groupName, &gitlab.GetGroupOptions{WithProjects: gitlab.Ptr(false)})
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return group, nil
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
		Title:                  strings.Join([]string{epic.Title, "Epics", group.Name, "GitLab"}, titleSeparator),
		Description:            description,
		AuthorName:             authorName,
		AuthorAvatarURL:        authorAvatarURL,
		AvatarURL:              group.AvatarURL,
		CanTruncateDescription: true,
		FooterTitle:            group.Path,
		FooterURL:              group.WebURL,
		FooterTime:             footerTime,
	}

	return page, nil
}
