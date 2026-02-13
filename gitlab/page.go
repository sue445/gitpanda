package gitlab

import (
	"time"
)

// Page represents info of GitLab page
type Page struct {
	Title                  string
	Description            string
	AuthorName             string
	AuthorAvatarURL        string
	AvatarURL              string
	CanTruncateDescription bool
	FooterTitle            string
	FooterURL              string
	FooterTime             *time.Time
	Color                  string
}
