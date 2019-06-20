package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"regexp"
	"strconv"
	"strings"
)

type mergeRequestFetcher struct {
}

func (f *mergeRequestFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	re := regexp.MustCompile("^([^/]+)/([^/]+)/merge_requests/(\\d+)")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := matched[1] + "/" + matched[2]
	mrID, _ := strconv.Atoi(matched[3])
	mr, _, err := client.MergeRequests.GetMergeRequest(projectName, mrID, nil)

	if err != nil {
		return nil, err
	}

	if isDebugLogging {
		fmt.Printf("[DEBUG] fetchMergeRequestURL: mr=%+v\n", mr)
	}

	project, _, err := client.Projects.GetProject(projectName, nil)

	if err != nil {
		return nil, err
	}

	if isDebugLogging {
		fmt.Printf("[DEBUG] fetchMergeRequestURL: project=%+v\n", project)
	}

	description := mr.Description
	authorName := mr.Author.Name
	authorAvatarURL := mr.Author.AvatarURL
	footerTime := mr.CreatedAt

	re2 := regexp.MustCompile("#note_(\\d+)$")
	matched2 := re2.FindStringSubmatch(path)

	if matched2 != nil {
		noteID, _ := strconv.Atoi(matched2[1])
		note, _, err := client.Notes.GetMergeRequestNote(projectName, mrID, noteID)

		if err != nil {
			return nil, err
		}

		if isDebugLogging {
			fmt.Printf("[DEBUG] fetchMergeRequestURL: note=%+v\n", note)
		}

		description = note.Body
		authorName = note.Author.Name
		authorAvatarURL = note.Author.AvatarURL
		footerTime = note.CreatedAt
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
