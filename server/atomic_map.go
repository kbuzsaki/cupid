package server

import "sync"

type AtomicMap interface {
	Get(k uint64) interface{}
	Put(k uint64, v interface{})
	Keys() []uint64
}

type atomicMapImpl struct {
	lock sync.RWMutex
	data map[uint64]interface{}
}

func NewAtomicMap() AtomicMap {
	return &atomicMapImpl{
		data: make(map[uint64]interface{}),
	}
}

func (am *atomicMapImpl) Get(k uint64) interface{} {
	am.lock.RLock()
	defer am.lock.RUnlock()
	return am.data[k]
}

func (am *atomicMapImpl) Put(k uint64, v interface{}) {
	am.lock.Lock()
	defer am.lock.Unlock()
	am.data[k] = v
}

func (am *atomicMapImpl) Keys() []uint64 {
	am.lock.RLock()
	defer am.lock.RUnlock()
	var keys []uint64

	for key := range am.data {
		keys = append(keys, key)
	}

	return keys
}
