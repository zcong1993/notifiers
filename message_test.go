package notifiers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zcong1993/notifiers/v2"
)

func TestMessageFromContent(t *testing.T) {
	assert.Equal(t, notifiers.Message{Content: "a"}, notifiers.MessageFromContent("a"))
	assert.Equal(t, notifiers.Message{Content: "a", Title: "b"}, notifiers.MessageFromContent("a", notifiers.WithTitle("b")))
}
