package server

import "time"

type Server interface {
	KeepAlive() error

	Open(path string, readOnly bool) (NodeDescriptor, error)

	Acquire(node NodeDescriptor) error
	TryAcquire(node NodeDescriptor) (bool, error)
	Release(node NodeDescriptor) error

	GetContentAndStat(node NodeDescriptor) (NodeContentAndStat, error)
	GetStat(node NodeDescriptor) (NodeStat, error)
	SetContent(node NodeDescriptor, content string, generation uint64) (bool, error)
}

type NodeDescriptor struct {
	descriptor string
}

type NodeContentAndStat struct {
	Content string
	Stat    NodeStat
}

type NodeStat struct {
	Generation   uint64
	LastModified time.Time
}
