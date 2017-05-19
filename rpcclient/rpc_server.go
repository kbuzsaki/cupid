package rpcclient

import (
	"github.com/kbuzsaki/cupid/server"
)

type RPCServer interface {
	KeepAlive(lease *LeaseInfo, _ *int) error
	Open(args *OpenArgs, nd *string) error

	Acquire(node string, _ *int) error
	TryAcquire(node string, success *bool) error
	Release(node string, _ *int) error

	GetContentAndStat(node string, cas *server.NodeContentAndStat) error
	GetStat(node string, stat *server.NodeStat) error
	SetContent(args *SetContentArgs, success *bool) error
}

type LeaseInfo struct {
}

type OpenArgs struct {
	Path     string
	ReadOnly bool
}

type SetContentArgs struct {
	SNode      string
	Content    string
	Generation uint64
}

type rpcServer struct {
	delegate server.Server
}

func NewServer(delegate server.Server) RPCServer {
	return &rpcServer{delegate}
}

func (rs *rpcServer) KeepAlive(lease *LeaseInfo, _ *int) error {
	return nil
}

func (rs *rpcServer) Open(args *OpenArgs, nd *string) error {
	descriptor, err := rs.delegate.Open(args.Path, args.ReadOnly)
	if err != nil {
		return err
	}

	*nd = descriptor.Serialize()
	return nil
}

func (rs *rpcServer) Acquire(snd string, _ *int) error {
	node, err := server.DeserializeNodeDescriptor(snd)
	if err != nil {
		return err
	}

	return rs.delegate.Acquire(node)
}

func (rs *rpcServer) TryAcquire(snd string, success *bool) error {
	node, err := server.DeserializeNodeDescriptor(snd)
	if err != nil {
		*success = false
		return err
	}

	succ, err := rs.delegate.TryAcquire(node)
	if err != nil {
		*success = false
		return err
	}

	*success = succ
	return nil
}

func (rs *rpcServer) Release(snd string, _ *int) error {
	node, err := server.DeserializeNodeDescriptor(snd)
	if err != nil {
		return err
	}

	return rs.delegate.Release(node)
}

func (rs *rpcServer) GetContentAndStat(snd string, cas *server.NodeContentAndStat) error {
	node, err := server.DeserializeNodeDescriptor(snd)
	if err != nil {
		return err
	}

	nodeCas, err := rs.delegate.GetContentAndStat(node)
	if err != nil {
		return err
	}

	*cas = nodeCas
	return nil
}

func (rs *rpcServer) GetStat(snd string, stat *server.NodeStat) error {
	node, err := server.DeserializeNodeDescriptor(snd)
	if err != nil {
		return err
	}

	nodeStat, err := rs.delegate.GetStat(node)
	if err != nil {
		return err
	}

	*stat = nodeStat
	return nil
}

func (rs *rpcServer) SetContent(args *SetContentArgs, success *bool) error {
	node, err := server.DeserializeNodeDescriptor(args.SNode)
	if err != nil {
		return err
	}

	succ, err := rs.delegate.SetContent(node, args.Content, args.Generation)
	if err != nil {
		*success = false
		return err
	}

	*success = succ
	return nil
}
