package notifiers

import (
	"context"
	"strings"
	"sync"

	"go.uber.org/multierr"
)

type Combine struct {
	name      string
	notifiers []Notifier
}

func NewCombine(notifiers ...Notifier) *Combine {
	c := &Combine{notifiers: notifiers}
	c.name = c.getName()

	return c
}

func (c *Combine) GetName() string {
	return c.name
}

func (c *Combine) Close() error {
	errs := make([]error, 0)
	for _, n := range c.notifiers {
		errs = append(errs, n.Close())
	}
	return multierr.Combine(errs...)
}

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
