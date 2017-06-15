package client

import (
	"sync"
	"time"

	"encoding/json"
	"log"
	"net/rpc"

	"github.com/kbuzsaki/cupid/server"
)

type RedirectServer struct {
	delegates []server.Server

	lock          sync.RWMutex
	leader        int
	pendingLeader int
}

func unmarshalRedirectError(se string) (server.LeaderRedirectError, error) {
	var lre server.LeaderRedirectError
	err := json.Unmarshal([]byte(se), &lre)
	return lre, err
}

func (rs *RedirectServer) getLeader() server.Server {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	// stable leader
	if rs.leader != 0 {
		return rs.delegates[rs.leader-1]
	}

	// leader not known, try pending leader
	return rs.delegates[rs.pendingLeader-1]
}

func (rs *RedirectServer) stabilizeLeader() {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	// only stabilize if we're currently aborting
	if rs.leader == 0 {
		rs.leader = rs.pendingLeader
		//log.Println("stabilize leader:", rs.leader)
	}
}

func (rs *RedirectServer) abortLeader() {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	//log.Println("aborting leader:", rs.leader, rs.pendingLeader)

	if rs.leader == 0 {
		if rs.pendingLeader < len(rs.delegates) {
			rs.pendingLeader++
		} else {
			// totally partitioned, maybe retry again and eventually go into jeopardy?
			panic("jeopardy!")
		}
	} else {
		rs.leader = 0
		rs.pendingLeader = 1
	}

	//log.Println("new pending leader:", rs.pendingLeader)
}

// setLeader sets the known stable leader
func (rs *RedirectServer) setLeader(leader int) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.leader = leader
}

func (rs *RedirectServer) KeepAlive(li server.LeaseInfo, eis []server.EventInfo, keepAliveDelay time.Duration) ([]server.Event, error) {
	events, err := rs.getLeader().KeepAlive(li, eis, keepAliveDelay)
	if err == nil {
		rs.stabilizeLeader()
		return events, nil
	} else if se, ok := err.(rpc.ServerError); ok {
		if lre, err := unmarshalRedirectError(se.Error()); err == nil {
			rs.setLeader(lre.LeaderID)
			return rs.KeepAlive(li, eis, keepAliveDelay)
		}
		return events, err
	} else {
		// it's an error from the rpc library
		rs.abortLeader()
		return rs.KeepAlive(li, eis, keepAliveDelay)
	}
}

func (rs *RedirectServer) OpenSession() (server.SessionDescriptor, error) {
	sd, err := rs.getLeader().OpenSession()
	if err == nil {
		rs.stabilizeLeader()
		return sd, nil
	} else if se, ok := err.(rpc.ServerError); ok {
		if lre, err := unmarshalRedirectError(se.Error()); err == nil {
			rs.setLeader(lre.LeaderID)
			return rs.OpenSession()
		}
		return sd, err
	} else {
		rs.abortLeader()
		return rs.OpenSession()
	}
}

func (rs *RedirectServer) CloseSession(sd server.SessionDescriptor) error {
	err := rs.getLeader().CloseSession(sd)
	if err == nil {
		rs.stabilizeLeader()
		return nil
	} else if se, ok := err.(rpc.ServerError); ok {
		if lre, err := unmarshalRedirectError(se.Error()); err == nil {
			rs.setLeader(lre.LeaderID)
			return rs.CloseSession(sd)
		}
		log.Println("server error:", se)
		return err
	} else {
		rs.abortLeader()
		return rs.CloseSession(sd)
	}
}

func (rs *RedirectServer) Open(sd server.SessionDescriptor, path string, readOnly bool, config server.EventsConfig) (server.NodeDescriptor, error) {
	nd, err := rs.getLeader().Open(sd, path, readOnly, config)
	if err == nil {
		rs.stabilizeLeader()
		return nd, nil
	} else if se, ok := err.(rpc.ServerError); ok {
		if lre, err := unmarshalRedirectError(se.Error()); err == nil {
			rs.setLeader(lre.LeaderID)
			return rs.Open(sd, path, readOnly, config)
		}
		log.Println("server error:", se)
		return nd, err
	} else {
		rs.abortLeader()
		return rs.Open(sd, path, readOnly, config)
	}
}

func (rs *RedirectServer) CloseNode(nd server.NodeDescriptor) error {
	err := rs.getLeader().CloseNode(nd)
	if err == nil {
		rs.stabilizeLeader()
		return nil
	} else if se, ok := err.(rpc.ServerError); ok {
		if lre, err := unmarshalRedirectError(se.Error()); err == nil {
			rs.setLeader(lre.LeaderID)
			return rs.CloseNode(nd)
		}
		log.Println("server error:", se)
		return err
	} else {
		rs.abortLeader()
		return rs.CloseNode(nd)
	}
}

func (rs *RedirectServer) Acquire(node server.NodeDescriptor) error {
	err := rs.getLeader().Acquire(node)
	if err == nil {
		rs.stabilizeLeader()
		return nil
	} else if se, ok := err.(rpc.ServerError); ok {
		if lre, err := unmarshalRedirectError(se.Error()); err == nil {
			rs.setLeader(lre.LeaderID)
			return rs.Acquire(node)
		}
		log.Println("server error:", se)
		return err
	} else {
		rs.abortLeader()
		return rs.Acquire(node)
	}
}

func (rs *RedirectServer) TryAcquire(node server.NodeDescriptor) (bool, error) {
	ok, err := rs.getLeader().TryAcquire(node)
	if err == nil {
		rs.stabilizeLeader()
		return ok, nil
	} else if se, ok := err.(rpc.ServerError); ok {
		if lre, err := unmarshalRedirectError(se.Error()); err == nil {
			rs.setLeader(lre.LeaderID)
			return rs.TryAcquire(node)
		}
		log.Println("server error:", se)
		return ok, err
	} else {
		rs.abortLeader()
		return rs.TryAcquire(node)
	}
	return ok, err
}

func (rs *RedirectServer) Release(node server.NodeDescriptor) error {
	err := rs.getLeader().Release(node)
	if err == nil {
		rs.stabilizeLeader()
		return nil
	} else if se, ok := err.(rpc.ServerError); ok {
		if lre, err := unmarshalRedirectError(se.Error()); err == nil {
			rs.setLeader(lre.LeaderID)
			return rs.Release(node)
		}
		log.Println("server error:", se)
		return err
	} else {
		rs.abortLeader()
		return rs.Release(node)
	}
	return err
}

func (rs *RedirectServer) GetContentAndStat(node server.NodeDescriptor) (server.NodeContentAndStat, error) {
	cas, err := rs.getLeader().GetContentAndStat(node)
	if err == nil {
		rs.stabilizeLeader()
		return cas, nil
	} else if se, ok := err.(rpc.ServerError); ok {
		if lre, err := unmarshalRedirectError(se.Error()); err == nil {
			rs.setLeader(lre.LeaderID)
			return rs.GetContentAndStat(node)
		}
		log.Println("server error:", se)
		return cas, err
	} else {
		rs.abortLeader()
		return rs.GetContentAndStat(node)
	}
	return cas, err
}

func (rs *RedirectServer) SetContent(node server.NodeDescriptor, content string, generation uint64) (bool, error) {
	ok, err := rs.getLeader().SetContent(node, content, generation)
	if err == nil {
		rs.stabilizeLeader()
		return ok, nil
	} else if se, ok := err.(rpc.ServerError); ok {
		if lre, err := unmarshalRedirectError(se.Error()); err == nil {
			rs.setLeader(lre.LeaderID)
			return rs.SetContent(node, content, generation)
		}
		log.Println("server error:", se)
		return ok, err
	} else {
		rs.abortLeader()
		return rs.SetContent(node, content, generation)
	}
	return ok, err
}

func (rs *RedirectServer) Nop(numOps uint64) error {
	err := rs.getLeader().Nop(numOps)
	if err == nil {
		rs.stabilizeLeader()
		return nil
	} else if se, ok := err.(rpc.ServerError); ok {
		if lre, err := unmarshalRedirectError(se.Error()); err == nil {
			rs.setLeader(lre.LeaderID)
			return rs.Nop(numOps)
		}
		return err
	} else {
		rs.abortLeader()
		return rs.Nop(numOps)
	}

}
