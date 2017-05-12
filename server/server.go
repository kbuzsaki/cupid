package server

import "errors"

var (
	ErrInvalidNodeDescriptor = errors.New("Invalid node descriptor")
)

type serverImpl struct {
	nodes *nodeInfoMap
}

// TODO: accept a config file instead?
func New() (Server, error) {
	return &serverImpl{
		nodes: makeNodeInfoMap(),
	}, nil
}

func (s *serverImpl) KeepAlive() error {
	return nil
}

func (s *serverImpl) Open(path string, readOnly bool) (NodeDescriptor, error) {
	// TODO: actually use readOnly to restrict SetContent
	// TODO: semantics regarding deletion and whether a descriptor is still valid after deletion
	s.nodes.GetOrCreateNode(path)
	return NodeDescriptor{descriptor: path}, nil
}

func (s *serverImpl) Acquire(nd NodeDescriptor) error {
	ni := s.nodes.GetNode(nd.descriptor)
	if ni == nil {
		return ErrInvalidNodeDescriptor
	}

	ni.Acquire()

	return nil
}

func (s *serverImpl) TryAcquire(nd NodeDescriptor) (bool, error) {
	ni := s.nodes.GetNode(nd.descriptor)
	if ni == nil {
		return false, ErrInvalidNodeDescriptor
	}

	return ni.TryAcquire(), nil
}

func (s *serverImpl) Release(nd NodeDescriptor) error {
	ni := s.nodes.GetNode(nd.descriptor)
	if ni == nil {
		return ErrInvalidNodeDescriptor
	}

	ni.Release()

	return nil
}

func (s *serverImpl) GetContentAndStat(nd NodeDescriptor) (NodeContentAndStat, error) {
	ni := s.nodes.GetNode(nd.descriptor)
	if ni == nil {
		return NodeContentAndStat{}, ErrInvalidNodeDescriptor
	}

	content, lastModified, generation := ni.GetContent()

	return NodeContentAndStat{
		Content: content,
		Stat: NodeStat{
			Generation:   generation,
			LastModified: lastModified,
		},
	}, nil
}

func (s *serverImpl) GetStat(nd NodeDescriptor) (NodeStat, error) {
	ni := s.nodes.GetNode(nd.descriptor)
	if ni == nil {
		return NodeStat{}, ErrInvalidNodeDescriptor
	}

	_, lastModified, generation := ni.GetContent()

	return NodeStat{
		Generation:   generation,
		LastModified: lastModified,
	}, nil
}

func (s *serverImpl) SetContent(nd NodeDescriptor, content string, generation uint64) (bool, error) {
	ni := s.nodes.GetNode(nd.descriptor)
	if ni == nil {
		return false, ErrInvalidNodeDescriptor
	}

	return ni.SetContent(content, generation), nil
}
