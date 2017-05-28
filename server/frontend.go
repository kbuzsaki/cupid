package server

import (
	"sync"
	"time"
)

type pendingEvent struct {
	event Event
	ackC  chan struct{}
}

// TODO: when an event times out, put the session into a "dead" state until next keepalive so we don't waste timouts on it
type sessionConn struct {
	inKeepAlive   bool
	lastKeepAlive time.Time

	lock    sync.Mutex
	pending []pendingEvent
	ackCs   []chan struct{}

	signaler Signaler
}

func NewSessionConn() *sessionConn {
	return &sessionConn{
		signaler: NewSignaler(),
	}
}

func (sc *sessionConn) EnterKeepAlive() {
	sc.inKeepAlive = true
}

func (sc *sessionConn) ExitKeepAlive() {
	sc.inKeepAlive = false
	sc.lastKeepAlive = time.Now()
}

// SendEvent sends an event to this session and either blocks until the session acks it or times out
func (sc *sessionConn) SendEvent(event Event) bool {
	ac := make(chan struct{})

	sc.lock.Lock()
	sc.pending = append(sc.pending, pendingEvent{event, ac})
	sc.lock.Unlock()

	sc.signaler.Signal()

	select {
	case <-ac:
		return true
	case <-time.Tick(timeoutThreshold):
		return false
	}
}

func (sc *sessionConn) ReadEvents() []Event {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	var events []Event
	for _, pe := range sc.pending {
		events = append(events, pe.event)
		sc.ackCs = append(sc.ackCs, pe.ackC)
	}
	sc.pending = nil

	return events
}

func (sc *sessionConn) AckEvents() {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	for _, ackC := range sc.ackCs {
		close(ackC)
	}
	sc.ackCs = nil
}

// TODO: error handling
type frontendImpl struct {
	fsm      FSM
	sessions AtomicMap
}

func NewFrontend() (Server, error) {
	fsm, err := NewFSM()
	if err != nil {
		return nil, err
	}

	c := make(chan string)
	fsm = NewRaftFSM(c, c, fsm)
	return &frontendImpl{fsm: fsm, sessions: NewAtomicMap()}, nil
}

func (fe *frontendImpl) KeepAlive(li LeaseInfo, eis []EventInfo, keepAliveDelay time.Duration) ([]Event, error) {
	sc := fe.sessions.Get(uint64(li.Session.Descriptor)).(*sessionConn)
	sc.EnterKeepAlive()
	defer sc.ExitKeepAlive()
	sc.AckEvents()

	var events []Event
	select {
	case <-time.Tick(minTime(keepAliveDelay)):
		// timeout, no events
	case <-sc.signaler.DoneChan():
		sc.signaler.Reset()
		events = sc.ReadEvents()
	}

	return events, nil
}

func (fe *frontendImpl) OpenSession() (SessionDescriptor, error) {
	sd := fe.fsm.OpenSession()
	fe.sessions.Put(uint64(sd.Descriptor), NewSessionConn())
	return sd, nil
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

	ok := fe.fsm.SetContent(nd, NodeContentAndStat{content, NodeStat{generation, time.Now()}})
	if !ok {
		return false, nil
	}

	cas := fe.fsm.GetContentAndStat(nd)
	sds := fe.sessions.Keys()
	for _, sd := range sds {
		session := fe.sessions.Get(sd).(*sessionConn)
		cs := fe.fsm.GetSession(SessionDescriptor{descriptorKey(sd)})
		keys := cs.GetDescriptorKeys(nd.Path)
		for _, key := range keys {
			snd := NodeDescriptor{SessionDescriptor{cs.key}, key, nd.Path}
			event := ContentInvalidationPushEvent{snd, cas}
			session.SendEvent(event)
		}
	}

	return true, nil
}
