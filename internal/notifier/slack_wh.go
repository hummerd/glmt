package notifier

import (
	"context"

	"github.com/slack-go/slack"
	"gitlab.com/gitlab-merge-tool/glmt/internal/templating"
)

func NewSlackWebHookNotifier(url, user, messageTemplate string) *SlackWebHookNotifier {
	return &SlackWebHookNotifier{
		url:             url,
		user:            user,
		messageTemplate: messageTemplate,
	}
}

type SlackWebHookNotifier struct {
	url             string
	user            string
	messageTemplate string
}

func (sn *SlackWebHookNotifier) Send(ctx context.Context, args map[string]string, add string) error {
	templ := sn.messageTemplate
	if templ == "" {
		templ = "<!here>\n{{.Description}}\n{{.MergeRequestURL}}"
	}

	m := templating.CreateText("slack_wh_message", templ, args)

	if add != "" {
		m += "\n" + add
	}

	msg := &slack.WebhookMessage{
		Text:     m,
		Username: sn.user,
	}
	return slack.PostWebhookContext(ctx, sn.url, msg)
}
