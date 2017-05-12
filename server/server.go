package server

import "fmt"

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

func (s *serverImpl) GetContentAndStat(nd NodeDescriptor) (NodeContentAndStat, error) {
	ni := s.nodes.GetNode(nd.descriptor)
	if ni == nil {
		return NodeContentAndStat{}, fmt.Errorf("Node %#v does not exist", nd.descriptor)
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
		return NodeStat{}, fmt.Errorf("Node %#v does not exist", nd.descriptor)
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
		return false, fmt.Errorf("Node %#v does not exist", nd.descriptor)
	}

	return ni.SetContent(content, generation), nil
}
