package server

import "time"

// TODO: error handling
type frontendImpl struct {
	fsm FSM
}

func NewFrontend() (Server, error) {
	c := make(chan string)
	fsm, err := NewFSM()
	fsm = NewRaftFSM(c, c, fsm)
	if err != nil {
		return nil, err
	}

	return &frontendImpl{fsm}, nil
}

func (fe *frontendImpl) KeepAlive(li LeaseInfo, eis []EventInfo, keepAliveDelay time.Duration) ([]Event, error) {
	return nil, nil
}

func (fe *frontendImpl) OpenSession() (SessionDescriptor, error) {
	return fe.fsm.OpenSession(), nil
}

func (fe *frontendImpl) Open(sd SessionDescriptor, path string, readOnly bool, events EventsConfig) (NodeDescriptor, error) {
	if session := fe.fsm.GetSession(sd); session == nil {
		return NodeDescriptor{}, ErrInvalidSessionDescriptor
	}

	return fe.fsm.OpenNode(sd, path, readOnly), nil
}

func (fe *frontendImpl) Acquire(nd NodeDescriptor) error {
	if node := fe.fsm.GetNodeDescriptor(nd); node == nil {
		return ErrInvalidNodeDescriptor
	} else if node.readOnly {
		return ErrReadOnlyNodeDescriptor
	}

	locked := fe.fsm.TryAcquire(nd)
	for !locked {
		locked = fe.fsm.TryAcquire(nd)
	}
	return nil
}

func (fe *frontendImpl) TryAcquire(nd NodeDescriptor) (bool, error) {
	if node := fe.fsm.GetNodeDescriptor(nd); node == nil {
		return false, ErrInvalidNodeDescriptor
	} else if node.readOnly {
		return false, ErrReadOnlyNodeDescriptor
	}

	return fe.fsm.TryAcquire(nd), nil
}

func (fe *frontendImpl) Release(nd NodeDescriptor) error {
	if node := fe.fsm.GetNodeDescriptor(nd); node == nil {
		return ErrInvalidNodeDescriptor
	} else if node.readOnly {
		return ErrReadOnlyNodeDescriptor
	}

	if ok := fe.fsm.Release(nd); !ok {
		return ErrLockNotHeld
	}
	return nil
}

func (fe *frontendImpl) GetContentAndStat(nd NodeDescriptor) (NodeContentAndStat, error) {
	if node := fe.fsm.GetNodeDescriptor(nd); node == nil {
		return NodeContentAndStat{}, ErrInvalidNodeDescriptor
	}

	return fe.fsm.GetContentAndStat(nd), nil
}

func (fe *frontendImpl) SetContent(nd NodeDescriptor, content string, generation uint64) (bool, error) {
	if node := fe.fsm.GetNodeDescriptor(nd); node == nil {
		return false, ErrInvalidNodeDescriptor
	} else if node.readOnly {
		return false, ErrReadOnlyNodeDescriptor
	}

	return fe.fsm.SetContent(nd, NodeContentAndStat{content, NodeStat{generation, time.Now()}}), nil
}
