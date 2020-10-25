package glmt

import "context"

type Notifier interface {
	Send(ctx context.Context, message string) error
}
