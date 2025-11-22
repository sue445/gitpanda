package gitlab

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/sue445/gitpanda/util"
	"gitlab.com/gitlab-org/api/client-go"
	"golang.org/x/sync/errgroup"
)

type projectSnippetFetcher struct {
}

func (f *projectSnippetFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(reProjectName + `/snippets/(\d+)`).FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := sanitizeProjectName(matched[1])

	var eg errgroup.Group

	var snippet *gitlab.Snippet
	authorName := ""
	authorAvatarURL := ""
	var footerTime *time.Time
	var note *gitlab.Note

	snippetID, _ := strconv.Atoi(matched[2])

	eg.Go(func() error {
		var err error
		snippet, err = util.WithDebugLogging("projectSnippetFetcher(GetSnippet)", isDebugLogging, func() (*gitlab.Snippet, error) {
			snippet, _, err := client.ProjectSnippets.GetSnippet(projectName, snippetID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return snippet, nil
		})

		if err != nil {
			return errors.WithStack(err)
		}

		authorName = snippet.Author.Name
		footerTime = snippet.CreatedAt

		matched2 := regexp.MustCompile(`#note_(\d+)$`).FindStringSubmatch(path)

		if matched2 != nil {
			note, err = util.WithDebugLogging("projectSnippetFetcher(GetSnippetNote)", isDebugLogging, func() (*gitlab.Note, error) {
				noteID, _ := strconv.Atoi(matched2[1])
				note, _, err := client.Notes.GetSnippetNote(projectName, snippetID, noteID)
				if err != nil {
					return nil, errors.WithStack(err)
				}
				return note, nil
			})

			if err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})

	content := ""

	eg.Go(func() error {
		body, err := util.WithDebugLogging("projectSnippetFetcher(SnippetContent)", isDebugLogging, func() (*string, error) {
			rawFile, _, err := client.ProjectSnippets.SnippetContent(projectName, snippetID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			body := string(rawFile)
			return &body, nil
		})
		if err != nil {
			return errors.WithStack(err)
		}

		content = strings.TrimSpace(*body)

		return nil
	})

	var project *gitlab.Project
	eg.Go(func() error {
		var err error
		project, err = util.WithDebugLogging("projectSnippetFetcher(GetProject)", isDebugLogging, func() (*gitlab.Project, error) {
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

	description := fmt.Sprintf("```\n%s\n```", content)
	canTruncateDescription := false

	if note != nil {
		description = note.Body
		authorName = note.Author.Name
		authorAvatarURL = note.Author.AvatarURL
		footerTime = note.CreatedAt
		canTruncateDescription = true
	}

	page := &Page{
		Title:                  snippet.FileName,
		Description:            description,
		AuthorName:             authorName,
		AuthorAvatarURL:        authorAvatarURL,
		AvatarURL:              project.AvatarURL,
		CanTruncateDescription: canTruncateDescription,
		FooterTitle:            project.PathWithNamespace,
		FooterURL:              project.WebURL,
		FooterTime:             footerTime,
	}

	return page, nil
}
