package server

import "time"

type Server interface {
	KeepAlive(li LeaseInfo) ([]Event, error)

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
	return nd.descriptor.Serialize()
}

func DeserializeNodeDescriptor(s string) (NodeDescriptor, error) {
	descriptor, err := DeserializeDescriptorKey(s)
	if err != nil {
		return NodeDescriptor{}, err
	}

	return NodeDescriptor{descriptor}, nil
}

// LeaseInfo represents a session and the locks tha a session thinks it holds
type LeaseInfo struct {
	// list of locks
	LockedNodes []NodeDescriptor
}

type NodeContentAndStat struct {
	Content string
	Stat    NodeStat
}

type NodeStat struct {
	Generation   uint64
	LastModified time.Time
}
