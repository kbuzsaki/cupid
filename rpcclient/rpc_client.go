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
	var a int
	var b int
	return conn.Call("Cupid.Ping", &a, &b) == nil
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

func (cl *client) Ping(_, _ *int) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}
	var a int
	var b int

	return conn.Call("Cupid.Ping", &a, &b)
}

func (cl *client) KeepAlive(args *KeepAliveArgs, events *[]server.Event) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.KeepAlive", args, events)
}

func (cl *client) OpenSession(_ int, sd *server.SessionDescriptor) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.OpenSession", 0, sd)
}

func (cl *client) CloseSession(sd *server.SessionDescriptor, _ *int) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	i := 0
	return conn.Call("Cupid.CloseSession", sd, &i)
}

func (cl *client) Open(args *OpenArgs, nd *server.NodeDescriptor) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.Open", args, nd)
}

func (cl *client) CloseNode(nd *server.NodeDescriptor, _ *int) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	i := 0
	return conn.Call("Cupid.CloseNode", nd, &i)
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

func (cl *client) SetContent(args *SetContentArgs, success *bool) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.SetContent", args, success)
}

func (cl *client) Nop(numOps uint64, garbage *bool) error {
	conn, err := cl.getConn()
	if err != nil {
		return err
	}

	return conn.Call("Cupid.Nop", numOps, garbage)
}
