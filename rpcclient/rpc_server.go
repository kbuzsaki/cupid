package rpcclient

import (
	"errors"

	"github.com/kbuzsaki/cupid/server"
)

type RPCServer interface {
	KeepAlive(lease *LeaseInfo, events *[]server.Event) error
	Open(args *OpenArgs, nd *server.NodeDescriptor) error

	Acquire(node server.NodeDescriptor, _ *int) error
	TryAcquire(node server.NodeDescriptor, success *bool) error
	Release(node server.NodeDescriptor, _ *int) error

	GetContentAndStat(node server.NodeDescriptor, cas *server.NodeContentAndStat) error
	GetStat(node server.NodeDescriptor, stat *server.NodeStat) error
	SetContent(args *SetContentArgs, success *bool) error
}

type LeaseInfo struct {
	LockedNodes []string
}

type OpenArgs struct {
	Path         string
	ReadOnly     bool
	EventsConfig server.EventsConfig
}

type SetContentArgs struct {
	SNode        server.NodeDescriptor
	Content    string
	Generation uint64
}

type rpcServer struct {
	delegate server.Server
}

func NewServer(delegate server.Server) RPCServer {
	return &rpcServer{delegate}
}

func (rs *rpcServer) KeepAlive(li *LeaseInfo, events *[]server.Event) error {

	leases := server.LeaseInfo{}

	for _, l := range li.LockedNodes {
		nd, err := server.DeserializeNodeDescriptor(l)
		if err != nil {
			return errors.New("Failed to deserialize lease info")
		}

		leases.LockedNodes = append(leases.LockedNodes, nd)

	}

	tmp_events, err := rs.delegate.KeepAlive(leases)

	if err != nil {
		return err
	}

	*events = tmp_events
	return nil
}

func (rs *rpcServer) Open(args *OpenArgs, nd *server.NodeDescriptor) error {
	descriptor, err := rs.delegate.Open(args.Path, args.ReadOnly, args.EventsConfig)
	if err != nil {
		return err
	}

	*nd = descriptor
	return nil
}

func (rs *rpcServer) Acquire(snd server.NodeDescriptor, _ *int) error {
	return rs.delegate.Acquire(snd)
}

func (rs *rpcServer) TryAcquire(snd server.NodeDescriptor, success *bool) error {
	succ, err := rs.delegate.TryAcquire(snd)
	if err != nil {
		*success = false
		return err
	}

	*success = succ
	return nil
}

func (rs *rpcServer) Release(snd server.NodeDescriptor, _ *int) error {
	return rs.delegate.Release(snd)
}

func (rs *rpcServer) GetContentAndStat(snd server.NodeDescriptor, cas *server.NodeContentAndStat) error {
	nodeCas, err := rs.delegate.GetContentAndStat(snd)
	if err != nil {
		return err
	}

	*cas = nodeCas
	return nil
}

func (rs *rpcServer) GetStat(snd server.NodeDescriptor, stat *server.NodeStat) error {
	nodeStat, err := rs.delegate.GetStat(snd)
	if err != nil {
		return err
	}

	*stat = nodeStat
	return nil
}

func (rs *rpcServer) SetContent(args *SetContentArgs, success *bool) error {
	succ, err := rs.delegate.SetContent(args.SNode, args.Content, args.Generation)
	if err != nil {
		*success = false
		return err
	}

	*success = succ
	return nil
}
