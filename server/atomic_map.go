package server

import "sync"

type AtomicMap interface {
	Get(k uint64) interface{}
	Put(k uint64, v interface{})
	Delete(k uint64)
	Keys() []uint64
}

type atomicMapImpl struct {
	lock sync.RWMutex
	def  func(uint64) interface{}
	data map[uint64]interface{}
}

func NewAtomicMap() AtomicMap {
	return NewAtomicMapWithDefault(nil)
}

func NewAtomicMapWithDefault(def func(uint64) interface{}) AtomicMap {
	return &atomicMapImpl{
		def:  def,
		data: make(map[uint64]interface{}),
	}
}

func (am *atomicMapImpl) Get(k uint64) interface{} {
	am.lock.RLock()
	defer am.lock.RUnlock()

	if _, ok := am.data[k]; !ok && am.def != nil {
		am.data[k] = am.def(k)
	}
	return am.data[k]
}

func (am *atomicMapImpl) Put(k uint64, v interface{}) {
	am.lock.Lock()
	defer am.lock.Unlock()

	am.data[k] = v
}

func (am *atomicMapImpl) Delete(k uint64) {
	am.lock.Lock()
	defer am.lock.Unlock()

	delete(am.data, k)
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

type AtomicStringMap interface {
	Get(k string) interface{}
	Put(k string, v interface{})
	Delete(k string)
	Keys() []string
}

type atomicStringMapImpl struct {
	lock sync.RWMutex
	def  func(string) interface{}
	data map[string]interface{}
}

func NewAtomicStringMap() AtomicStringMap {
	return NewAtomicStringMapWithDefault(nil)
}

func NewAtomicStringMapWithDefault(def func(string) interface{}) AtomicStringMap {
	return &atomicStringMapImpl{
		def:  def,
		data: make(map[string]interface{}),
	}
}

func (am *atomicStringMapImpl) Get(k string) interface{} {
	am.lock.RLock()
	defer am.lock.RUnlock()

	if _, ok := am.data[k]; !ok && am.def != nil {
		am.data[k] = am.def(k)
	}
	return am.data[k]
}

func (am *atomicStringMapImpl) Put(k string, v interface{}) {
	am.lock.Lock()
	defer am.lock.Unlock()

	am.data[k] = v
}

func (am *atomicStringMapImpl) Delete(k string) {
	am.lock.Lock()
	defer am.lock.Unlock()

	delete(am.data, k)
}

func (am *atomicStringMapImpl) Keys() []string {
	am.lock.RLock()
	defer am.lock.RUnlock()

	var keys []string
	for key := range am.data {
		keys = append(keys, key)
	}
	return keys
}
