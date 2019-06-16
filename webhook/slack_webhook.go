package webhook

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/sue445/gitpanda"
	"github.com/sue445/gitpanda/gitlab"
	"github.com/sue445/gitpanda/util"
)

// SlackWebhook represents Slack webhook
type SlackWebhook struct {
	slackOAuthAccessToken string
	gitLabURLParserParams *gitlab.GitLabURLParserParams
}

// NewSlackWebhook create new SlackWebhook instance
func NewSlackWebhook(slackOAuthAccessToken string, gitLabURLParserParams *gitlab.GitLabURLParserParams) *SlackWebhook {
	return &SlackWebhook{slackOAuthAccessToken: slackOAuthAccessToken, gitLabURLParserParams: gitLabURLParserParams}
}

// Request handles Slack event
func (s *SlackWebhook) Request(body string, truncateLines int) (string, error) {
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())

	if err != nil {
		return "Failed: slackevents.ParseEvent", err
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			return "Failed: json.Unmarshal", err
		}
		return r.Challenge, nil

	case slackevents.CallbackEvent:
		p, err := gitlab.NewGitlabURLParser(s.gitLabURLParserParams)

		if err != nil {
			return "Failed: NewGitlabURLParser", err
		}

		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.LinkSharedEvent:
			unfurls := map[string]slack.Attachment{}

			for _, link := range ev.Links {
				page, err := p.FetchURL(link.URL)

				if err != nil {
					return "Failed: FetchURL", err
				}

				if page == nil {
					continue
				}

				if main.isDebugLogging {
					fmt.Printf("[DEBUG] FetchURL: page=%v\n", page)
				}

				description := page.Description
				if page.canTruncateDescription {
					description = util.TruncateWithLine(description, truncateLines)
				}

				attachment := slack.Attachment{
					Title:      page.Title,
					TitleLink:  link.URL,
					AuthorName: page.AuthorName,
					AuthorIcon: page.AuthorAvatarURL,
					Text:       description,
					Color:      "#fc6d26", // c.f. https://brandcolors.net/b/gitlab
				}

				if page.AvatarURL != "" {
					attachment.ThumbURL = page.AvatarURL
				}

				unfurls[link.URL] = attachment
			}

			if len(unfurls) == 0 {
				return "do nothing", nil
			}

			api := slack.New(s.slackOAuthAccessToken)
			_, _, _, err := api.UnfurlMessage(ev.Channel, ev.MessageTimeStamp.String(), unfurls)

			if err != nil {
				return "Failed: UnfurlMessage", err
			}

			return "ok", nil
		}
	}

	return "", fmt.Errorf("Unknown event type: %s", eventsAPIEvent.Type)
}
