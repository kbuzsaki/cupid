package client

import (
	"fmt"
	"log"
	"time"

	"github.com/kbuzsaki/cupid/rpcclient"
	"github.com/kbuzsaki/cupid/server"
)

const (
	minimumKeepAliveDelay = 100 * time.Millisecond
)

type clientImpl struct {
	s server.Server
}

// TODO: accept a config file instead?
func New(addr string) (Client, error) {
	s := rpcclient.New(addr)
	return newFromServer(s)
}

func newFromServer(s server.Server) (Client, error) {
	cl := &clientImpl{s: s}
	go cl.keepAlive()
	return cl, nil
}

// background does background KeepAlive processing in a goroutine
func (cl *clientImpl) keepAlive() {
	for {
		before := time.Now()

		events, err := cl.s.KeepAlive(server.LeaseInfo{})
		if err != nil {
			log.Println("KeepAlive error:", err)
		}

		if len(events) > 0 {
			fmt.Println("events:", events)
			// TODO: pass events to callback / goroutine / whatever
			// TODO: cache invalidation
		}

		// before looping around, make sure that it's been at least the minimum amount of time
		// before we start another KeepAlive
		duration := time.Since(before)
		if duration < minimumKeepAliveDelay {
			time.Sleep(minimumKeepAliveDelay - duration)
		}
	}
}

func (cl *clientImpl) Open(path string, readOnly bool, events server.EventsConfig) (NodeHandle, error) {
	nd, err := cl.s.Open(path, readOnly, events)
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
