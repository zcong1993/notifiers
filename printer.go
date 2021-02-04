package notifiers

import (
	"context"
	"encoding/json"
	"io"
)

type Printer struct {
	writer io.Writer
}

func NewPrinter(writer io.Writer) *Printer {
	return &Printer{writer: writer}
}

func (p *Printer) GetName() string {
	return "printer"
}

func (p *Printer) Close() error {
	if closer, ok := p.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (p *Printer) Notify(ctx context.Context, to string, msg Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = p.writer.Write(b)
	return err
}
