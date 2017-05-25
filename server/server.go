package server

import (
	"errors"
	"time"
)

const (
	keepAliveSleep = 1 * time.Second
)

var (
	ErrInvalidNodeDescriptor  = errors.New("Invalid node descriptor")
	ErrReadOnlyNodeDescriptor = errors.New("Write from read-only node descriptor")
)

type serverImpl struct {
	nodes       *nodeInfoMap
	descriptors *nodeDescriptorMap
}

// TODO: accept a config file instead?
func New() (Server, error) {
	return &serverImpl{
		nodes:       makeNodeInfoMap(),
		descriptors: makeNodeDescriptorMap(),
	}, nil
}

func (s *serverImpl) KeepAlive(li LeaseInfo, eis []EventInfo) ([]Event, error) {
	defer func() {
		for _, nd := range li.LockedNodes {
			nid := s.descriptors.GetDescriptor(nd.descriptor)
			if nid != nil && nid.ni.OwnedBy(nid) {
				nid.ni.ExitKeepAlive()
			}
		}
	}()
	// check if you own all the locks you think you own
	var events []Event
	for _, nd := range li.LockedNodes {
		nid := s.descriptors.GetDescriptor(nd.descriptor)
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

	// TODO: maybe wake up early when events happen?
	time.Sleep(keepAliveSleep)

	for _, ei := range eis {
		// TODO: maybe check nd validity earlier so that the error doesn't wait a second before being sent
		cas, err := s.GetContentAndStat(ei.Descriptor)
		if err != nil {
			return nil, err
		}

		if cas.Stat.Generation > ei.Generation {
			if ei.Push {
				events = append(events, ContentInvalidationPushEvent{ei.Descriptor, cas})
			} else {
				events = append(events, ContentInvalidationEvent{ei.Descriptor})
			}
		}
	}

	return events, nil
}

func (s *serverImpl) Open(path string, readOnly bool, events EventsConfig) (NodeDescriptor, error) {
	// TODO: semantics regarding deletion and whether a descriptor is still valid after deletion
	ni := s.nodes.GetOrCreateNode(path)
	key := s.descriptors.OpenDescriptor(ni, readOnly)
	return NodeDescriptor{key}, nil
}

func (s *serverImpl) Acquire(nd NodeDescriptor) error {
	nid := s.descriptors.GetDescriptor(nd.descriptor)
	if nid == nil {
		return ErrInvalidNodeDescriptor
	} else if nid.readOnly {
		return ErrReadOnlyNodeDescriptor
	}

	nid.ni.Acquire(nid)

	return nil
}

func (s *serverImpl) TryAcquire(nd NodeDescriptor) (bool, error) {
	nid := s.descriptors.GetDescriptor(nd.descriptor)
	if nid == nil {
		return false, ErrInvalidNodeDescriptor
	} else if nid.readOnly {
		return false, ErrReadOnlyNodeDescriptor
	}

	return nid.ni.TryAcquire(nid), nil
}

// TODO: prevent releases by an unrelated node descriptor
func (s *serverImpl) Release(nd NodeDescriptor) error {
	nid := s.descriptors.GetDescriptor(nd.descriptor)
	if nid == nil {
		return ErrInvalidNodeDescriptor
	} else if nid.readOnly {
		return ErrReadOnlyNodeDescriptor
	}

	return nid.ni.Release(nid)
}

func (s *serverImpl) GetContentAndStat(nd NodeDescriptor) (NodeContentAndStat, error) {
	nid := s.descriptors.GetDescriptor(nd.descriptor)
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

func (s *serverImpl) GetStat(nd NodeDescriptor) (NodeStat, error) {
	nid := s.descriptors.GetDescriptor(nd.descriptor)
	if nid == nil {
		return NodeStat{}, ErrInvalidNodeDescriptor
	}

	_, lastModified, generation := nid.ni.GetContent()

	return NodeStat{
		Generation:   generation,
		LastModified: lastModified,
	}, nil
}

func (s *serverImpl) SetContent(nd NodeDescriptor, content string, generation uint64) (bool, error) {
	nid := s.descriptors.GetDescriptor(nd.descriptor)
	if nid == nil {
		return false, ErrInvalidNodeDescriptor
	} else if nid.readOnly {
		return false, ErrReadOnlyNodeDescriptor
	}

	return nid.ni.SetContent(content, generation), nil
}
