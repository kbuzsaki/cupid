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

	contents := []string{"some content", "dog", "foo bar", "foo bar"}

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
	if cas.Stat.Generation != 0 {
		t.Error("Generation not initialized to 0")
	}

	oldCas := cas
	for _, content := range contents {
		ok, err := s.SetContent(nd, content, 16)
		ne("Error SetContent /foo/bar:", err)
		if !ok {
			t.Error("Failed write for SetContent")
		}

		cas, err = s.GetContentAndStat(nd)
		ne("Error GetContentAndStat /foo/bar:", err)
		if cas.Content != content {
			t.Errorf("Content not set correctly, was %#v, expected \"some content\"\n", cas.Content)
		}
		if !cas.Stat.LastModified.After(oldCas.Stat.LastModified) {
			t.Error("LastModified not higher for new cas: new:", cas.Stat.LastModified, "old:", oldCas.Stat.LastModified)
		}
		if cas.Stat.Generation <= oldCas.Stat.Generation {
			t.Error("Generation not higher for new cas: new:", cas.Stat.Generation, "old:", oldCas.Stat.Generation)
		}

		stat, err := s.GetStat(nd)
		ne("Error GetStat /foo/bar:", err)
		if cas.Stat != stat {
			t.Error("GetStat did not equal GetContentAndStat")
		}

		oldCas = cas
	}
}
