package server

import "sync"

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
		nim.data[path] = &nodeInfo{path: path, finalized: true}
	}

	return nim.data[path]
}

func (nim *nodeInfoMap) GetOrCreateNode(path string) *nodeInfo {
	if node := nim.GetNode(path); node != nil {
		return node
	}

	return nim.CreateNode(path)
}

func (nim *nodeInfoMap) GetUnfinalizedNodes() []*nodeInfo {
	nim.lock.RLock()
	defer nim.lock.RUnlock()

	var unfinalized []*nodeInfo
	for key := range nim.data {
		ni := nim.data[key]
		if !ni.finalized {
			unfinalized = append(unfinalized, ni)
		}
	}
	return unfinalized
}
