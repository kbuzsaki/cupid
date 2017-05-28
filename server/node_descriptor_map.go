package server

import (
	"log"
	"sync"
	"time"
)

const (
	sessionTimeout = 5 * time.Second
)

type descriptorKey uint64

type sessionEvent struct {
	event   ContentInvalidationPushEvent
	ackChan chan struct{}
}

type clientSession struct {
	key descriptorKey

	lock      sync.RWMutex
	data      map[descriptorKey]*nodeDescriptor
	ndsByPath map[string][]descriptorKey
	nextKey   descriptorKey

	events   chan sessionEvent
	ackChans []chan struct{}
}

func newClientSession(key descriptorKey) *clientSession {
	return &clientSession{
		key:       key,
		data:      make(map[descriptorKey]*nodeDescriptor),
		ndsByPath: make(map[string][]descriptorKey),
		events:    make(chan sessionEvent),
	}
}

func (cs *clientSession) GetSD() SessionDescriptor {
	return SessionDescriptor{cs.key}
}

func (cs *clientSession) InvalidateCache(nd NodeDescriptor, cas NodeContentAndStat) {
	ackChan := make(chan struct{})
	se := sessionEvent{
		event:   ContentInvalidationPushEvent{nd, cas},
		ackChan: ackChan,
	}

	// we have sent the event to KeepAlive
	cs.events <- se

	// now wait for ack or timeout
	select {
	case <-time.Tick(sessionTimeout):
		// TODO: do we have to clean anything up if we time out?
		log.Println("Timed out for session:", cs.key, nd)
	case <-ackChan:
		// success!
	}
}

func (cs *clientSession) GetDescriptor(key descriptorKey) *nodeDescriptor {
	if cs == nil {
		return nil
	}

	cs.lock.RLock()
	defer cs.lock.RUnlock()

	return cs.data[key]
}

func (cs *clientSession) GetDescriptorKeys(path string) []descriptorKey {
	if cs == nil {
		return nil
	}

	cs.lock.RLock()
	defer cs.lock.RUnlock()

	return cs.ndsByPath[path]
}

func (cs *clientSession) OpenDescriptor(ni *nodeInfo, readOnly bool) descriptorKey {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.nextKey++
	nd := &nodeDescriptor{cs, cs.nextKey, ni, readOnly}
	cs.data[cs.nextKey] = nd
	cs.ndsByPath[ni.path] = append(cs.ndsByPath[ni.path], cs.nextKey)
	return cs.nextKey
}

func (cs *clientSession) CloseDescriptor(key descriptorKey) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	delete(cs.data, key)
}

type sessionDescriptorMap struct {
	lock    sync.RWMutex
	data    map[descriptorKey]*clientSession
	nextKey descriptorKey
}

func makeSessionDescriptorMap() *sessionDescriptorMap {
	return &sessionDescriptorMap{data: make(map[descriptorKey]*clientSession)}
}

func (sdm *sessionDescriptorMap) GetDescriptor(nd NodeDescriptor) *nodeDescriptor {
	return sdm.GetSession(nd.Session.Descriptor).GetDescriptor(nd.Descriptor)
}

func (sdm *sessionDescriptorMap) GetSession(key descriptorKey) *clientSession {
	sdm.lock.RLock()
	defer sdm.lock.RUnlock()
	return sdm.data[key]
}

func (sdm *sessionDescriptorMap) GetSessionDescriptors() []descriptorKey {
	sdm.lock.RLock()
	defer sdm.lock.RUnlock()

	var sds []descriptorKey
	for key := range sdm.data {
		sds = append(sds, key)
	}
	return sds
}

func (sdm *sessionDescriptorMap) OpenSession() descriptorKey {
	sdm.lock.Lock()
	defer sdm.lock.Unlock()
	sdm.nextKey++
	sdm.data[sdm.nextKey] = newClientSession(sdm.nextKey)
	return sdm.nextKey
}

type nodeDescriptor struct {
	cs       *clientSession
	key      descriptorKey
	ni       *nodeInfo
	readOnly bool
}

func (nd *nodeDescriptor) GetND() NodeDescriptor {
	return NodeDescriptor{nd.cs.GetSD(), nd.key, nd.ni.path}
}
