package client

import "github.com/kbuzsaki/cupid/server"

type clientImpl struct {
	s server.Server
}

// TODO: accept a config file instead?
func New(addr string) (Client, error) {
	// TODO: actually instantiate server over rpc transport
	return &clientImpl{s: nil}, nil
}

func (cl *clientImpl) Open(path string, readOnly bool) (NodeHandle, error) {
	nd, err := cl.s.Open(path, readOnly)
	if err != nil {
		return nil, err
	}

	return &nodeHandleImpl{cl, nd}, nil
}

type nodeHandleImpl struct {
	cl *clientImpl
	nd server.NodeDescriptor
}

func (nh *nodeHandleImpl) Close() error {
	panic("implement me")
}

func (nh *nodeHandleImpl) Acquire() error {
	return nh.cl.s.Acquire(nh.nd)
}

func (nh *nodeHandleImpl) TryAcquire() (bool, error) {
	return nh.cl.s.TryAcquire(nh.nd)
}

func (nh *nodeHandleImpl) Release() error {
	return nh.cl.s.Release(nh.nd)
}

func (nh *nodeHandleImpl) GetContentAndStat() (server.NodeContentAndStat, error) {
	return nh.cl.s.GetContentAndStat(nh.nd)
}

func (nh *nodeHandleImpl) GetStat() (server.NodeStat, error) {
	return nh.cl.s.GetStat(nh.nd)
}

func (nh *nodeHandleImpl) SetContent(contents string, generation uint64) (bool, error) {
	return nh.cl.s.SetContent(nh.nd, contents, generation)
}

func (nh *nodeHandleImpl) Delete() error {
	panic("implement me")
}
