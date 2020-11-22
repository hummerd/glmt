package notifier

import (
	"context"

	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
)

type Notifier interface {
	Send(ctx context.Context, args map[string]string, add string, mentions []*team.Member) error
}
