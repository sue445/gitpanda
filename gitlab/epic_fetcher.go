package gitlab

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"gitlab.com/gitlab-org/api/client-go"
	"golang.org/x/sync/errgroup"
	"regexp"
	"strconv"
	"strings"
	"time"
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
		epicID, _ := strconv.Atoi(matched[2])
		start := time.Now()
		epic, _, err = client.Epics.GetEpic(groupName, epicID)

		if err != nil {
			return errors.WithStack(err)
		}

		if isDebugLogging {
			duration := time.Since(start)
			fmt.Printf("[DEBUG] epicFetcher (%s): epic=%+v\n", duration, epic)
		}

		description = epic.Description
		authorName = epic.Author.Name
		authorAvatarURL = epic.Author.AvatarURL
		footerTime = epic.CreatedAt

		matched2 := regexp.MustCompile(`#note_(\d+)$`).FindStringSubmatch(path)

		if matched2 != nil {
			noteID, _ := strconv.Atoi(matched2[1])
			start := time.Now()
			note, _, err := client.Notes.GetEpicNote(groupName, epicID, noteID)

			if err != nil {
				return errors.WithStack(err)
			}

			if isDebugLogging {
				duration := time.Since(start)
				fmt.Printf("[DEBUG] epicFetcher (%s): note=%+v\n", duration, note)
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
		start := time.Now()
		group, _, err = client.Groups.GetGroup(groupName, &gitlab.GetGroupOptions{WithProjects: gitlab.Ptr(false)})

		if err != nil {
			return errors.WithStack(err)
		}

		if isDebugLogging {
			duration := time.Since(start)
			fmt.Printf("[DEBUG] epicFetcher (%s): group=%+v\n", duration, group)
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
