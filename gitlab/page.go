package gitlab

// Page represents info of GitLab page
type Page struct {
	Title                  string
	Description            string
	AuthorName             string
	AuthorAvatarURL        string
	AvatarURL              string
	CanTruncateDescription bool
}
