package server

import "sync"

type descriptorKey uint64

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
