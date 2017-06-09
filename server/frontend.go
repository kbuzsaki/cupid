package server

import (
	"errors"
	"sync"
	"time"
)

const (
	maxKeepAliveDelay = 3 * time.Second
)

var (
	ErrInvalidSessionDescriptor = errors.New("Invalid session descriptor")
	ErrInvalidNodeDescriptor    = errors.New("Invalid node descriptor")
	ErrReadOnlyNodeDescriptor   = errors.New("Write from read-only node descriptor")
)

func minTime(keepAliveDelay time.Duration) time.Duration {
	if keepAliveDelay > maxKeepAliveDelay {
		return maxKeepAliveDelay
	}

	return keepAliveDelay
}

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
		lastKeepAlive: time.Now(),
		signaler:      NewSignaler(),
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
	fsm FSM

	csLock sync.RWMutex
	cs     ClusterState

	sessions  AtomicMap
	lockLocks AtomicStringMap
	setLocks  AtomicStringMap
}

func NewFrontend() (Server, error) {
	fsm, err := NewFSM()
	if err != nil {
		return nil, err
	}

	c := make(chan string, 1)
	c2 := make(chan *string, 1)
	go func() {
		for {
			s := <-c
			c2 <- &s
		}
	}()

	fsm = NewRaftFSM(c, c2, fsm)
	stateC := make(chan ClusterState, 1)
	stateC <- ClusterState{true, 1, ""}
	return NewFrontendWithFSM(fsm, stateC)
}

func NewFrontendWithFSM(fsm FSM, stateChanges <-chan ClusterState) (Server, error) {
	fe := &frontendImpl{
		fsm:       fsm,
		sessions:  NewAtomicMap(),
		lockLocks: NewAtomicStringMapWithDefault(func(string) interface{} { return &sync.Mutex{} }),
		setLocks:  NewAtomicStringMapWithDefault(func(string) interface{} { return &sync.Mutex{} }),
	}

	cs := <-stateChanges
	fe.setClusterState(cs)

	go func() {
		for cs := range stateChanges {
			fe.setClusterState(cs)
		}
	}()

	return fe, nil
}

func (fe *frontendImpl) setClusterState(cs ClusterState) {
	fe.csLock.Lock()
	defer fe.csLock.Unlock()
	wasLeader := fe.cs.IsLeader
	fe.cs = cs

	if wasLeader && !cs.IsLeader {
		fe.sessions = NewAtomicMap()
		fe.lockLocks = NewAtomicStringMapWithDefault(func(string) interface{} { return &sync.Mutex{} })
		fe.setLocks = NewAtomicStringMapWithDefault(func(string) interface{} { return &sync.Mutex{} })
	} else if !wasLeader && cs.IsLeader {
		sds := fe.fsm.GetSessionDescriptors()
		for _, sd := range sds {
			fe.sessions.Put(uint64(sd.Descriptor), NewSessionConn())
		}

		// TODO: finish propagating events
		// TODO: make double failover work properly
		unfinalized := fe.fsm.GetUnfinalizedNodes()
		for _, ni := range unfinalized {
			mut := fe.setLocks.Get(ni.path).(*sync.Mutex)
			mut.Lock()
			go fe.finalizeSetContent(ni)
		}
	}
}

func (fe *frontendImpl) getClusterState() ClusterState {
	fe.csLock.RLock()
	defer fe.csLock.RUnlock()
	return fe.cs
}

func (fe *frontendImpl) KeepAlive(li LeaseInfo, eis []EventInfo, keepAliveDelay time.Duration) ([]Event, error) {
	if cs := fe.getClusterState(); !cs.IsLeader {
		return nil, cs.MakeRedirectError()
	}

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
	if cs := fe.getClusterState(); !cs.IsLeader {
		return SessionDescriptor{}, cs.MakeRedirectError()
	}

	sd := fe.fsm.OpenSession()
	fe.sessions.Put(uint64(sd.Descriptor), NewSessionConn())
	return sd, nil
}

func (fe *frontendImpl) CloseSession(sd SessionDescriptor) error {
	if cs := fe.getClusterState(); !cs.IsLeader {
		return cs.MakeRedirectError()
	}

	fe.fsm.CloseSession(sd)
	// TODO: internal cleanup?
	fe.sessions.Delete(uint64(sd.Descriptor))
	return nil
}

func (fe *frontendImpl) Open(sd SessionDescriptor, path string, readOnly bool, config EventsConfig) (NodeDescriptor, error) {
	if cs := fe.getClusterState(); !cs.IsLeader {
		return NodeDescriptor{}, cs.MakeRedirectError()
	}

	if session := fe.fsm.GetSession(sd); session == nil {
		return NodeDescriptor{}, ErrInvalidSessionDescriptor
	}

	return fe.fsm.OpenNode(sd, path, readOnly, config), nil
}

func (fe *frontendImpl) CloseNode(nd NodeDescriptor) error {
	if cs := fe.getClusterState(); !cs.IsLeader {
		return cs.MakeRedirectError()
	}

	fe.fsm.CloseNode(nd)
	// TODO: internal cleanup?
	return nil
}

func (fe *frontendImpl) Acquire(nd NodeDescriptor) error {
	if cs := fe.getClusterState(); !cs.IsLeader {
		return cs.MakeRedirectError()
	}

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
	if cs := fe.getClusterState(); !cs.IsLeader {
		return false, cs.MakeRedirectError()
	}

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
	if cs := fe.getClusterState(); !cs.IsLeader {
		return cs.MakeRedirectError()
	}

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
	if cs := fe.getClusterState(); !cs.IsLeader {
		return NodeContentAndStat{}, cs.MakeRedirectError()
	}

	if node := fe.fsm.GetNodeDescriptor(nd); node == nil {
		return NodeContentAndStat{}, ErrInvalidNodeDescriptor
	}

	return fe.fsm.GetContentAndStat(nd), nil
}

func createInvalidationEvent(nd NodeDescriptor, cas NodeContentAndStat, config EventsConfig) Event {
	if config.ContentModified {
		return ContentInvalidationPushEvent{nd, cas}
	}
	return ContentInvalidationEvent{nd}
}

func (fe *frontendImpl) SetContent(nd NodeDescriptor, content string, generation uint64) (bool, error) {
	if cs := fe.getClusterState(); !cs.IsLeader {
		return false, cs.MakeRedirectError()
	}

	var nid *nodeDescriptor
	if nid = fe.fsm.GetNodeDescriptor(nd); nid == nil {
		return false, ErrInvalidNodeDescriptor
	} else if nid.readOnly {
		return false, ErrReadOnlyNodeDescriptor
	}

	mut := fe.setLocks.Get(nid.ni.path).(*sync.Mutex)
	mut.Lock()

	ok := fe.fsm.PrepareSetContent(nd, NodeContentAndStat{content, NodeStat{generation, time.Now()}})
	if !ok {
		mut.Unlock()
		return false, nil
	}

	fe.finalizeSetContent(nid.ni)

	return true, nil
}

func (fe *frontendImpl) finalizeSetContent(ni *nodeInfo) {
	wg := sync.WaitGroup{}

	cas := ni.GetContentAndStat()
	sds := fe.sessions.Keys()
	for _, sd := range sds {
		if sd == uint64(nd.Session.Descriptor) {
			continue
		}
		// if the session has been closed since, just ignore it
		session, ok := fe.sessions.Get(sd).(*sessionConn)
		if !ok {
			continue
		}

		cs := fe.fsm.GetSession(SessionDescriptor{descriptorKey(sd)})
		keys := cs.GetDescriptorKeys(ni.path)
		wg.Add(len(keys))

		for _, key := range keys {
			go func(key descriptorKey) {
				snd := NodeDescriptor{SessionDescriptor{cs.key}, key, ni.path}
				event := createInvalidationEvent(snd, cas, fe.fsm.GetNodeDescriptor(snd).config)
				session.SendEvent(event)
				wg.Done()
			}(key)
		}
	}

	wg.Wait()

	fe.fsm.FinalizeSetContent(ni.path)

	mut := fe.setLocks.Get(ni.path).(*sync.Mutex)
	mut.Unlock()
}
