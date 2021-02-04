package notifiers

import (
	"context"

	"github.com/pkg/errors"
)

var ErrNotify = errors.New("notify error")

type Notifier interface {
	GetName() string
	Notify(ctx context.Context, to string, msg Message) error
	Close() error
}

type NoopCloser struct{}

func (nc *NoopCloser) Close() error {
	return nil
}
