package glmt

import "context"

type Notifier interface {
	Send(ctx context.Context, args map[string]string, add string) error
}
