package notifier

import (
	"context"
	"strings"

	"github.com/slack-go/slack"
	"gitlab.com/gitlab-merge-tool/glmt/internal/glmt"
	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
	"gitlab.com/gitlab-merge-tool/glmt/internal/templating"
)

const (
	slackMemberKey = "slack_member_id"
)

func NewSlackWebHookNotifier(url, user, message string) *SlackWebHookNotifier {
	return &SlackWebHookNotifier{
		url:     url,
		user:    user,
		message: message,
	}
}

type SlackWebHookNotifier struct {
	url     string
	user    string
	message string
}

func (sn *SlackWebHookNotifier) Send(ctx context.Context, args map[string]string, add string, mentions []*team.Member) error {
	templ := sn.message
	if templ == "" {
		templ = "<!here>\n{{.Description}}\n{{.MergeRequestURL}}"
	}

	args[glmt.TmpVarNotificationMentions] = getSlackMentions(mentions)

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

func getSlackMentions(mentions []*team.Member) string {
	ms := make([]string, 0, len(mentions))
	for _, m := range mentions {
		n := m.Names[slackMemberKey]
		if n != "" {
			ms = append(ms, "<@"+n+">")
		}
	}

	return strings.Join(ms, ", ")
}
