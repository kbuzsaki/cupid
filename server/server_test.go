package server

import "testing"

// These test functions just immediately delegate to the shared tester functions in server_test_common.go

func TestServerImpl_KeepAlive(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_KeepAlive(t, s)
}

func TestServerImpl_OpenGetSet(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_OpenGetSet(t, s)
}

func TestServerImpl_OpenReadOnly(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_OpenReadOnly(t, s)
}

func TestServerImpl_SetContentGeneration(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_SetContentGeneration(t, s)
}

func TestServerImpl_ConcurrentOpen(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_ConcurrentOpen(t, s)
}

func TestServerImpl_TryAcquire(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	DoServerTest_TryAcquire(t, s)
}
