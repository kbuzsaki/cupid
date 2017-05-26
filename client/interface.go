package client

import (
	"io"

	"github.com/kbuzsaki/cupid/server"
)

type Client interface {
	Open(path string, readOnly bool, events server.EventsConfig) (NodeHandle, error)
	GetEventsOut() <-chan server.Event
}

type Locker interface {
	Acquire() error
	TryAcquire() (bool, error)
	Release() error
}

type File interface {
	GetContentAndStat() (server.NodeContentAndStat, error)
	SetContent(contents string, generation uint64) (bool, error)
}

type NodeHandle interface {
	io.Closer
	Locker
	File
	Delete() error
	Path() string
	Register(cb SubscriberCallback)
}
