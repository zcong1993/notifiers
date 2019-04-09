package types

import (
	"time"

	"github.com/google/uuid"
)

// Message is message struct
type Message struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Tags    []string  `json:"tags"`
	Content string    `json:"content"`
	URL     string    `json:"url"`
	Time    time.Time `json:"time"`
	Exts    []string  `json:"exts"`
}

// Notifier is notifier interface
type Notifier interface {
	Notify(msg *Message) error
}

func WrapMsg(msg *Message) *Message {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}

	if msg.Time.IsZero() {
		msg.Time = time.Now()
	}

	return msg
}
