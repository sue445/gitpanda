package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type userOrGroupFetcher struct {
	baseURL string
}

func (f *userOrGroupFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	re := regexp.MustCompile("^([^/]+)/?$")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	name := matched[1]
	userPage, err := f.fetchUserPath(name, client, isDebugLogging)

	if err != nil {
		return nil, err
	}

	if userPage != nil {
		return userPage, nil
	}

	groupPage, err := f.fetchGroupPath(name, client, isDebugLogging)

	if err != nil {
		return nil, err
	}

	if groupPage != nil {
		return groupPage, nil
	}

	return nil, fmt.Errorf("%s is not found", name)
}

func (f *userOrGroupFetcher) fetchUserPath(name string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	start := time.Now()
	users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{Username: &name})

	if err != nil {
		return nil, err
	}

	if isDebugLogging {
		duration := time.Now().Sub(start)
		fmt.Printf("[DEBUG] fetchUserPath (%s): users=%+v\n", duration, users)
	}

	if len(users) < 1 {
		return nil, nil
	}

	user := users[0]

	page := &Page{
		Title:                  strings.Join([]string{user.Name, "GitLab"}, titleSeparator),
		Description:            user.Name,
		AuthorName:             user.Name,
		AuthorAvatarURL:        user.AvatarURL,
		AvatarURL:              user.AvatarURL,
		CanTruncateDescription: true,
		FooterTitle:            fmt.Sprintf("@%s", user.Username),
		FooterURL:              fmt.Sprintf("%s/%s", f.baseURL, user.Username),
		FooterTime:             user.CreatedAt,
	}

	return page, nil
}

func (f *userOrGroupFetcher) fetchGroupPath(name string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	start := time.Now()
	group, _, err := client.Groups.GetGroup(name, func(req *http.Request) error {
		req.URL.RawQuery = "with_projects=false"
		return nil
	})

	if err != nil {
		return nil, err
	}

	if isDebugLogging {
		duration := time.Now().Sub(start)
		fmt.Printf("[DEBUG] fetchGroupPath (%s): group=%+v\n", duration, group)
	}

	page := &Page{
		Title:                  strings.Join([]string{group.Name, "GitLab"}, titleSeparator),
		Description:            group.Description,
		AuthorName:             "",
		AuthorAvatarURL:        "",
		AvatarURL:              group.AvatarURL,
		CanTruncateDescription: true,
		FooterTitle:            fmt.Sprintf("@%s", group.Path),
		FooterURL:              group.WebURL,
		FooterTime:             nil,
	}

	return page, nil
}
