package server

import "sync"

type descriptorKey uint64

type sessionDescriptor struct {
	*nodeDescriptorMap
}

type sessionDescriptorMap struct {
	lock    sync.RWMutex
	data    map[descriptorKey]*sessionDescriptor
	nextKey descriptorKey
}

func makeSessionDescriptorMap() *sessionDescriptorMap {
	return &sessionDescriptorMap{data: make(map[descriptorKey]*sessionDescriptor)}
}

func (sdm *sessionDescriptorMap) GetDescriptor(nd NodeDescriptor) *nodeDescriptor {
	return sdm.GetSession(nd.Session.Descriptor).GetDescriptor(nd.Descriptor)
}

func (sdm *sessionDescriptorMap) GetSession(key descriptorKey) *sessionDescriptor {
	sdm.lock.RLock()
	defer sdm.lock.RUnlock()
	return sdm.data[key]
}

func (sdm *sessionDescriptorMap) OpenSession() descriptorKey {
	sdm.lock.Lock()
	defer sdm.lock.Unlock()
	sdm.nextKey++
	sdm.data[sdm.nextKey] = &sessionDescriptor{nodeDescriptorMap: makeNodeDescriptorMap()}
	return sdm.nextKey
}

type nodeDescriptor struct {
	ni       *nodeInfo
	readOnly bool
}

// maps aren't safe for concurrent access, so guard mutations with a RWMutex
type nodeDescriptorMap struct {
	lock    sync.RWMutex
	data    map[descriptorKey]*nodeDescriptor
	nextKey descriptorKey
}

func makeNodeDescriptorMap() *nodeDescriptorMap {
	return &nodeDescriptorMap{data: make(map[descriptorKey]*nodeDescriptor)}
}

func (nid *nodeDescriptorMap) GetDescriptor(key descriptorKey) *nodeDescriptor {
	if nid == nil {
		return nil
	}

	nid.lock.RLock()
	defer nid.lock.RUnlock()

	return nid.data[key]
}

func (nid *nodeDescriptorMap) OpenDescriptor(ni *nodeInfo, readOnly bool) descriptorKey {
	nid.lock.Lock()
	defer nid.lock.Unlock()
	nid.nextKey++
	nid.data[nid.nextKey] = &nodeDescriptor{ni, readOnly}
	return nid.nextKey
}

func (nid *nodeDescriptorMap) CloseDescriptor(key descriptorKey) {
	nid.lock.Lock()
	defer nid.lock.Unlock()
	delete(nid.data, key)
}
