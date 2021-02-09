package notifiers

type Message struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}

type Option func(message *Message)

func WithTitle(title string) Option {
	return func(message *Message) {
		message.Title = title
	}
}

func MessageFromContent(content string, opts ...Option) Message {
	msg := &Message{Content: content}

	for _, f := range opts {
		f(msg)
	}

	return *msg
}
