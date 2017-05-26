package rpcclient

import (
	"net/rpc"
	"sync"

	"time"

	"github.com/kbuzsaki/cupid/server"
)

var (
	lock = sync.Mutex{}
	pool = make(map[string]*rpc.Client)
)

type client struct {
	addr           string
	keepAliveDelay time.Duration
}

func NewClient(addr string, keepAliveDelay time.Duration) RPCServer {
	return &client{addr, keepAliveDelay}
}

func connAlive(conn *rpc.Client) bool {
	return conn.Call("Cupid.KeepAlive", nil, nil) == nil
}

func (cl *client) getConn() (*rpc.Client, error) {
	lock.Lock()
	defer lock.Unlock()

	if client, ok := pool[cl.addr]; ok && connAlive(client) {
		return client, nil
	}

	client, err := rpc.Dial("tcp", cl.addr)
	if err != nil {
		return nil, err
	}

	pool[cl.addr] = client
	return client, nil
}

func (cl *client) KeepAlive(args *KeepAliveArgs, events *[]server.Event) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.KeepAlive", args, events)
}

func (cl *client) Open(args *OpenArgs, nd *server.NodeDescriptor) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.Open", args, nd)
}

func (cl *client) Acquire(node server.NodeDescriptor, _ *int) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.Acquire", node, nil)
}

func (cl *client) TryAcquire(node server.NodeDescriptor, success *bool) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.TryAcquire", node, success)
}

func (cl *client) Release(node server.NodeDescriptor, _ *int) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.Release", node, nil)
}

func (cl *client) GetContentAndStat(node server.NodeDescriptor, cas *server.NodeContentAndStat) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.GetContentAndStat", node, cas)
}

func (cl *client) GetStat(node server.NodeDescriptor, stat *server.NodeStat) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.GetStat", node, stat)
}

func (cl *client) SetContent(args *SetContentArgs, success *bool) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.SetContent", args, success)
}
