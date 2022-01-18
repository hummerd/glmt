package impl

import (
	"context"
	"strings"

	"gitlab.com/gitlab-merge-tool/glmt/internal/config"
	"gitlab.com/gitlab-merge-tool/glmt/internal/glmt"
	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
	"gitlab.com/gitlab-merge-tool/glmt/internal/templating"

	"github.com/nafisfaysal/matterhook"
)

const (
	mattermostMemberKey = "mattermost_member_id"

	mattermostDefaultMessageTmpl = "@here\n{{.Description}}\n{{.MergeRequestURL}}"
)

func NewMattermostWebHookNotifier(cfg config.MattermostWebHook) *MattermostWebHookNotifier {
	return &MattermostWebHookNotifier{
		url:         cfg.URL,
		user:        cfg.User,
		messageTpml: cfg.MessageTmpl,
	}
}

type MattermostWebHookNotifier struct {
	url         string
	user        string
	messageTpml string
}

func (mn *MattermostWebHookNotifier) Send(ctx context.Context, args map[string]string, add string, mentions []*team.Member) error {
	templ := mn.messageTpml
	if templ == "" {
		templ = mattermostDefaultMessageTmpl
	}

	args[glmt.TmpVarNotificationMentions] = getMattermostMentions(mentions)

	m := templating.CreateText("mattermost_wh_message", templ, args)

	if add != "" {
		m += "\n" + add
	}

	message := matterhook.Message{
		Text:     m,
		Username: mn.user,
	}
	return matterhook.Send(mn.url, message)
}

func getMattermostMentions(mentions []*team.Member) string {
	ms := make([]string, 0, len(mentions))
	for _, m := range mentions {
		if n := m.Names[mattermostMemberKey]; n != "" {
			ms = append(ms, "@"+n)
		}
	}

	return strings.Join(ms, ", ")
}
