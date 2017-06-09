package server

import (
	"math"
	"testing"
)

func BenchmarkFsmImpl_SetContent(b *testing.B) {
	fsm, err := NewFSM()
	if err != nil {
		b.Fatal("unable to create fsm")
	}

	sd := fsm.OpenSession()
	nd := fsm.OpenNode(sd, "/foo/bar", false, EventsConfig{})

	cas := NodeContentAndStat{
		Content: "some content",
		Stat: NodeStat{
			Generation: math.MaxUint64,
		},
	}

	// initial set content to check correctness
	fsm.PrepareSetContent(nd, cas)
	getCas := fsm.GetContentAndStat(nd)
	if cas.Content != getCas.Content {
		b.Error("set content failed")
	}

	cas.Content = "some other content"

	for n := 0; n < b.N; n++ {
		fsm.PrepareSetContent(nd, cas)
	}

	getCas2 := fsm.GetContentAndStat(nd)
	if cas.Content != getCas2.Content {
		b.Error("looped set content failed")
	}
}

func BenchmarkFsmImpl_SetContentNoOp(b *testing.B) {
	fsm, err := NewFSM()
	if err != nil {
		b.Fatal("unable to create fsm")
	}

	sd := fsm.OpenSession()
	nd := fsm.OpenNode(sd, "/foo/bar", false, EventsConfig{})

	cas := NodeContentAndStat{
		Content: "some content",
		Stat: NodeStat{
			Generation: 0,
		},
	}

	// initial set content to check correctness
	fsm.PrepareSetContent(nd, cas)
	getCas := fsm.GetContentAndStat(nd)
	if cas.Content != getCas.Content {
		b.Error("set content failed")
	}

	cas.Content = "some other content"

	for n := 0; n < b.N; n++ {
		fsm.PrepareSetContent(nd, cas)
	}

	getCas2 := fsm.GetContentAndStat(nd)
	if "some content" != getCas2.Content {
		b.Error("looped set content failed")
	}
}
