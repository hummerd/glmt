// Package impl implements notifier.Notifier
package impl

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/gitlab-merge-tool/glmt/internal/notifier"
	"gitlab.com/gitlab-merge-tool/glmt/internal/team"

	"github.com/rs/zerolog/log"
)

const (
	memberKeySlack    = "slack_member_id"
	memberKeyTelegram = "telegram_member_id"
)

type MultiNotifier struct {
	notifiers []notifier.Notifier
}

func NewMultiNotifier(notifiers ...notifier.Notifier) *MultiNotifier {
	return &MultiNotifier{
		notifiers: notifiers,
	}
}

func (mn *MultiNotifier) Send(
	ctx context.Context,
	args map[string]string,
	add string,
	mentions []*team.Member,
) (err error) {
	for _, n := range mn.notifiers {
		err = n.Send(ctx, args, add, mentions)
		if err != nil {
			return fmt.Errorf("multi send: %T: %w", n, err)
		}
		log.Ctx(ctx).Debug().
			Str("add", add).
			Interface("args", args).
			Interface("mentions", mentions).
			Msgf("sent notification with %T", n)
	}

	return nil
}

func getMentions(mentions []*team.Member, memberKey string, format string) string {
	ms := make([]string, 0, len(mentions))
	for _, m := range mentions {
		n := m.Names[memberKey]
		if n != "" {
			ms = append(ms, fmt.Sprintf(format, n))
		}
	}

	return strings.Join(ms, ", ")
}
