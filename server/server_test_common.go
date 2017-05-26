package server

import (
	"sync"
	"testing"
)

func DoServerTest_KeepAlive(t *testing.T, s Server) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	events, err := s.KeepAlive(LeaseInfo{}, nil, 1)
	ne("KeepAlive with Empty LeaseInfo", err)
	if len(events) != 0 {
		t.Error("Got events, expected none: ", events)
	}

	// test with LeaseInfo that claims to have a lock we don't have
	nd, err := s.Open("/foo/bar", false, EventsConfig{})
	ne("Error opening /foo/bar:", err)

	// expect to get lock invalidation event
	bogusLeaseInfo := LeaseInfo{[]NodeDescriptor{nd}}
	events, err = s.KeepAlive(bogusLeaseInfo, nil, 1)
	ne("KeepAlive with bogus LeaseInfo", err)
	if len(events) != 1 {
		t.Error("Got incorrect number of events, expected 1: ", events)
	} else {
		if event, ok := events[0].(LockInvalidationEvent); !ok {
			t.Error("Expected LockInvalidationEvent, got:", events[0])
		} else {
			if event.Descriptor != nd {
				t.Error("Got wrong Descriptor, expected:", nd, ", was:", event.Descriptor)
			}
		}
	}

	// test case where we do have a lock
	ok, err := s.TryAcquire(nd)
	ne("TryAcquire", err)
	if !ok {
		t.Error("Failed to acquire lock")
	}

	goodLeaseInfo := LeaseInfo{[]NodeDescriptor{nd}}
	events, err = s.KeepAlive(goodLeaseInfo, nil, 1)
	ne("KeepAlive with good LeaseInfo", err)
	if len(events) != 0 {
		t.Error("Got events with good LeaseInfo, expected none: ", events)
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
	ae("Expected error from SetContent with read only Descriptor", err)

	oldCas := cas

	cas, err = s.GetContentAndStat(nd)
	ne("Error GetContentAndStat /foo/bar:", err)
	if cas != oldCas {
		t.Error("Erroneously modified node from read only Descriptor")
	}

	err = s.Acquire(nd)
	ae("Expected error from Acquire with read only Descriptor", err)

	_, err = s.TryAcquire(nd)
	ae("Expected error from TryAcquire with read only Descriptor", err)

	err = s.Release(nd)
	ae("Expected error from Release with read only Descriptor", err)
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

	err = s.Release(nd)
	if err == nil {
		t.Error("Failed to error when releasing a lock not held")
	}

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

func DoServerTest_BadRelease(t *testing.T, s Server) {
	ne := func(m string, e error) {
		if e != nil {
			t.Error(m, e)
		}
	}

	nd1, err := s.Open("/foo/bar", false, EventsConfig{})
	ne("Error opening /foo/bar nd1:", err)

	nd2, err := s.Open("/foo/bar", false, EventsConfig{})
	ne("Error opening /foo/bar nd2:", err)

	ok, err := s.TryAcquire(nd1)
	ne("Acquiring nd1", err)
	if !ok {
		t.Error("Unable to acquire lock")
	}

	err = s.Release(nd2)
	if err == nil {
		t.Error("Erroneously released lock that we do not own")
	}
}
