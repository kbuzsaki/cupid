package server

import "testing"

func TestServerImpl_KeepAlive(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	err = s.KeepAlive()
	if err != nil {
		t.Error("KeepAlive failed:", err)
	}
}

func TestServerImpl_OpenGetSet(t *testing.T) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	s, err := New()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	nd, err := s.Open("/foo/bar", false)
	ne("Error opening /foo/bar:", err)

	cas, err := s.GetContentAndStat(nd)
	ne("Error GetContentAndStat /foo/bar:", err)
	if cas.Content != "" {
		t.Errorf("Default content for new node was %#v, expected empty string\n", cas.Content)
	}
	// TODO: check generation and other stat info

	ok, err := s.SetContent(nd, "some content", 16)
	ne("Error SetContent /foo/bar:", err)
	if !ok {
		t.Error("Failed write for SetContent")
	}

	cas, err = s.GetContentAndStat(nd)
	ne("Error GetContentAndStat /foo/bar:", err)
	if cas.Content != "some content" {
		t.Errorf("Content not set correctly, was %#v, expected \"some content\"\n", cas.Content)
	}
}
