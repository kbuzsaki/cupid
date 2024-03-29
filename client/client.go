package client

import (
	"log"
	"sync"
	"time"

	"github.com/kbuzsaki/cupid/rpcclient"
	"github.com/kbuzsaki/cupid/server"
)

const (
	minimumKeepAliveDelay  = 100 * time.Millisecond
	connectionErrorBackoff = 5 * time.Second
)

type clientImpl struct {
	s  server.Server
	sd server.SessionDescriptor

	closingLock sync.Mutex
	closing     bool

	eventsIn  chan<- server.Event
	eventsOut <-chan server.Event

	nodeCache nodeCache
	locks     lockSet

	keepAliveDelay time.Duration
	subscriber     Subscriber
}

// TODO: accept a config file instead?
func New(addr string, keepAliveDelay time.Duration) (Client, error) {
	s := rpcclient.New(addr, keepAliveDelay)
	return newFromServer(s, keepAliveDelay)
}

func NewRaft(addrs []string, keepAliveDelay time.Duration) (Client, error) {
	var delegates []server.Server
	for _, addr := range addrs {
		delegates = append(delegates, rpcclient.New(addr, keepAliveDelay))
	}
	s := &RedirectServer{delegates: delegates, leader: 1, pendingLeader: 1}
	return newFromServer(s, keepAliveDelay)
}

func newFromServer(s server.Server, keepAliveDelay time.Duration) (Client, error) {
	eventsIn := make(chan server.Event)
	eventsOut := make(chan server.Event)
	go BufferEvents(eventsIn, eventsOut)

	cl := &clientImpl{
		s:              s,
		eventsIn:       eventsIn,
		eventsOut:      eventsOut,
		nodeCache:      newNodeCache(),
		locks:          newLockSet(),
		keepAliveDelay: keepAliveDelay,
	}
	subscriber, err := NewSubscriber(cl)
	if err != nil {
		return nil, err
	}
	cl.subscriber = subscriber

	sd, err := s.OpenSession()
	if err != nil {
		return nil, err
	}
	cl.sd = sd

	go cl.keepAlive()

	return cl, nil
}

func (cl *clientImpl) GetEventsOut() <-chan server.Event {
	return cl.eventsOut
}

func (cl *clientImpl) register(nd server.NodeDescriptor, cb SubscriberCallback) {
	cl.subscriber.Register(nd.Path, cb)
}

func (cl *clientImpl) handleEvents(events []server.Event) {
	// do things with those functions
	for _, rawEvent := range events {
		switch event := rawEvent.(type) {
		case server.LockInvalidationEvent:
			log.Println("handling lock invalidation event:", event)
			cl.locks.Remove(event.Descriptor)
		case server.ContentInvalidationEvent:
			cl.nodeCache.Delete(event.Descriptor)
		case server.ContentInvalidationPushEvent:
			cl.nodeCache.Put(event.Descriptor, event.NodeContentAndStat)
		default:
			log.Println("Unrecognized event:", rawEvent)
		}

		//
		cl.eventsIn <- rawEvent
	}
}

func (cl *clientImpl) isClosing() bool {
	cl.closingLock.Lock()
	defer cl.closingLock.Unlock()
	return cl.closing
}

func (cl *clientImpl) setClosing() {
	cl.closingLock.Lock()
	defer cl.closingLock.Unlock()
	cl.closing = true
}

// background does background KeepAlive processing in a goroutine
func (cl *clientImpl) keepAlive() {
	for !cl.isClosing() {
		li := cl.locks.GetLeaseInfo()
		li.Session = cl.sd
		//start := time.Now()
		events, err := cl.s.KeepAlive(li, cl.nodeCache.GetEventInfos(), cl.keepAliveDelay)
		//fmt.Println(time.Since(start))
		if err != nil {
			log.Println("KeepAlive error:", err)
			time.Sleep(connectionErrorBackoff)
		}

		if len(events) > 0 {
			cl.handleEvents(events)
		}
	}
}

func (cl *clientImpl) Open(path string, readOnly bool, config server.EventsConfig) (NodeHandle, error) {
	nd, err := cl.s.Open(cl.sd, path, readOnly, config)
	if err != nil {
		return nil, err
	}

	return &nodeHandleImpl{cl, nd}, nil
}

func (cl *clientImpl) Close() error {
	err := cl.s.CloseSession(cl.sd)
	if err != nil {
		return err
	}

	cl.setClosing()
	return nil
}

type nodeHandleImpl struct {
	cl *clientImpl
	nd server.NodeDescriptor
}

func (nh *nodeHandleImpl) Close() error {
	err := nh.cl.s.CloseNode(nh.nd)
	if err != nil {
		return err
	}

	nh.cl.nodeCache.Delete(nh.nd)
	nh.cl.nodeCache.Delete(nh.nd)
	return nil
}

func (nh *nodeHandleImpl) Acquire() error {
	if nh.cl.locks.Contains(nh.nd) {
		return nil
	}

	err := nh.cl.s.Acquire(nh.nd)
	if err != nil {
		return err
	}

	nh.cl.locks.Add(nh.nd)

	return nil
}

func (nh *nodeHandleImpl) TryAcquire() (bool, error) {
	if nh.cl.locks.Contains(nh.nd) {
		return true, nil
	}

	ok, err := nh.cl.s.TryAcquire(nh.nd)
	if err != nil {
		return false, err
	}

	if ok {
		nh.cl.locks.Add(nh.nd)
	}

	return ok, nil
}

func (nh *nodeHandleImpl) Release() error {
	err := nh.cl.s.Release(nh.nd)
	if err != nil {
		return err
	}

	nh.cl.locks.Remove(nh.nd)

	return nil
}

func (nh *nodeHandleImpl) GetContentAndStat() (server.NodeContentAndStat, error) {
	if cas, ok := nh.cl.nodeCache.Get(nh.nd); ok {
		return cas, nil
	}

	cas, err := nh.cl.s.GetContentAndStat(nh.nd)
	if err != nil {
		return server.NodeContentAndStat{}, err
	}

	nh.cl.nodeCache.Put(nh.nd, cas)

	return cas, nil
}

func (nh *nodeHandleImpl) SetContent(contents string, generation uint64) (bool, error) {
	return nh.cl.s.SetContent(nh.nd, contents, generation)
}

func (nh *nodeHandleImpl) Delete() error {
	panic("implement me")
}

func (nh *nodeHandleImpl) Path() string {
	return nh.nd.Path
}

func (nh *nodeHandleImpl) Register(cb SubscriberCallback) {
	nh.cl.register(nh.nd, cb)
}

func (nh *nodeHandleImpl) Nop(numOps uint64) error {
	return nh.cl.s.Nop(numOps)
}
