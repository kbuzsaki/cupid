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

// maps aren't safe for concurrent access, so guard mutations with a RWMutex
type nodeInfoMap struct {
	data map[string]*nodeInfo
	lock sync.RWMutex
}

func makeNodeInfoMap() *nodeInfoMap {
	return &nodeInfoMap{data: make(map[string]*nodeInfo)}
}

func (nim *nodeInfoMap) GetNode(path string) *nodeInfo {
	nim.lock.RLock()
	defer nim.lock.RUnlock()
	return nim.data[path]
}

func (nim *nodeInfoMap) CreateNode(path string) *nodeInfo {
	nim.lock.Lock()
	defer nim.lock.Unlock()

	if _, ok := nim.data[path]; !ok {
		nim.data[path] = &nodeInfo{}
	}

	return nim.data[path]
}

func (nim *nodeInfoMap) GetOrCreateNode(path string) *nodeInfo {
	if node := nim.GetNode(path); node != nil {
		return node
	}

	return nim.CreateNode(path)
}
