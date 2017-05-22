package server

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrLockNotHeld = errors.New("Attempting to release a lock that is not held")
)

type nodeInfo struct {
	content      string
	lastModified time.Time
	generation   uint64
	lock         sync.RWMutex
	locker       *nodeDescriptor
}

func (ni *nodeInfo) GetContent() (string, time.Time, uint64) {
	ni.lock.RLock()
	defer ni.lock.RUnlock()
	return ni.content, ni.lastModified, ni.generation
}

func (ni *nodeInfo) SetContent(content string, generation uint64) bool {
	ni.lock.Lock()
	defer ni.lock.Unlock()

	if generation < ni.generation {
		return false
	}

	ni.content = content
	ni.lastModified = time.Now()
	ni.generation += 1
	return true
}

func (ni *nodeInfo) Acquire(locker *nodeDescriptor) {
	locked := ni.TryAcquire(locker)
	for !locked {
		locked = ni.TryAcquire(locker)
	}
}

func (ni *nodeInfo) TryAcquire(locker *nodeDescriptor) bool {
	ni.lock.Lock()
	defer ni.lock.Unlock()

	if ni.locker != nil {
		return false
	}

	ni.locker = locker
	return true
}

// TODO: add error checking for whether the caller has the lock
func (ni *nodeInfo) Release(unlocker *nodeDescriptor) error {
	ni.lock.Lock()
	defer ni.lock.Unlock()

	if unlocker == nil {
		return ErrInvalidNodeDescriptor
	}

	if ni.locker != unlocker {
		return ErrLockNotHeld
	}

	ni.locker = nil

	return nil
}
