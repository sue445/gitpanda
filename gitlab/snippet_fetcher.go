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

type snippetFetcher struct {
}

func (f *snippetFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile("^(?:-/)?snippets/(\\d+)").FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	var eg errgroup.Group

	var snippet *gitlab.Snippet
	authorName := ""
	var footerTime *time.Time

	snippetID, _ := strconv.Atoi(matched[1])

	eg.Go(func() error {
		var err error
		start := time.Now()
		snippet, _, err = client.Snippets.GetSnippet(snippetID)

		if err != nil {
			return errors.WithStack(err)
		}

		if isDebugLogging {
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] snippetFetcher (%s): snippet=%+v\n", duration, snippet)
		}

		authorName = snippet.Author.Name
		footerTime = snippet.CreatedAt

		return nil
	})

	content := ""

	eg.Go(func() error {
		var err error
		start := time.Now()
		rawFile, _, err := client.Snippets.SnippetContent(snippetID)

		if err != nil {
			return errors.WithStack(err)
		}

		content = strings.TrimSpace(string(rawFile))

		if isDebugLogging {
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] snippetFetcher (%s): content=%+v\n", duration, content)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	page := &Page{
		Title:                  snippet.FileName,
		Description:            fmt.Sprintf("```\n%s\n```", content),
		AuthorName:             authorName,
		AuthorAvatarURL:        "",
		AvatarURL:              "",
		CanTruncateDescription: false,
		FooterTitle:            "",
		FooterURL:              "",
		FooterTime:             footerTime,
	}

	return page, nil
}
