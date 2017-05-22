package rpcclient

import "github.com/kbuzsaki/cupid/server"

type clientGlue struct {
	delegate RPCServer
}

func New(addr string) server.Server {
	return &clientGlue{NewClient(addr)}
}

func (cg *clientGlue) KeepAlive(li server.LeaseInfo) ([]server.Event, error) {

	leases := []string{}

	for _, l := range li.LockedNodes {
		lease := l.Serialize()
		leases = append(leases, lease)
	}

	leaseInfo := LeaseInfo{leases}
	events := []server.Event{}
	err := cg.delegate.KeepAlive(&leaseInfo, &events)

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (cg *clientGlue) Open(path string, readOnly bool, events server.EventsConfig) (server.NodeDescriptor, error) {
	args := OpenArgs{path, readOnly, events}
	nd := ""
	err := cg.delegate.Open(&args, &nd)
	if err != nil {
		return server.NodeDescriptor{}, err
	}

	return server.DeserializeNodeDescriptor(nd)
}

func (cg *clientGlue) Acquire(node server.NodeDescriptor) error {
	return cg.delegate.Acquire(node.Serialize(), nil)
}

func (cg *clientGlue) TryAcquire(node server.NodeDescriptor) (bool, error) {
	ok := false
	err := cg.delegate.TryAcquire(node.Serialize(), &ok)
	return ok, err
}

func (cg *clientGlue) Release(node server.NodeDescriptor) error {
	return cg.delegate.Release(node.Serialize(), nil)
}

func (cg *clientGlue) GetContentAndStat(node server.NodeDescriptor) (server.NodeContentAndStat, error) {
	cas := server.NodeContentAndStat{}
	err := cg.delegate.GetContentAndStat(node.Serialize(), &cas)
	return cas, err
}

func (cg *clientGlue) GetStat(node server.NodeDescriptor) (server.NodeStat, error) {
	stat := server.NodeStat{}
	err := cg.delegate.GetStat(node.Serialize(), &stat)
	return stat, err
}

func (cg *clientGlue) SetContent(node server.NodeDescriptor, content string, generation uint64) (bool, error) {
	setContentArgs := SetContentArgs{node.Serialize(), content, generation}
	ok := false
	err := cg.delegate.SetContent(&setContentArgs, &ok)
	return ok, err
}
