package notifiers

import (
	"context"
	"encoding/json"
	"io"
)

// Printer impl notifier, notify msg to a io.Writer.
type Printer struct {
	NoopWaiter
	writer io.Writer
}

// NewPrinter create a instance.
func NewPrinter(writer io.Writer) *Printer {
	return &Printer{writer: writer}
}

// GetName impl Notifier.GetName.
func (p *Printer) GetName() string {
	return "printer"
}

// Close impl Notifier.Close.
func (p *Printer) Close() error {
	return nil
}

// Notify impl Notifier.Notify.
// write json encoded msg to writer.
func (p *Printer) Notify(ctx context.Context, to string, msg Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = p.writer.Write(b)
	return err
}

var _ Notifier = (*Printer)(nil)
