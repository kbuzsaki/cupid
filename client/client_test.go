package client

import (
	"testing"

	"errors"

	"github.com/kbuzsaki/cupid/mocks"
	"github.com/kbuzsaki/cupid/server"
)

// TODO: do dependency injection here with mocks so that the client doesn't talk to an actual server
func TestClientImpl_Open(t *testing.T) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}
	mockServer := &mocks.Server{}
	sd := server.SessionDescriptor{3}
	cl := clientImpl{s: mockServer, sd: sd}

	// success case
	nd := server.NodeDescriptor{sd, 4, "/foo/bar"}
	mockServer.On("Open", sd, "/foo/bar", false, server.EventsConfig{}).Return(nd, nil)

	// test success
	nh, err := cl.Open("/foo/bar", false, server.EventsConfig{})
	ne("Opening /foo/bar from client", err)
	if nh == nil {
		t.Error("Got nil NodeHandle from Open")
	}
	mockServer.AssertExpectations(t)

	// failure case
	someError := errors.New("some error")
	mockServer.On("Open", sd, "/bad/file", false, server.EventsConfig{}).Return(server.NodeDescriptor{}, someError)

	nh, err = cl.Open("/bad/file", false, server.EventsConfig{})
	if err == nil {
		t.Error("Client performed bad open")
	}
	mockServer.AssertExpectations(t)
}
