package rpcclient

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/kbuzsaki/cupid/server"
)

func randaddr() string {
	n := rand.Intn(10000) + 10001
	s := "127.0.0.1:" + strconv.Itoa(n)
	return s
}

func TestRPC_KeepAlive(t *testing.T) {

	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	s, err := server.New()
	ne("Could not instantiate server", err)
	addr := randaddr()

	ready := make(chan bool)
	go ServeCupidRPC(s, addr, ready)
	v := <-ready
	if !v {
		t.Fatal("Could not launch rpc server")
	}

	cl := New(addr, 1)

	server.DoServerTest_KeepAlive(t, cl)
}

func TestRPC_OpenGetSet(t *testing.T) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	s, err := server.New()
	ne("Could not instantiate server", err)
	addr := randaddr()

	ready := make(chan bool)
	go ServeCupidRPC(s, addr, ready)
	v := <-ready
	if !v {
		t.Fatal("Could not launch rpc server")
	}

	cl := New(addr, 1)

	server.DoServerTest_OpenGetSet(t, cl)
}

func TestRPC_OpenReadOnly(t *testing.T) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	s, err := server.New()
	ne("Could not instantiate server", err)
	addr := randaddr()

	ready := make(chan bool)
	go ServeCupidRPC(s, addr, ready)
	v := <-ready
	if !v {
		t.Fatal("Could not launch rpc server")
	}

	cl := New(addr, 1)

	server.DoServerTest_OpenReadOnly(t, cl)
}

func TestRPC_SetContentGeneration(t *testing.T) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	s, err := server.New()
	ne("Could not instantiate server", err)
	addr := randaddr()

	ready := make(chan bool)
	go ServeCupidRPC(s, addr, ready)
	v := <-ready
	if !v {
		t.Fatal("Could not launch rpc server")
	}

	cl := New(addr, 1)

	server.DoServerTest_SetContentGeneration(t, cl)
}

func TestRPC_ConcurrentOpen(t *testing.T) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	s, err := server.New()
	ne("Could not instantiate server", err)
	addr := randaddr()

	ready := make(chan bool)
	go ServeCupidRPC(s, addr, ready)
	v := <-ready
	if !v {
		t.Fatal("Could not launch rpc server")
	}

	cl := New(addr, 1)

	server.DoServerTest_ConcurrentOpen(t, cl)
}

func TestRPC_TryAcquire(t *testing.T) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	s, err := server.New()
	ne("Could not instantiate server", err)
	addr := randaddr()

	ready := make(chan bool)
	go ServeCupidRPC(s, addr, ready)
	v := <-ready
	if !v {
		t.Fatal("Could not launch rpc server")
	}

	cl := New(addr, 1)

	server.DoServerTest_TryAcquire(t, cl)
}

func TestRPC_BadRelease(t *testing.T) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	s, err := server.New()
	ne("Could not instantiate server", err)
	addr := randaddr()

	ready := make(chan bool)
	go ServeCupidRPC(s, addr, ready)
	v := <-ready
	if !v {
		t.Fatal("Could not launch rpc server")
	}

	cl := New(addr, 1)

	server.DoServerTest_BadRelease(t, cl)
}
