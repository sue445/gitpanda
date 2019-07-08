package gitlab

import (
	"fmt"
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

// FormatFooter returns formatted footer for slack
func (p *Page) FormatFooter() string {
	if p.FooterURL != "" {
		if p.FooterTitle != "" {
			return fmt.Sprintf("<%s|%s>", p.FooterURL, p.FooterTitle)
		}
		return p.FooterURL
	}

	if p.FooterTitle != "" {
		return p.FooterTitle
	}
	return ""
}
