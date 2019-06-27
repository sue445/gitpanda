package gitlab

import (
	"fmt"
	"github.com/sue445/go-gitlab"
	"regexp"
	"strings"
	"time"
)

type userFetcher struct {
	baseURL string
}

func (f *userFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	re := regexp.MustCompile("^([^/]+)/?$")
	matched := re.FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	username := matched[1]
	start := time.Now()
	users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{Username: &username})

	if err != nil {
		return nil, err
	}

	if isDebugLogging {
		duration := time.Now().Sub(start)
		fmt.Printf("[DEBUG] projectFetcher (%s): users=%+v\n", duration, users)
	}

	if len(users) < 1 {
		return nil, fmt.Errorf("NotFound user: %s", username)
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
