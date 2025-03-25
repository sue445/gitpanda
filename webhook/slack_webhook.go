package webhook

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/sue445/gitpanda/gitlab"
	"github.com/sue445/gitpanda/util"
	"golang.org/x/sync/errgroup"
)

// SlackWebhook represents Slack webhook
type SlackWebhook struct {
	slackOAuthAccessToken  string
	slackVerificationToken string
	gitLabURLParserParams  *gitlab.URLParserParams
}

// NewSlackWebhook create new SlackWebhook instance
func NewSlackWebhook(slackOAuthAccessToken string, slackVerificationToken string, gitLabURLParserParams *gitlab.URLParserParams) *SlackWebhook {
	return &SlackWebhook{slackOAuthAccessToken: slackOAuthAccessToken, slackVerificationToken: slackVerificationToken, gitLabURLParserParams: gitLabURLParserParams}
}

// Request handles Slack event
func (s *SlackWebhook) Request(body string, truncateLines int) (string, error) {
	option := slackevents.OptionNoVerifyToken()

	if s.slackVerificationToken != "" {
		option = slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: s.slackVerificationToken})
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), option)

	if err != nil {
		return "Failed: slackevents.ParseEvent", errors.WithStack(err)
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			return "Failed: json.Unmarshal", errors.WithStack(err)
		}
		return r.Challenge, nil

	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.LinkSharedEvent:
			return s.requestLinkSharedEvent(ev, truncateLines)
		}
	}

	return "", fmt.Errorf("unknown event type: %s", eventsAPIEvent.Type)
}

func (s *SlackWebhook) requestLinkSharedEvent(ev *slackevents.LinkSharedEvent, truncateLines int) (string, error) {
	p, err := gitlab.NewGitlabURLParser(s.gitLabURLParserParams)

	if err != nil {
		return "Failed: NewGitlabURLParser", errors.WithStack(err)
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
				return errors.WithStack(err)
			}

			if page == nil {
				return nil
			}

			if s.gitLabURLParserParams.IsDebugLogging {
				duration := time.Since(start)
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
		if len(unfurls) > 0 {
			// NOTE: Don't returns error when contains 1 or more valid urls
			fmt.Printf("[WARN] FetchURL error=%+v\n", err)
		} else {
			return "Failed: FetchURL", errors.WithStack(err)
		}
	}

	if len(unfurls) == 0 {
		return "do nothing", nil
	}

	start := time.Now()
	api := slack.New(s.slackOAuthAccessToken)
	_, _, _, err = api.UnfurlMessage(ev.Channel, ev.MessageTimeStamp, unfurls)

	if err != nil {
		return "Failed: UnfurlMessage", errors.WithStack(err)
	}

	if s.gitLabURLParserParams.IsDebugLogging {
		duration := time.Since(start)
		fmt.Printf("[DEBUG] UnfurlMessage (%s)\n", duration)
	}

	return "ok", nil
}
