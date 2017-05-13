package server

import (
	"sync"
	"testing"
)

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

func TestServerImpl_SetContentGeneration(t *testing.T) {
	// TODO: test skipping writes for old generation
}

func TestServerImpl_ConcurrentOpen(t *testing.T) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	s, err := New()
	if err != nil {
		t.Fatal("Unable to start server:", err)
	}

	content := "some content"
	mainDone1 := sync.WaitGroup{}
	mainDone2 := sync.WaitGroup{}
	childrenDone := sync.WaitGroup{}
	childCount := 100

	mainDone1.Add(1)
	for i := 0; i < childCount; i++ {
		go func() {
			// wait until all children are created to do concurrent open
			mainDone1.Wait()

			nd, err := s.Open("/foo/baz", false)
			ne("Error opening file in child:", err)

			// signal that this child is done and wait until main is done
			childrenDone.Done()
			mainDone2.Wait()

			// verify that the node descriptor from earlier reads the correct data
			cas, err := s.GetContentAndStat(nd)
			ne("Error getting content in child:", err)

			if cas.Content != content {
				t.Errorf("Read incorrect content in child: %#v\n", cas.Content)
			}

			// finally tell main that this child is done
			childrenDone.Done()
		}()
	}

	// once all of the children are created, signal them to open and wait until they all do
	childrenDone.Add(childCount)
	mainDone2.Add(1)
	mainDone1.Done()
	childrenDone.Wait()

	// set the content so that all of the children can read it
	nd, err := s.Open("/foo/baz", false)
	ne("Error opening file in main:", err)

	ok, err := s.SetContent(nd, content, 10)
	ne("Error setting content from main:", err)
	if !ok {
		t.Error("Failed to set content from main")
	}

	// then wait for all of the children to finish reading
	childrenDone.Add(childCount)
	mainDone2.Done()
	childrenDone.Wait()
}

func TestServerImpl_TryAcquire(t *testing.T) {
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

	ok, err := s.TryAcquire(nd)
	ne("Error TryAcquire:", err)
	if !ok {
		t.Error("Failed to acquire lock with TryAcquire starting from default state")
	}

	ok, err = s.TryAcquire(nd)
	ne("Error TryAcquire after TryAcquire:", err)
	if ok {
		t.Error("Acquired lock for a second time")
	}

	err = s.Release(nd)
	ne("Error Release:", err)

	ok, err = s.TryAcquire(nd)
	ne("Error TryAcquire after Release:", err)
	if !ok {
		t.Error("Failed to acquire lock after Release")
	}

	ok, err = s.TryAcquire(nd)
	ne("Error TryAcquire after TryAcquire:", err)
	if ok {
		t.Error("Acquired lock for a second time")
	}
}
