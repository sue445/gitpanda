package gitlab

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/sue445/gitpanda/util"
	"gitlab.com/gitlab-org/api/client-go"
	"golang.org/x/sync/errgroup"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type snippetFetcher struct {
}

func (f *snippetFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(`^(?:-/)?snippets/(\d+)`).FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	var eg errgroup.Group

	var snippet *gitlab.Snippet
	authorName := ""
	var footerTime *time.Time

	snippetID, _ := strconv.ParseInt(matched[1], 10, 64)

	eg.Go(func() error {
		var err error
		snippet, err = util.WithDebugLogging("snippetFetcher(GetSnippet)", isDebugLogging, func() (*gitlab.Snippet, error) {
			snippet, _, err := client.Snippets.GetSnippet(snippetID)
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

		return nil
	})

	content := ""

	eg.Go(func() error {
		body, err := util.WithDebugLogging("snippetFetcher(SnippetContent)", isDebugLogging, func() (*string, error) {
			rawFile, _, err := client.Snippets.SnippetContent(snippetID)
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
