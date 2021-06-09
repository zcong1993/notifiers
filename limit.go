package notifiers

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

// Limiter impl notifier, can add a interval between messages.
// Limiter run a background goroutine to handle message notify.
type Limiter struct {
	notifier Notifier
	interval time.Duration
	doneCh   chan struct{}
	closed   int32
	msgCh    chan *msgWithTo
	errCh    chan error
	wg       sync.WaitGroup
}

type msgWithTo struct {
	to  string
	msg Message
}

// NewLimiter create a instance.
func NewLimiter(notifier Notifier, interval time.Duration, msgChSize int) *Limiter {
	l := &Limiter{
		notifier: notifier,
		interval: interval,
		doneCh:   make(chan struct{}),
		msgCh:    make(chan *msgWithTo, msgChSize),
		errCh:    make(chan error, 10),
	}

	go l.run()

	return l
}

// GetName impl Notifier.GetName.
// Name is "limiter " + inner notifier name.
func (l *Limiter) GetName() string {
	return "limiter " + l.notifier.GetName()
}

// Close impl Notifier.Close.
// It will wait unfinished messages before close.
func (l *Limiter) Close() error {
	atomic.StoreInt32(&l.closed, 1)
	l.wg.Wait()
	close(l.doneCh)
	close(l.errCh)
	return l.notifier.Close()
}

// Notify impl Notifier.Notify.
// This function is unblock, so return error always be nil.
// If you need error message, see Limiter.GetErrorCh().
func (l *Limiter) Notify(ctx context.Context, to string, msg Message) error {
	if atomic.LoadInt32(&l.closed) == 1 {
		return errors.New("limiter closed")
	}

	l.msgCh <- &msgWithTo{
		to:  to,
		msg: msg,
	}
	l.wg.Add(1)

	return nil
}

// GetErrorCh return error message channel.
func (l *Limiter) GetErrorCh() <-chan error {
	return l.errCh
}

func (l *Limiter) run() {
	for {
		select {
		case <-l.doneCh:
			return
		case msg := <-l.msgCh:
			err := l.notifier.Notify(context.Background(), msg.to, msg.msg)
			l.wg.Done()
			if err != nil {
				select {
				case l.errCh <- err:
				default:
				}
			}

			time.Sleep(l.interval)
		}
	}
}
