package client

import "github.com/kbuzsaki/cupid/server"

// TODO: ensure correctness in the face of closed channels
func BufferEvents(in <-chan server.Event, rawOut chan<- server.Event) {
	var buffer []server.Event

	var out chan<- server.Event
	var nextSend server.Event

	for {
		select {
		case event := <-in:
			buffer = append(buffer, event)
			if out == nil {
				nextSend = buffer[0]
				buffer = buffer[1:]
				out = rawOut
			}
		case out <- nextSend:
			if len(buffer) > 0 {
				nextSend = buffer[0]
				buffer = buffer[1:]
			} else {
				out = nil
			}
		}
	}
}
