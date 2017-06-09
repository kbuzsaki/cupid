package server

import (
	"errors"
	"sync"
	"time"
)

const (
	timeoutThreshold = 3 * maxKeepAliveDelay
)

var (
	ErrLockNotHeld = errors.New("Attempting to release a lock that is not held")
)

type nodeInfo struct {
	path         string
	content      string
	lastModified time.Time
	generation   uint64
	finalized    bool
	lock         sync.RWMutex
	locker       *nodeDescriptor
}

func (ni *nodeInfo) GetContentAndStat() NodeContentAndStat {
	ni.lock.RLock()
	defer ni.lock.RUnlock()

	return NodeContentAndStat{
		ni.content,
		NodeStat{
			ni.generation,
			ni.lastModified,
		},
	}
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
	ni.finalized = false
	return true
}

func (ni *nodeInfo) FinalizeSetContent() {
	ni.lock.Lock()
	defer ni.lock.Unlock()
	ni.finalized = true
}

func (ni *nodeInfo) SetLocked(locker *nodeDescriptor) {
	ni.lock.Lock()
	defer ni.lock.Unlock()

	ni.locker = locker
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
