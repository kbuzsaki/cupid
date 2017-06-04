package server

import "sync"

type descriptorKey uint64

type clientSession struct {
	key descriptorKey

	lock      sync.RWMutex
	data      map[descriptorKey]*nodeDescriptor
	ndsByPath map[string][]descriptorKey
	nextKey   descriptorKey
}

func newClientSession(key descriptorKey) *clientSession {
	return &clientSession{
		key:       key,
		data:      make(map[descriptorKey]*nodeDescriptor),
		ndsByPath: make(map[string][]descriptorKey),
	}
}

func (cs *clientSession) GetSD() SessionDescriptor {
	return SessionDescriptor{cs.key}
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

func (cs *clientSession) OpenDescriptor(ni *nodeInfo, readOnly bool, config EventsConfig) descriptorKey {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.nextKey++
	nd := &nodeDescriptor{cs, cs.nextKey, ni, readOnly, config}
	cs.data[cs.nextKey] = nd
	cs.ndsByPath[ni.path] = append(cs.ndsByPath[ni.path], cs.nextKey)
	return cs.nextKey
}

func filterNdsKey(nds []descriptorKey, key descriptorKey) []descriptorKey {
	var nnds []descriptorKey
	for _, nd := range nds {
		if nd != key {
			nnds = append(nnds, nd)
		}
	}
	return nnds
}

func (cs *clientSession) CloseDescriptor(key descriptorKey) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	if nid, ok := cs.data[key]; ok {
		nds := cs.ndsByPath[nid.ni.path]
		cs.ndsByPath[nid.ni.path] = filterNdsKey(nds, nid.key)
	}

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

func (sdm *sessionDescriptorMap) CloseSession(sd SessionDescriptor) {
	sdm.lock.Lock()
	defer sdm.lock.Unlock()
	delete(sdm.data, sd.Descriptor)
}

type nodeDescriptor struct {
	cs       *clientSession
	key      descriptorKey
	ni       *nodeInfo
	readOnly bool
	config   EventsConfig
}

func (nd *nodeDescriptor) GetND() NodeDescriptor {
	return NodeDescriptor{nd.cs.GetSD(), nd.key, nd.ni.path}
}
