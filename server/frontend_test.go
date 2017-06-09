package server

import (
	"testing"
	"time"
)

// These test functions just immediately delegate to the shared tester functions in server_test_common.go

func TestFrontend_KeepAlive(t *testing.T) {
	s, err := NewFrontend()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_KeepAlive(t, s)
}

func TestFrontend_OpenGetSet(t *testing.T) {
	s, err := NewFrontend()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_OpenGetSet(t, s)
}

func TestFrontend_OpenReadOnly(t *testing.T) {
	s, err := NewFrontend()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_OpenReadOnly(t, s)
}

func TestFrontend_SetContentGeneration(t *testing.T) {
	s, err := NewFrontend()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_SetContentGeneration(t, s)
}

func TestFrontend_ConcurrentOpen(t *testing.T) {
	s, err := NewFrontend()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_ConcurrentOpen(t, s)
}

func TestFrontend_TryAcquire(t *testing.T) {
	s, err := NewFrontend()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_TryAcquire(t, s)
}

func TestFrontend_BadRelease(t *testing.T) {
	s, err := NewFrontend()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_BadRelease(t, s)
}

func TestFrontendImpl_SetContentFailover(t *testing.T) {
	fsm, err := NewFSM()
	if err != nil {
		t.Fatal("unable to create fsm:", err)
	}

	// grab a session and node descriptor
	sd := fsm.OpenSession()
	nd := fsm.OpenNode(sd, "/foo", false, EventsConfig{ContentModified: true})

	nid := fsm.GetNodeDescriptor(nd)

	// add an incomplete SetContent to run failover for
	fsm.PrepareSetContent(nd, NodeContentAndStat{Content: "new set", Stat: NodeStat{Generation: 10, LastModified: time.Now()}})

	stateC := make(chan ClusterState, 1)
	stateC <- ClusterState{true, 1, ""}
	s, err := NewFrontendWithFSM(fsm, stateC)
	if err != nil {
		t.Fatal("unable to create frontend with fsm:", err)
	}

	// check in for events
	events, err := s.KeepAlive(LeaseInfo{Session: sd}, []EventInfo{{nd, 0, true}}, time.Second)
	if len(events) != 1 {
		t.Error("events slice was not unary:", events)
	} else {
		if e, ok := events[0].(ContentInvalidationPushEvent); !ok {
			t.Errorf("event was not content invalidation push, was: %#v", events[0])
		} else {
			if e.Descriptor != nd {
				t.Error("wrong descriptor:", e.Descriptor, "expected:", nd)
			}
			if e.Content != "new set" {
				t.Error("wrong content:", e.Content, "expected: 'new set'")
			}
		}
	}

	ok, err := s.SetContent(nd, "final set", 10)
	if err != nil || !ok {
		t.Error("failed to do final set")
	}

	if !nid.ni.finalized {
		t.Error("not finalized")
	}
}
