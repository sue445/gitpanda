package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"regexp"
	"strconv"
	"strings"
)

type issueFetcher struct {
}

func (f *issueFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	re := regexp.MustCompile("^([^/]+)/([^/]+)/issues/(\\d+)")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := matched[1] + "/" + matched[2]
	issueID, _ := strconv.Atoi(matched[3])
	issue, _, err := client.Issues.GetIssue(projectName, issueID)

	if err != nil {
		return nil, err
	}

	if isDebugLogging {
		fmt.Printf("[DEBUG] fetchIssueURL: issue=%+v\n", issue)
	}

	project, _, err := client.Projects.GetProject(projectName, nil)

	if err != nil {
		return nil, err
	}

	if isDebugLogging {
		fmt.Printf("[DEBUG] fetchIssueURL: project=%+v\n", project)
	}

	description := issue.Description
	authorName := issue.Author.Name
	authorAvatarURL := issue.Author.AvatarURL
	footerTime := issue.CreatedAt

	re2 := regexp.MustCompile("#note_(\\d+)$")
	matched2 := re2.FindStringSubmatch(path)

	if matched2 != nil {
		noteID, _ := strconv.Atoi(matched2[1])
		note, _, err := client.Notes.GetIssueNote(projectName, issueID, noteID)

		if err != nil {
			return nil, err
		}

		if isDebugLogging {
			fmt.Printf("[DEBUG] fetchIssueURL: note=%+v\n", note)
		}

		description = note.Body
		authorName = note.Author.Name
		authorAvatarURL = note.Author.AvatarURL
		footerTime = note.CreatedAt
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
