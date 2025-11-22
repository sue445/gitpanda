package gitlab

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sue445/gitpanda/util"
	"gitlab.com/gitlab-org/api/client-go"
)

type userOrGroupFetcher struct {
}

func (f *userOrGroupFetcher) fetchPath(path string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	matched := regexp.MustCompile(`^([^/]+)/?$`).FindStringSubmatch(path)

	if matched == nil {
		return nil, nil
	}

	name := matched[1]
	userPage, err := f.fetchUserPath(name, client, isDebugLogging)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if userPage != nil {
		return userPage, nil
	}

	groupPage, err := f.fetchGroupPath(name, client, isDebugLogging)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if groupPage != nil {
		return groupPage, nil
	}

	return nil, fmt.Errorf("%s is not found", name)
}

func (f *userOrGroupFetcher) fetchUserPath(name string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	user, err := util.WithDebugLogging("userOrGroupFetcher(fetchUserPath)", isDebugLogging, func() (*gitlab.User, error) {
		users, _, err := client.Users.ListUsers(&gitlab.ListUsersOptions{Username: &name})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if len(users) < 1 {
			return nil, nil
		}
		return users[0], nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if user == nil {
		return nil, nil
	}

	page := &Page{
		Title:                  strings.Join([]string{user.Name, "GitLab"}, titleSeparator),
		Description:            user.Name,
		AuthorName:             user.Name,
		AuthorAvatarURL:        user.AvatarURL,
		AvatarURL:              user.AvatarURL,
		CanTruncateDescription: true,
		FooterTitle:            fmt.Sprintf("@%s", user.Username),
		FooterURL:              user.WebURL,
		FooterTime:             user.CreatedAt,
	}

	return page, nil
}

func (f *userOrGroupFetcher) fetchGroupPath(name string, client *gitlab.Client, isDebugLogging bool) (*Page, error) {
	group, err := util.WithDebugLogging("userOrGroupFetcher(fetchGroupPath)", isDebugLogging, func() (*gitlab.Group, error) {
		// FIXME: `WithProjects` is deprecated and will be removed since API v5.
		//        However `with_projects` is default to `true`.
		//        c.f. https://docs.gitlab.com/api/groups/#get-a-single-group
		//        So `with_projects=false` is needed to get the group without including Projects.
		//        `WithProjects` can be deleted after API v5.
		group, _, err := client.Groups.GetGroup(name, &gitlab.GetGroupOptions{WithProjects: gitlab.Ptr(false)})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return group, nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
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
