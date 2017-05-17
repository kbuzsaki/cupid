package rpcclient

import (
	"time"

	"github.com/kbuzsaki/cupid/server"
)

type RPCServer interface {
	KeepAlive(lease *LeaseInfo, _ *int) error
	Open(args *OpenArgs, node *NodeDescriptor) error

	Acquire(node *NodeDescriptor, _ *int) error
	TryAcquire(node *NodeDescriptor, success *bool) error
	Release(node *NodeDescriptor, _ *int) error

	GetContentAndStat(node *NodeDescriptor, cas *NodeContentAndStat) error
	GetStat(node *NodeDescriptor, stat *NodeStat) error
	SetContent(node *NodeDescriptor, cas *NodeContentAndStat) error
}

type LeaseInfo struct {
}

type OpenArgs struct {
	Path     string
	ReadOnly bool
}

// TODO: make this an arbitrary int associated with the permissions on the descriptor
type NodeDescriptor struct {
	Descriptor string
}

type NodeContentAndStat struct {
	Content string
	Stat    NodeStat
}

type NodeStat struct {
	Generation   uint64
	LastModified time.Time
}

type rpcServer struct {
	delegate server.Server
}

func (rs *rpcServer) KeepAlive(lease *LeaseInfo, _ *int) error {
	return rs.delegate.KeepAlive()
}

func (rs *rpcServer) Open(args *OpenArgs, node *NodeDescriptor) error {
	nd, err := rs.delegate.Open(args.Path, args.ReadOnly)
	if err != nil {
		return err
	}

	*node = nd
	return nil
}

func (rs *rpcServer) Acquire(node *NodeDescriptor, _ *int) error {
	panic("implement me")
}

func (rs *rpcServer) TryAcquire(node *NodeDescriptor, success *bool) error {
	panic("implement me")
}

func (rs *rpcServer) Release(node *NodeDescriptor, _ *int) error {
	panic("implement me")
}

func (rs *rpcServer) GetContentAndStat(node *NodeDescriptor, cas *NodeContentAndStat) error {
	res, err := rs.delegate.GetContentAndStat(*node)
	if err != nil {
		return err
	}

	*cas = res
	return nil
}

func (rs *rpcServer) GetStat(node *NodeDescriptor, stat *NodeStat) error {
	res, err := rs.delegate.GetStat(*node)
	if err != nil {
		return err
	}

	*stat = res
	return nil
}

func (rs *rpcServer) SetContent(node *NodeDescriptor, cas *NodeContentAndStat) error {
	return rs.delegate.SetContent(*node, cas.Content, cas.Stat.Generation)
}
