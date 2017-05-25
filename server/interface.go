package server

import (
	"log"
	"strings"
	"time"
)

type Server interface {
	KeepAlive(li LeaseInfo, eis []EventInfo) ([]Event, error)

	Open(path string, readOnly bool, events EventsConfig) (NodeDescriptor, error)

	Acquire(node NodeDescriptor) error
	TryAcquire(node NodeDescriptor) (bool, error)
	Release(node NodeDescriptor) error

	GetContentAndStat(node NodeDescriptor) (NodeContentAndStat, error)
	GetStat(node NodeDescriptor) (NodeStat, error)
	SetContent(node NodeDescriptor, content string, generation uint64) (bool, error)
}

type NodeDescriptor struct {
	descriptor descriptorKey
	path       string
}

func (nd NodeDescriptor) Path() string {
	return nd.path
}

func (nd NodeDescriptor) GobEncode() ([]byte, error) {
	return []byte(nd.Serialize()), nil
}

func (ndp *NodeDescriptor) GobDecode(b []byte) error {
	nd, err := DeserializeNodeDescriptor(string(b))
	if err != nil {
		return err
	}

	*ndp = nd

	return nil
}

func (nd NodeDescriptor) Serialize() string {
	return nd.descriptor.Serialize() + "|" + nd.path
}

func DeserializeNodeDescriptor(s string) (NodeDescriptor, error) {
	ss := strings.SplitN(s, "|", 2)
	if len(ss) != 2 {
		log.Printf("invalid node descriptor: %#v\n", s)
		return NodeDescriptor{}, ErrInvalidNodeDescriptor
	}

	descriptor, err := DeserializeDescriptorKey(ss[0])
	if err != nil {
		return NodeDescriptor{}, err
	}

	path := ss[1]

	return NodeDescriptor{descriptor, path}, nil
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
