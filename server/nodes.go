package server

import (
	"sync"
	"time"
)

type nodeInfo struct {
	content      string
	lastModified time.Time
	generation   uint64
	lock         sync.RWMutex
	// TODO: keep track of who locked the node
	locked bool
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

func (ni *nodeInfo) Acquire() {
	locked := ni.TryAcquire()
	for !locked {
		locked = ni.TryAcquire()
	}
}

func (ni *nodeInfo) TryAcquire() bool {
	ni.lock.Lock()
	defer ni.lock.Unlock()

	if ni.locked {
		return false
	}

	ni.locked = true
	return true
}

// TODO: add error checking for whether the caller has the lock
func (ni *nodeInfo) Release() {
	ni.lock.Lock()
	defer ni.lock.Unlock()

	ni.locked = false
}
