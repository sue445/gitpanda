package gitlab

import (
	"fmt"
	"github.com/sue445/gitpanda/util"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type blobFetcher struct {
}

func (f *blobFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	re := regexp.MustCompile("^([^/]+)/([^/]+)/blob/([^/]+)/(.+)#L([0-9-]+)$")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	projectName := matched[1] + "/" + matched[2]
	sha1 := matched[3]
	fileName := matched[4]

	var eg errgroup.Group

	selectedFile := ""
	lineRange := ""
	eg.Go(func() error {
		start := time.Now()
		rawFile, _, err := client.RepositoryFiles.GetRawFile(projectName, fileName, &gitlab.GetRawFileOptions{Ref: &sha1})

		if err != nil {
			return err
		}

		fileBody := string(rawFile)

		if isDebugLogging {
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] blobFetcher (%s): fileBody=%s\n", duration, fileBody)
		}

		lineHash := matched[5]
		lines := strings.Split(lineHash, "-")

		switch len(lines) {
		case 1:
			line, _ := strconv.Atoi(lines[0])
			lineRange = lines[0]
			selectedFile = util.SelectLine(fileBody, line)
			return nil
		case 2:
			startLine, _ := strconv.Atoi(lines[0])
			endLine, _ := strconv.Atoi(lines[1])
			lineRange = fmt.Sprintf("%s-%s", lines[0], lines[1])
			selectedFile = util.SelectLines(fileBody, startLine, endLine)
			return nil
		default:
			return fmt.Errorf("Invalid line: L%s", lineHash)
		}
	})

	var project *gitlab.Project
	eg.Go(func() error {
		start := time.Now()
		_project, _, err := client.Projects.GetProject(projectName, nil)

		if err != nil {
			return err
		}

		project = _project

		if isDebugLogging {
			duration := time.Now().Sub(start)
			fmt.Printf("[DEBUG] blobFetcher (%s): project=%+v\n", duration, project)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	page := &Page{
		Title:                  fmt.Sprintf("%s:%s", fileName, lineRange),
		Description:            fmt.Sprintf("```\n%s\n```", selectedFile),
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
