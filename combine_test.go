package notifiers

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombine(t *testing.T) {
	var (
		b1 bytes.Buffer
		b2 bytes.Buffer
	)
	p1 := NewPrinter(&b1)
	p2 := NewPrinter(&b2)

	c := NewCombine(p1, p2)
	assert.Equal(t, "combine printer printer", c.GetName())

	c.Notify(context.Background(), "", MessageFromContent("test"))
	assert.Equal(t, `{"content":"test","title":""}`, b1.String())
	assert.Equal(t, `{"content":"test","title":""}`, b2.String())
}
