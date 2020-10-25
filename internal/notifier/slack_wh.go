package notifier

import (
	"context"

	"github.com/slack-go/slack"
)

func NewSlackWebHookNotifier(url string) *SlackWebHookNotifier {
	return &SlackWebHookNotifier{
		url: url,
	}
}

type SlackWebHookNotifier struct {
	url string
}

func (sn *SlackWebHookNotifier) Send(ctx context.Context, message string) error {
	msg := &slack.WebhookMessage{
		Text: message,
	}
	return slack.PostWebhookContext(ctx, sn.url, msg)
}
