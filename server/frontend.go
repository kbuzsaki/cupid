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
	aliveLock     sync.Mutex
	inKeepAlive   bool
	lastKeepAlive time.Time

	eventLock sync.Mutex
	pending   []pendingEvent
	ackCs     []chan struct{}

	signaler Signaler
}

func NewSessionConn() *sessionConn {
	return &sessionConn{
		signaler: NewSignaler(),
	}
}

func (sc *sessionConn) EnterKeepAlive() {
	sc.aliveLock.Lock()
	defer sc.aliveLock.Unlock()

	sc.inKeepAlive = true
}

func (sc *sessionConn) ExitKeepAlive() {
	sc.aliveLock.Lock()
	defer sc.aliveLock.Unlock()

	sc.inKeepAlive = false
	sc.lastKeepAlive = time.Now()
}

func (sc *sessionConn) IsAlive() bool {
	sc.aliveLock.Lock()
	defer sc.aliveLock.Unlock()

	if sc.inKeepAlive {
		return true
	}

	return time.Since(sc.lastKeepAlive) < timeoutThreshold
}

// SendEvent sends an event to this session and either blocks until the session acks it or times out
func (sc *sessionConn) SendEvent(event Event) bool {
	ac := make(chan struct{})

	sc.eventLock.Lock()
	sc.pending = append(sc.pending, pendingEvent{event, ac})
	sc.eventLock.Unlock()

	sc.signaler.Signal()

	if !sc.IsAlive() {
		return false
	}

	select {
	case <-ac:
		return true
	case <-time.Tick(timeoutThreshold):
		// TODO: maybe reduce this from timeoutThreshold to (timeoutThreshold - time.Since(sc.lastKeepAlive))
		return false
	}
}

func (sc *sessionConn) ReadEvents() []Event {
	sc.eventLock.Lock()
	defer sc.eventLock.Unlock()

	var events []Event
	for _, pe := range sc.pending {
		events = append(events, pe.event)
		sc.ackCs = append(sc.ackCs, pe.ackC)
	}
	sc.pending = nil

	return events
}

func (sc *sessionConn) AckEvents() {
	sc.eventLock.Lock()
	defer sc.eventLock.Unlock()

	for _, ackC := range sc.ackCs {
		close(ackC)
	}
	sc.ackCs = nil
}

// TODO: error handling
type frontendImpl struct {
	fsm       FSM
	sessions  AtomicMap
	lockLocks AtomicStringMap
}

func NewFrontend() (Server, error) {
	fsm, err := NewFSM()
	if err != nil {
		return nil, err
	}

	c := make(chan string)
	fsm = NewRaftFSM(c, c, fsm)
	return &frontendImpl{
		fsm:       fsm,
		sessions:  NewAtomicMap(),
		lockLocks: NewAtomicStringMapWithDefault(func(string) interface{} { return &sync.Mutex{} }),
	}, nil
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
	for {
		locked, err := fe.TryAcquire(nd)
		if err != nil {
			return err
		} else if locked {
			return nil
		}
	}
}

func (fe *frontendImpl) TryAcquire(nd NodeDescriptor) (bool, error) {
	var nid *nodeDescriptor
	if nid = fe.fsm.GetNodeDescriptor(nd); nid == nil {
		return false, ErrInvalidNodeDescriptor
	} else if nid.readOnly {
		return false, ErrReadOnlyNodeDescriptor
	}

	lock := fe.lockLocks.Get(nid.ni.path).(*sync.Mutex)
	lock.Lock()
	defer lock.Unlock()

	currentLocker := nid.ni.locker
	if currentLocker == nil {
		// there is no locker, so take the lock
		fe.fsm.SetLocked(nd)
		return true, nil
	}

	lockerSession := fe.sessions.Get(uint64(currentLocker.cs.key)).(*sessionConn)
	if !lockerSession.IsAlive() {
		// the locker died, so take the lock and send them an event
		// TODO: what if the leader dies here?
		fe.fsm.SetLocked(nd)
		lockInvalidationEvent := LockInvalidationEvent{currentLocker.GetND()}
		lockerSession.SendEvent(lockInvalidationEvent)
		return true, nil
	}

	// we don't get the lock :(
	return false, nil
}

func (fe *frontendImpl) Release(nd NodeDescriptor) error {
	if node := fe.fsm.GetNodeDescriptor(nd); node == nil {
		return ErrInvalidNodeDescriptor
	} else if node.readOnly {
		return ErrReadOnlyNodeDescriptor
	}

	if ok := fe.fsm.ReleaseLock(nd); !ok {
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
