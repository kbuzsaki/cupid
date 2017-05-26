package rpcclient

import (
	"time"

	"github.com/kbuzsaki/cupid/server"
)

type clientGlue struct {
	delegate RPCServer
}

func New(addr string, keepAliveDelay time.Duration) server.Server {
	return &clientGlue{NewClient(addr, keepAliveDelay)}
}

func (cg *clientGlue) KeepAlive(li server.LeaseInfo, eventsInfo []server.EventInfo, keepAliveDelay time.Duration) ([]server.Event, error) {
	events := []server.Event{}

	args := KeepAliveArgs{li, eventsInfo, keepAliveDelay}
	err := cg.delegate.KeepAlive(&args, &events)

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (cg *clientGlue) OpenSession() (server.SessionDescriptor, error) {
	sd := server.SessionDescriptor{}
	err := cg.delegate.OpenSession(0, &sd)
	if err != nil {
		return server.SessionDescriptor{}, err
	}

	return sd, nil
}

func (cg *clientGlue) Open(sd server.SessionDescriptor, path string, readOnly bool, events server.EventsConfig) (server.NodeDescriptor, error) {
	args := OpenArgs{sd, path, readOnly, events}
	nd := server.NodeDescriptor{}
	err := cg.delegate.Open(&args, &nd)
	if err != nil {
		return server.NodeDescriptor{}, err
	}

	return nd, nil
}

func (cg *clientGlue) Acquire(node server.NodeDescriptor) error {
	return cg.delegate.Acquire(node, nil)
}

func (cg *clientGlue) TryAcquire(node server.NodeDescriptor) (bool, error) {
	ok := false
	err := cg.delegate.TryAcquire(node, &ok)
	return ok, err
}

func (cg *clientGlue) Release(node server.NodeDescriptor) error {
	return cg.delegate.Release(node, nil)
}

func (cg *clientGlue) GetContentAndStat(node server.NodeDescriptor) (server.NodeContentAndStat, error) {
	cas := server.NodeContentAndStat{}
	err := cg.delegate.GetContentAndStat(node, &cas)
	return cas, err
}

func (cg *clientGlue) SetContent(node server.NodeDescriptor, content string, generation uint64) (bool, error) {
	setContentArgs := SetContentArgs{node, content, generation}
	ok := false
	err := cg.delegate.SetContent(&setContentArgs, &ok)
	return ok, err
}
