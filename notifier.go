package notifiers

import "context"

type Notifier interface {
	GetName() string
	Notify(ctx context.Context, to string, msg Message) error
	Close() error
}
