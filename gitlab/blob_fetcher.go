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
	"unicode/utf8"
)

type blobFetcher struct {
}

func (f *blobFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(reProjectName + "/blob/([^/]+)/(.+)$").FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := sanitizeProjectName(matched[1])
	sha1 := matched[2]
	fileName := matched[3]

	lineMatched := regexp.MustCompile("(.+)#L([0-9-]+)$").FindStringSubmatch(fileName)

	if lineMatched != nil {
		fileName = lineMatched[1]
	}

	paramsMatched := regexp.MustCompile(`(.+)\?(.+)$`).FindStringSubmatch(fileName)

	if paramsMatched != nil {
		fileName = paramsMatched[1]
	}

	var eg errgroup.Group

	selectedFile := ""
	lineRange := ""
	eg.Go(func() error {
		fileBody, err := util.WithDebugLogging("blobFetcher(GetRawFile)", isDebugLogging, func() (*string, error) {
			rawFile, _, err := client.RepositoryFiles.GetRawFile(projectName, fileName, &gitlab.GetRawFileOptions{Ref: &sha1})
			if err != nil {
				return nil, errors.WithStack(err)
			}
			fileBody := string(rawFile)
			return &fileBody, nil
		})
		if err != nil {
			return errors.WithStack(err)
		}

		if !utf8.ValidString(*fileBody) {
			return nil
		}

		if lineMatched == nil {
			selectedFile = *fileBody
			return nil
		}

		lineHash := lineMatched[2]
		lines := strings.Split(lineHash, "-")

		switch len(lines) {
		case 1:
			line, _ := strconv.Atoi(lines[0])
			lineRange = lines[0]
			selectedFile = util.SelectLine(*fileBody, line)
			return nil
		case 2:
			startLine, _ := strconv.Atoi(lines[0])
			endLine, _ := strconv.Atoi(lines[1])
			lineRange = fmt.Sprintf("%s-%s", lines[0], lines[1])
			selectedFile = util.SelectLines(*fileBody, startLine, endLine)
			return nil
		default:
			return fmt.Errorf("invalid line: L%s", lineHash)
		}
	})

	var project *gitlab.Project
	eg.Go(func() error {
		var err error
		project, err = util.WithDebugLogging("blobFetcher(GetProject)", isDebugLogging, func() (*gitlab.Project, error) {
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

	title := fileName
	if lineRange != "" {
		title = fmt.Sprintf("%s:%s", title, lineRange)
	}

	description := ""
	if selectedFile != "" {
		description = fmt.Sprintf("```\n%s\n```", selectedFile)
	}

	page := &Page{
		Title:                  title,
		Description:            description,
		AuthorName:             "",
		AuthorAvatarURL:        "",
		AvatarURL:              project.AvatarURL,
		CanTruncateDescription: false,
		FooterTitle:            project.PathWithNamespace,
		FooterURL:              project.WebURL,
		FooterTime:             nil,
	}

	return page, nil
}
