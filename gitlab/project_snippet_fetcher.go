package gitlab

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type projectSnippetFetcher struct {
}

func (f *projectSnippetFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(reProjectName + "/snippets/(\\d+)").FindStringSubmatch(path)

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
		start := time.Now()
		snippet, _, err = client.ProjectSnippets.GetSnippet(projectName, snippetID)

		if err != nil {
			return errors.WithStack(err)
		}

		if isDebugLogging {
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] projectSnippetFetcher (%s): snippet=%+v\n", duration, snippet)
		}

		authorName = snippet.Author.Name
		footerTime = snippet.CreatedAt

		matched2 := regexp.MustCompile("#note_(\\d+)$").FindStringSubmatch(path)

		if matched2 != nil {
			noteID, _ := strconv.Atoi(matched2[1])
			start := time.Now()
			note, _, err = client.Notes.GetSnippetNote(projectName, snippetID, noteID)

			if err != nil {
				return errors.WithStack(err)
			}

			if isDebugLogging {
				duration := time.Now().Sub(start)
				fmt.Printf("[DEBUG] projectSnippetFetcher (%s): note=%+v\n", duration, note)
			}
		}

		return nil
	})

	content := ""

	eg.Go(func() error {
		var err error
		start := time.Now()
		rawFile, _, err := client.ProjectSnippets.SnippetContent(projectName, snippetID)

		if err != nil {
			return errors.WithStack(err)
		}

		content = strings.TrimSpace(string(rawFile))

		if isDebugLogging {
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] projectSnippetFetcher (%s): content=%+v\n", duration, content)
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
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] projectSnippetFetcher (%s): project=%+v\n", duration, project)
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
