package notifiers

// Message is notify msg type.
type Message struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}

// Option is config function type.
type Option func(message *Message)

// WithTitle set title to msg.
func WithTitle(title string) Option {
	return func(message *Message) {
		message.Title = title
	}
}

// MessageFromContent is a helper function creating a message from content and some options.
func MessageFromContent(content string, opts ...Option) Message {
	msg := &Message{Content: content}

	for _, f := range opts {
		f(msg)
	}

	return *msg
}
