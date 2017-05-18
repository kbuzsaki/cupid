package server

import (
	"fmt"
	"strconv"
	"sync"
)

const descriptorKeyBase = 10

type descriptorKey uint64

func (k descriptorKey) Serialize() string {
	return strconv.FormatUint(uint64(k), descriptorKeyBase)
}

func DeserializeDescriptorKey(s string) (descriptorKey, error) {
	n, err := strconv.ParseUint(s, descriptorKeyBase, 0)
	if err != nil {
		return 0, fmt.Errorf("Unable to parse serialized descriptor: %v", err)
	}

	return descriptorKey(n), nil
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
