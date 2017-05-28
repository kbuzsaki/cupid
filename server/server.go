package server

import (
	"errors"
	"log"
	"time"
)

const (
	maxKeepAliveDelay = 3 * time.Second
)

var (
	ErrInvalidSessionDescriptor = errors.New("Invalid session descriptor")
	ErrInvalidNodeDescriptor    = errors.New("Invalid node descriptor")
	ErrReadOnlyNodeDescriptor   = errors.New("Write from read-only node descriptor")
)

type serverImpl struct {
	nodes    *nodeInfoMap
	sessions *sessionDescriptorMap
}

// TODO: accept a config file instead?
func New() (Server, error) {
	return &serverImpl{
		nodes:    makeNodeInfoMap(),
		sessions: makeSessionDescriptorMap(),
	}, nil
}

func minTime(keepAliveDelay time.Duration) time.Duration {
	if keepAliveDelay > maxKeepAliveDelay {
		return maxKeepAliveDelay
	}

	return keepAliveDelay
}

func (s *serverImpl) KeepAlive(li LeaseInfo, eis []EventInfo, keepAliveDelay time.Duration) ([]Event, error) {
	defer func() {
		for _, nd := range li.LockedNodes {
			nid := s.sessions.GetDescriptor(nd)
			if nid != nil && nid.ni.OwnedBy(nid) {
				nid.ni.ExitKeepAlive()
			}
		}
	}()

	session := s.sessions.GetSession(li.Session.Descriptor)

	// now that we're back, we have to ack everything
	for _, ackChan := range session.ackChans {
		log.Println("acking for session:", li)
		close(ackChan)
	}
	session.ackChans = nil

	// check if you own all the locks you think you own
	// TODO: move lock timeouts to event stuff?
	var events []Event
	for _, nd := range li.LockedNodes {
		nid := s.sessions.GetDescriptor(nd)
		if nid == nil {
			return nil, ErrInvalidNodeDescriptor
		}

		if !nid.ni.OwnedBy(nid) {
			events = append(events, LockInvalidationEvent{nd})
		} else {
			nid.ni.EnterKeepAlive()
		}
	}

	if len(events) > 0 {
		return events, nil
	}

	// TODO: check events after calling reset but before listening on DoneChan
	// in case an event happened in between KeepAlive

	select {
	case <-time.Tick(minTime(keepAliveDelay)):
		// timeout
	case se := <-session.events:
		log.Println("got se:", se)
		session.ackChans = append(session.ackChans, se.ackChan)
		events = append(events, se.event)
	}

	return events, nil
}

func (s *serverImpl) OpenSession() (SessionDescriptor, error) {
	key := s.sessions.OpenSession()
	return SessionDescriptor{key}, nil
}

func (s *serverImpl) Open(sd SessionDescriptor, path string, readOnly bool, events EventsConfig) (NodeDescriptor, error) {
	// TODO: semantics regarding deletion and whether a descriptor is still valid after deletion
	session := s.sessions.GetSession(sd.Descriptor)
	if session == nil {
		return NodeDescriptor{}, ErrInvalidSessionDescriptor
	}

	ni := s.nodes.GetOrCreateNode(path)
	key := session.OpenDescriptor(ni, readOnly)
	return NodeDescriptor{sd, key, path}, nil
}

func (s *serverImpl) Acquire(nd NodeDescriptor) error {
	nid := s.sessions.GetDescriptor(nd)
	if nid == nil {
		return ErrInvalidNodeDescriptor
	} else if nid.readOnly {
		return ErrReadOnlyNodeDescriptor
	}

	nid.ni.Acquire(nid)

	return nil
}

func (s *serverImpl) TryAcquire(nd NodeDescriptor) (bool, error) {
	nid := s.sessions.GetDescriptor(nd)
	if nid == nil {
		return false, ErrInvalidNodeDescriptor
	} else if nid.readOnly {
		return false, ErrReadOnlyNodeDescriptor
	}

	return nid.ni.TryAcquire(nid), nil
}

// TODO: prevent releases by an unrelated node descriptor
func (s *serverImpl) Release(nd NodeDescriptor) error {
	nid := s.sessions.GetDescriptor(nd)
	if nid == nil {
		return ErrInvalidNodeDescriptor
	} else if nid.readOnly {
		return ErrReadOnlyNodeDescriptor
	}

	return nid.ni.Release(nid)
}

func (s *serverImpl) GetContentAndStat(nd NodeDescriptor) (NodeContentAndStat, error) {
	nid := s.sessions.GetDescriptor(nd)
	if nid == nil {
		return NodeContentAndStat{}, ErrInvalidNodeDescriptor
	}

	content, lastModified, generation := nid.ni.GetContent()

	return NodeContentAndStat{
		Content: content,
		Stat: NodeStat{
			Generation:   generation,
			LastModified: lastModified,
		},
	}, nil
}

func (s *serverImpl) SetContent(nd NodeDescriptor, content string, generation uint64) (bool, error) {
	nid := s.sessions.GetDescriptor(nd)
	if nid == nil {
		return false, ErrInvalidNodeDescriptor
	} else if nid.readOnly {
		return false, ErrReadOnlyNodeDescriptor
	}

	// TODO: fix race condition with concurrent gets here
	ok := nid.ni.SetContent(content, generation)
	if !ok {
		return false, nil
	}

	cas := nid.ni.GetContentAndStat()

	sds := s.sessions.GetSessionDescriptors()
	for _, sd := range sds {
		session := s.sessions.GetSession(sd)
		keys := session.GetDescriptorKeys(nd.Path)

		for _, key := range keys {
			snd := NodeDescriptor{SessionDescriptor{session.key}, key, nd.Path}
			session.InvalidateCache(snd, cas)
		}
	}

	return true, nil
}
