package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

// SlackWebhook represents Slack webhook
type SlackWebhook struct {
	slackOAuthAccessToken string
	gitLabURLParserParams *GitLabURLParserParams
}

// NewSlackWebhook create new SlackWebhook instance
func NewSlackWebhook(slackOAuthAccessToken string, gitLabURLParserParams *GitLabURLParserParams) *SlackWebhook {
	return &SlackWebhook{slackOAuthAccessToken: slackOAuthAccessToken, gitLabURLParserParams: gitLabURLParserParams}
}

// Request handles Slack event
func (s *SlackWebhook) Request(body string) (string, error) {
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
		p, err := NewGitlabURLParser(s.gitLabURLParserParams)

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

				if isDebugLogging {
					fmt.Printf("[DEBUG] FetchURL: page=%v\n", page)
				}

				attachment := slack.Attachment{
					Title:      page.Title,
					TitleLink:  link.URL,
					AuthorName: page.AuthorName,
					AuthorIcon: page.AuthorAvatarURL,
					Text:       page.Description,
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
