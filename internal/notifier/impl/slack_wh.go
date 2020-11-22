package impl

import (
	"context"

	"gitlab.com/gitlab-merge-tool/glmt/internal/config"
	"gitlab.com/gitlab-merge-tool/glmt/internal/glmt"
	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
	"gitlab.com/gitlab-merge-tool/glmt/internal/templating"

	"github.com/slack-go/slack"
)

func NewSlackWebHookNotifier(cfg config.SlackWebHook) *SlackWebHookNotifier {
	return &SlackWebHookNotifier{
		url:         cfg.URL,
		user:        cfg.User,
		messageTmpl: cfg.MessageTmpl,
	}
}

type SlackWebHookNotifier struct {
	url         string
	user        string
	messageTmpl string
}

func (sn *SlackWebHookNotifier) Send(ctx context.Context, args map[string]string, add string, mentions []*team.Member) error {
	templ := sn.messageTmpl
	if templ == "" {
		templ = "<!here>\n{{.Description}}\n{{.MergeRequestURL}}"
	}

	args[glmt.TmpVarNotificationMentions] = getMentions(
		mentions,
		memberKeySlack,
		"<@%s>",
	)

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
