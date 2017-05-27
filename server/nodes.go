package server

import (
	"errors"
	"log"
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
	lock         sync.RWMutex
	locker       *nodeDescriptor
	keepAlive    bool
	lastPing     time.Time
}

func (ni *nodeInfo) GetContent() (string, time.Time, uint64) {
	ni.lock.RLock()
	defer ni.lock.RUnlock()
	return ni.content, ni.lastModified, ni.generation
}

func (ni *nodeInfo) GetContentAndStat() NodeContentAndStat {
	content, lastModified, generation := ni.GetContent()
	return NodeContentAndStat{
		content,
		NodeStat{
			generation,
			lastModified,
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

	if ni.locker != nil && (ni.keepAlive || time.Since(ni.lastPing) < timeoutThreshold) {
		log.Println("aborting try acquire: ", ni.keepAlive, time.Since(ni.lastPing))
		return false
	}
	log.Println("taking the lock from:", ni.locker, ni.keepAlive, time.Since(ni.lastPing))

	ni.locker = locker
	ni.lastPing = time.Now()

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

func (ni *nodeInfo) OwnedBy(nid *nodeDescriptor) bool {
	ni.lock.RLock()
	defer ni.lock.RUnlock()

	return ni.locker == nid
}

func (ni *nodeInfo) EnterKeepAlive() {
	ni.lock.Lock()
	defer ni.lock.Unlock()

	ni.keepAlive = true
}
func (ni *nodeInfo) ExitKeepAlive() {
	ni.lock.Lock()
	defer ni.lock.Unlock()

	ni.keepAlive = false
	ni.lastPing = time.Now()
}
