package client

import "testing"

// TODO: do dependency injection here with mocks so that the client doesn't talk to an actual server
func TestClientImpl_Open(t *testing.T) {
	cl, err := New("localhost:12345")
	if err != nil {
		t.Fatal("Unable to open client:", err)
	}

	_, err = cl.Open("/foo/bar")
	if err == nil {
		t.Error("Erroneously opened non-existant file")
	}
}
