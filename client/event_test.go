package client

import (
	"testing"

	"github.com/kbuzsaki/cupid/server"
	"github.com/stretchr/testify/assert"
)

func TestBufferEvents(t *testing.T) {
	in := make(chan server.Event)
	out := make(chan server.Event)

	go BufferEvents(in, out)

	send := []server.Event{"a", "b", "c"}
	for _, s := range send {
		in <- s
	}

	receive := []server.Event{}
	for range send {
		r := <-out
		receive = append(receive, r)
	}

	assert.Equal(t, len(send), len(receive), "sent and received must have the same length")

	for i := range send {
		assert.Equal(t, send[i], receive[i], "sent and receive must have same value at index %d", i)
	}
}
