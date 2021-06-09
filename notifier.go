package notifiers

import (
	"context"

	"github.com/pkg/errors"
)

// ErrNotify is notify error.
var ErrNotify = errors.New("notify error")

// Notifier is Notify interface.
type Notifier interface {
	// GetName return notifier type name
	GetName() string
	// Notify notify msg
	Notify(ctx context.Context, to string, msg Message) error
	// Close close notifier
	Close() error
	Wait()
}

// NoopCloser impl a noop io.Closer.
type NoopCloser struct{}

// Close impl io.Closer.
func (nc *NoopCloser) Close() error {
	return nil
}

// NoopWaiter impl a noop Waiter
type NoopWaiter struct{}

// Wait imple wait
func (nw *NoopWaiter) Wait() {}
