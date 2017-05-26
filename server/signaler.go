package server

import "sync"

type Signaler interface {
	Signal()
	Reset()
	DoneChan() <-chan struct{}
}

type signalerImpl struct {
	lock sync.Mutex
	c    chan struct{}
	done bool
}

func NewSignaler() Signaler {
	return &signalerImpl{c: make(chan struct{})}
}

func (s *signalerImpl) Signal() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.done {
		close(s.c)
		s.done = true
	}
}

func (s *signalerImpl) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.c = make(chan struct{})
	s.done = false
}

func (s *signalerImpl) DoneChan() <-chan struct{} {
	return s.c
}
