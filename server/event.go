package server

import "encoding/gob"

func init() {
	gob.Register(LockInvalidationEvent{})
	gob.Register(ContentInvalidationEvent{})
	gob.Register(ContentInvalidationPushEvent{})
}

type EventsConfig struct {
	ContentModified bool
	LockInvalidated bool
	MasterFailed    bool
}

type Event interface {
}

type LockInvalidationEvent struct {
	Descriptor NodeDescriptor
}

type ContentInvalidationEvent struct {
	Descriptor NodeDescriptor
}

type ContentInvalidationPushEvent struct {
	Descriptor NodeDescriptor
	NodeContentAndStat
}
