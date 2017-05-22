package server

import "encoding/gob"

func init() {
	gob.Register(LockInvalidationEvent{})
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
