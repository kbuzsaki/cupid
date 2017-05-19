package server

import (
	"sync"
	"testing"
)

func DoServerTest_KeepAlive(t *testing.T, s Server) {
	err := s.KeepAlive()
	if err != nil {
		t.Error("KeepAlive failed:", err)
	}
}

func DoServerTest_OpenGetSet(t *testing.T, s Server) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	contents := []string{"some content", "dog", "foo bar", "foo bar"}

	nd, err := s.Open("/foo/bar", false, EventsConfig{})
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

func DoServerTest_OpenReadOnly(t *testing.T, s Server) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}
	ae := func(m string, e error) {
		if e == nil {
			t.Error(m)
		}
	}

	nd, err := s.Open("/foo/bar", true, EventsConfig{})
	ne("Error opening /foo/bar:", err)

	cas, err := s.GetContentAndStat(nd)
	ne("Error GetContentAndStat /foo/bar:", err)
	if cas.Content != "" {
		t.Errorf("Default content for new node was %#v, expected empty string\n", cas.Content)
	}
	if cas.Stat.Generation != 0 {
		t.Error("Generation not initialized to 0")
	}

	_, err = s.SetContent(nd, "foo", 16)
	ae("Expected error from SetContent with read only descriptor", err)

	oldCas := cas

	cas, err = s.GetContentAndStat(nd)
	ne("Error GetContentAndStat /foo/bar:", err)
	if cas != oldCas {
		t.Error("Erroneously modified node from read only descriptor")
	}

	stat, err := s.GetStat(nd)
	ne("Error GetStat /foo/bar:", err)
	if cas.Stat != stat {
		t.Error("GetStat did not equal GetContentAndStat")
	}

	err = s.Acquire(nd)
	ae("Expected error from Acquire with read only descriptor", err)

	_, err = s.TryAcquire(nd)
	ae("Expected error from TryAcquire with read only descriptor", err)

	err = s.Release(nd)
	ae("Expected error from Release with read only descriptor", err)
}

func DoServerTest_SetContentGeneration(t *testing.T, s Server) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	nd, err := s.Open("/foo/bar", false, EventsConfig{})
	ne("Error opening /foo/bar:", err)

	cas, err := s.GetContentAndStat(nd)
	ne("Error GetContentAndStat /foo/bar:", err)
	if cas.Content != "" {
		t.Errorf("Default content for new node was %#v, expected empty string\n", cas.Content)
	}
	if cas.Stat.Generation != 0 {
		t.Error("Generation not initialized to 0")
	}

	// when a node is first created, any generation value should succeed:
	ok, err := s.SetContent(nd, "foo", 0)
	ne("Error SetContent generation 0", err)
	if !ok {
		t.Error("Failed initial SetContent with generation 0")
	}

	// after this the generation should be 1
	cas, err = s.GetContentAndStat(nd)
	ne("Error GetContentAndStat after SetContentAndStat", err)
	if cas.Stat.Generation != 1 {
		t.Error("Wrong generation number:", cas.Stat.Generation, "expected 1")
	}
	if cas.Content != "foo" {
		t.Errorf("Wrong content %#v, expected \"foo\"", cas.Content)
	}

	// generation is now 1, so a generation value of 0 should nop
	ok, err = s.SetContent(nd, "bar", 0)
	ne("Error SetContent generation 0 second round", err)
	if ok {
		t.Error("Erroneously succeeded in setting content when generation was too low")
	}

	// confirm that previous set was a nop and content is same
	cas, err = s.GetContentAndStat(nd)
	ne("Error GetContentAndStat after second SetContent", err)
	if cas.Content != "foo" {
		t.Errorf("Wrong content after nop SetContent, was %#v expected \"foo\"", cas.Content)
	}

	// now do a set with the correct minimum generation value, 1
	ok, err = s.SetContent(nd, "bar", 1)
	ne("Error SetContent generation 1", err)
	if !ok {
		t.Error("Failed SetContent with generation 1 (high enough that should not nop)")
	}

	// and confirm that the previous set went through
	cas, err = s.GetContentAndStat(nd)
	ne("Error GetContentAndStat after SetContent with generation 1:", err)
	if cas.Content != "bar" {
		t.Errorf("Wrong content after SetContent with generation 1, was %#v expected \"bar\"", cas.Content)
	}
	if cas.Stat.Generation != 2 {
		t.Error("Wrong generation number:", cas.Stat.Generation, "expected 2")
	}
}

func DoServerTest_ConcurrentOpen(t *testing.T, s Server) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
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

			nd, err := s.Open("/foo/baz", false, EventsConfig{})
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
	nd, err := s.Open("/foo/baz", false, EventsConfig{})
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

func DoServerTest_TryAcquire(t *testing.T, s Server) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	nd, err := s.Open("/foo/bar", false, EventsConfig{})
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
