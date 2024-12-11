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

type issueFetcher struct {
}

func (f *issueFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(reProjectName + "/issues/(\\d+)").FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := sanitizeProjectName(matched[1])

	var eg errgroup.Group

	var issue *gitlab.Issue
	description := ""
	authorName := ""
	authorAvatarURL := ""
	var footerTime *time.Time

	eg.Go(func() error {
		var err error
		issueID, _ := strconv.Atoi(matched[2])
		start := time.Now()
		issue, _, err = client.Issues.GetIssue(projectName, issueID)

		if err != nil {
			return errors.WithStack(err)
		}

		if isDebugLogging {
			duration := time.Since(start)
			fmt.Printf("[DEBUG] issueFetcher (%s): issue=%+v\n", duration, issue)
		}

		description = issue.Description
		authorName = issue.Author.Name
		authorAvatarURL = issue.Author.AvatarURL
		footerTime = issue.CreatedAt

		matched2 := regexp.MustCompile(`#note_(\d+)$`).FindStringSubmatch(path)

		if matched2 != nil {
			noteID, _ := strconv.Atoi(matched2[1])
			start := time.Now()
			note, _, err := client.Notes.GetIssueNote(projectName, issueID, noteID)

			if err != nil {
				return errors.WithStack(err)
			}

			if isDebugLogging {
				duration := time.Since(start)
				fmt.Printf("[DEBUG] issueFetcher (%s): note=%+v\n", duration, note)
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
			fmt.Printf("[DEBUG] issueFetcher (%s): project=%+v\n", duration, project)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	page := &Page{
		Title:                  strings.Join([]string{issue.Title, "Issues", project.NameWithNamespace, "GitLab"}, titleSeparator),
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
