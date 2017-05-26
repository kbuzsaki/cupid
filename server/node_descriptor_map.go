package server

import "sync"

type descriptorKey uint64

type clientSession struct {
	lock      sync.RWMutex
	data      map[descriptorKey]*nodeDescriptor
	ndsByPath map[string][]*nodeDescriptor
	nextKey   descriptorKey

	signaler Signaler
}

func newClientSession() *clientSession {
	return &clientSession{
		data:      make(map[descriptorKey]*nodeDescriptor),
		ndsByPath: make(map[string][]*nodeDescriptor),
		signaler:  NewSignaler(),
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

func (cs *clientSession) GetDescriptors(path string) []*nodeDescriptor {
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
	nd := &nodeDescriptor{cs, ni, readOnly}
	cs.data[cs.nextKey] = nd
	cs.ndsByPath[ni.path] = append(cs.ndsByPath[ni.path], nd)
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
	sdm.data[sdm.nextKey] = newClientSession()
	return sdm.nextKey
}

type nodeDescriptor struct {
	cs       *clientSession
	ni       *nodeInfo
	readOnly bool
}
