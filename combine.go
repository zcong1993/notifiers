package notifiers

import (
	"context"
	"strings"
	"sync"

	"go.uber.org/multierr"
)

// Combine impl notifier, can combine multi notifiers to one.
type Combine struct {
	name      string
	notifiers []Notifier
}

// NewCombine create a instance from other notifiers.
func NewCombine(notifiers ...Notifier) *Combine {
	c := &Combine{notifiers: notifiers}
	c.name = c.getName()

	return c
}

// GetName impl Notifier.GetName.
// Name is "combine " + inner notifier names.
func (c *Combine) GetName() string {
	return c.name
}

// Close impl Notifier.Close.
// It return multierr, use multierr.Errors unwrap to multi errors slice.
func (c *Combine) Close() error {
	errs := make([]error, 0)
	for _, n := range c.notifiers {
		errs = append(errs, n.Close())
	}
	return multierr.Combine(errs...)
}

// Wait impl Notifier.Wait
func (c *Combine) Wait() {
	var wg sync.WaitGroup
	for _, n := range c.notifiers {
		wg.Add(1)
		n := n
		go func() {
			n.Wait()
			wg.Done()
		}()
	}
	wg.Wait()
}

// Notify impl Notifier.Notify.
// Call inner notifiers parallelly, and return a multierr, use multierr.Errors unwrap to multi errors slice.
func (c *Combine) Notify(ctx context.Context, to string, msg Message) error {
	var wg sync.WaitGroup
	errs := make([]error, len(c.notifiers))
	for i, n := range c.notifiers {
		wg.Add(1)

		n := n
		i := i

		go func() {
			errs[i] = n.Notify(ctx, to, msg)
			wg.Done()
		}()
	}

	wg.Wait()
	return multierr.Combine(errs...)
}

func (c *Combine) getName() string {
	sb := strings.Builder{}

	sb.WriteString("combine")
	for _, n := range c.notifiers {
		sb.WriteString(" ")
		sb.WriteString(n.GetName())
	}
	return sb.String()
}

var _ Notifier = (*Combine)(nil)
