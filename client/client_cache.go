package client

import (
	"sync"

	"github.com/kbuzsaki/cupid/server"
)

type nodeCache struct {
	casLock  sync.RWMutex
	casCache map[server.NodeDescriptor]server.NodeContentAndStat
}

func newNodeCache() nodeCache {
	return nodeCache{casCache: make(map[server.NodeDescriptor]server.NodeContentAndStat)}
}

func (nc *nodeCache) Get(nd server.NodeDescriptor) (server.NodeContentAndStat, bool) {
	nc.casLock.RLock()
	defer nc.casLock.RUnlock()

	cas, ok := nc.casCache[nd]
	return cas, ok
}

func (nc *nodeCache) Put(nd server.NodeDescriptor, cas server.NodeContentAndStat) {
	nc.casLock.Lock()
	defer nc.casLock.Unlock()

	nc.casCache[nd] = cas
}

func (nc *nodeCache) Delete(nd server.NodeDescriptor) {
	nc.casLock.Lock()
	defer nc.casLock.Unlock()

	delete(nc.casCache, nd)
}

func (nc *nodeCache) GetEventInfos() []server.EventInfo {
	var eis []server.EventInfo
	nc.casLock.RLock()
	defer nc.casLock.RUnlock()

	for nd := range nc.casCache {
		cas := nc.casCache[nd]
		eis = append(eis, server.EventInfo{nd, cas.Stat.Generation, true})
	}

	return eis
}

// lockSet keeps track of which locks a client holds
type lockSet struct {
	heldLocksLock sync.RWMutex
	heldLocks     map[server.NodeDescriptor]struct{}
}

func newLockSet() lockSet {
	return lockSet{heldLocks: make(map[server.NodeDescriptor]struct{})}
}

func (ls *lockSet) Contains(nd server.NodeDescriptor) bool {
	ls.heldLocksLock.RLock()
	defer ls.heldLocksLock.RUnlock()

	_, ok := ls.heldLocks[nd]
	return ok
}

func (ls *lockSet) Add(nd server.NodeDescriptor) {
	ls.heldLocksLock.Lock()
	defer ls.heldLocksLock.Unlock()

	ls.heldLocks[nd] = struct{}{}
}

func (ls *lockSet) Remove(nd server.NodeDescriptor) {
	ls.heldLocksLock.Lock()
	defer ls.heldLocksLock.Unlock()

	delete(ls.heldLocks, nd)
}

func (ls *lockSet) GetLeaseInfo() server.LeaseInfo {
	var li server.LeaseInfo
	ls.heldLocksLock.RLock()
	defer ls.heldLocksLock.RUnlock()

	for nd := range ls.heldLocks {
		li.LockedNodes = append(li.LockedNodes, nd)
	}

	return li
}
