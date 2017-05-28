package server

import "testing"

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
