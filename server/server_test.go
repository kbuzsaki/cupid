package server

import "testing"

func TestServer_KeepAlive(t *testing.T) {
	s, err := New("localhost:12345")
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	err = s.KeepAlive()
	if err != nil {
		t.Error("KeepAlive failed:", err)
	}
}
