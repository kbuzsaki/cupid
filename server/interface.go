package server

import "time"

type Server interface {
	KeepAlive(li LeaseInfo, eis []EventInfo, keepAliveDelay time.Duration) ([]Event, error)

	OpenSession() (SessionDescriptor, error)
	Open(sd SessionDescriptor, path string, readOnly bool, events EventsConfig) (NodeDescriptor, error)

	Acquire(node NodeDescriptor) error
	TryAcquire(node NodeDescriptor) (bool, error)
	Release(node NodeDescriptor) error

	GetContentAndStat(node NodeDescriptor) (NodeContentAndStat, error)
	SetContent(node NodeDescriptor, content string, generation uint64) (bool, error)
}

type SessionDescriptor struct {
	Descriptor descriptorKey
}

type NodeDescriptor struct {
	Session    SessionDescriptor
	Descriptor descriptorKey
	Path       string
}

// LeaseInfo represents a session and the locks tha a session thinks it holds
type LeaseInfo struct {
	// list of locks
	LockedNodes []NodeDescriptor
}

type EventInfo struct {
	Descriptor NodeDescriptor
	Generation uint64
	Push       bool
}

type NodeContentAndStat struct {
	Content string
	Stat    NodeStat
}

type NodeStat struct {
	Generation   uint64
	LastModified time.Time
}
