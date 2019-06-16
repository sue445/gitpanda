package gitlab

// GitLabPage represents info of GitLab page
type GitLabPage struct {
	Title                  string
	Description            string
	AuthorName             string
	AuthorAvatarURL        string
	AvatarURL              string
	CanTruncateDescription bool
}
