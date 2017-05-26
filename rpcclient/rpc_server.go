package rpcclient

import (
	"errors"

	"time"

	"github.com/kbuzsaki/cupid/server"
)

type RPCServer interface {
	Ping(_, _ *int) error

	KeepAlive(args *KeepAliveArgs, events *[]server.Event) error
	Open(args *OpenArgs, nd *server.NodeDescriptor) error

	Acquire(node server.NodeDescriptor, _ *int) error
	TryAcquire(node server.NodeDescriptor, success *bool) error
	Release(node server.NodeDescriptor, _ *int) error

	GetContentAndStat(node server.NodeDescriptor, cas *server.NodeContentAndStat) error
	SetContent(args *SetContentArgs, success *bool) error
}

type KeepAliveArgs struct {
	LeaseInfo      []string
	EventsInfo     []server.EventInfo
	KeepAliveDelay time.Duration
}

type OpenArgs struct {
	Path         string
	ReadOnly     bool
	EventsConfig server.EventsConfig
}

type SetContentArgs struct {
	SNode      server.NodeDescriptor
	Content    string
	Generation uint64
}

type rpcServer struct {
	delegate server.Server
}

func NewServer(delegate server.Server) RPCServer {
	return &rpcServer{delegate}
}

func (rs *rpcServer) Ping(_, _ *int) error {
	return nil
}

func (rs *rpcServer) KeepAlive(args *KeepAliveArgs, events *[]server.Event) error {

	leases := server.LeaseInfo{}

	for _, l := range args.LeaseInfo {
		nd, err := server.DeserializeNodeDescriptor(l)
		if err != nil {
			return errors.New("Failed to deserialize lease info")
		}

		leases.LockedNodes = append(leases.LockedNodes, nd)

	}

	tmp_events, err := rs.delegate.KeepAlive(leases, args.EventsInfo, args.KeepAliveDelay)

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

func (rs *rpcServer) SetContent(args *SetContentArgs, success *bool) error {
	succ, err := rs.delegate.SetContent(args.SNode, args.Content, args.Generation)
	if err != nil {
		*success = false
		return err
	}

	*success = succ
	return nil
}
