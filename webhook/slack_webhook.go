package webhook

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/sue445/gitpanda/gitlab"
	"github.com/sue445/gitpanda/util"
	"golang.org/x/sync/errgroup"
	"strconv"
	"sync"
	"time"
)

// SlackWebhook represents Slack webhook
type SlackWebhook struct {
	slackOAuthAccessToken string
	gitLabURLParserParams *gitlab.URLParserParams
}

// NewSlackWebhook create new SlackWebhook instance
func NewSlackWebhook(slackOAuthAccessToken string, gitLabURLParserParams *gitlab.URLParserParams) *SlackWebhook {
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
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.LinkSharedEvent:
			return s.requestLinkSharedEvent(ev, truncateLines)
		}
	}

	return "", fmt.Errorf("Unknown event type: %s", eventsAPIEvent.Type)
}

func (s *SlackWebhook) requestLinkSharedEvent(ev *slackevents.LinkSharedEvent, truncateLines int) (string, error) {
	p, err := gitlab.NewGitlabURLParser(s.gitLabURLParserParams)

	if err != nil {
		return "Failed: NewGitlabURLParser", err
	}

	unfurls := map[string]slack.Attachment{}

	var mu sync.Mutex
	var eg errgroup.Group
	for _, link := range ev.Links {
		url := link.URL
		eg.Go(func() error {
			start := time.Now()
			page, err := p.FetchURL(url)

			if err != nil {
				return err
			}

			if page == nil {
				return nil
			}

			if s.gitLabURLParserParams.IsDebugLogging {
				duration := time.Now().Sub(start)
				fmt.Printf("[DEBUG] FetchURL (%s): page=%v\n", duration, page)
			}

			description := util.FormatMarkdownForSlack(page.Description)

			if page.CanTruncateDescription {
				description = util.TruncateWithLine(description, truncateLines)
			}

			attachment := slack.Attachment{
				Title:      page.Title,
				TitleLink:  url,
				AuthorName: page.AuthorName,
				AuthorIcon: page.AuthorAvatarURL,
				Text:       description,
				Footer:     page.FormatFooter(),
				ThumbURL:   page.AvatarURL,
			}

			if page.FooterTime != nil {
				attachment.Ts = json.Number(strconv.FormatInt(page.FooterTime.Unix(), 10))
			}

			if page.Color == "" {
				attachment.Color = gitlab.BrandColor
			} else {
				attachment.Color = page.Color
			}

			mu.Lock()
			defer mu.Unlock()
			unfurls[url] = attachment

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return "Failed: FetchURL", err
	}

	if len(unfurls) == 0 {
		return "do nothing", nil
	}

	start := time.Now()
	api := slack.New(s.slackOAuthAccessToken)
	_, _, _, err = api.UnfurlMessage(ev.Channel, ev.MessageTimeStamp.String(), unfurls)

	if err != nil {
		return "Failed: UnfurlMessage", err
	}

	if s.gitLabURLParserParams.IsDebugLogging {
		duration := time.Now().Sub(start)
		fmt.Printf("[DEBUG] UnfurlMessage (%s)\n", duration)
	}

	return "ok", nil
}
