package notifiers

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
)

type w struct {
	bt *bytes.Buffer
	t  *testing.T
}

func (ww *w) Write(p []byte) (n int, err error) {
	ww.t.Log(string(p))
	ww.bt.Write(p)
	return 0, nil
}

type errNotifier struct {
}

func (e *errNotifier) GetName() string {
	return "err"
}

func (e *errNotifier) Close() error {
	return nil
}

func (e *errNotifier) Notify(ctx context.Context, to string, msg Message) error {
	return errors.New(msg.Content)
}

func TestLimiter(t *testing.T) {
	l := NewLimiter(NewPrinter(&w{bt: &bytes.Buffer{}, t: t}), time.Second, 10)

	l.Notify(context.Background(), "", MessageFromContent("test1"))
	l.Notify(context.Background(), "", MessageFromContent("test2"))
	l.Notify(context.Background(), "", MessageFromContent("test3"))

	go func() {
		time.Sleep(time.Second * 4)
		l.Close()
	}()

	for e := range l.GetErrorCh() {
		t.Log(e)
	}
}

func TestLimiter_Error(t *testing.T) {
	l := NewLimiter(&errNotifier{}, time.Millisecond*50, 10)

	i := 0

	for i < 20 {
		l.Notify(context.Background(), "", MessageFromContent("test1"))
		i++
	}

	go func() {
		for e := range l.GetErrorCh() {
			t.Log(e)
		}
	}()

	//time.Sleep(time.Second * 2)
	l.Close()
}
