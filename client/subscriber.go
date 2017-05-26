package client

import (
	"errors"
	"log"

	"github.com/kbuzsaki/cupid/server"
)

var (
	ErrInvalidClient = errors.New("invalid client")
)

type SubscriberCallback func(path string, cas server.NodeContentAndStat)

type Subscriber interface {
	Register(path string, cb SubscriberCallback)
}

type subscriber struct {
	cl        Client
	callbacks map[string]SubscriberCallback
}

func NewSubscriber(cl Client) (*subscriber, error) {
	if cl == nil {
		return nil, ErrInvalidClient
	}

	s := &subscriber{cl, make(map[string]SubscriberCallback)}

	go s.handleEvents()

	return s, nil
}

func (s *subscriber) Register(path string, cb SubscriberCallback) {
	s.callbacks[path] = cb
}

func (s *subscriber) handleEvents() {
	for event := range s.cl.GetEventsOut() {
		//log.Println("got event:", event)
		s.handleEvent(event)
	}
}

func (s *subscriber) handleEvent(rawEvent server.Event) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered from:", r)
		}
	}()

	switch event := rawEvent.(type) {
	case server.ContentInvalidationEvent:
		nh := &nodeHandleImpl{s.cl.(*clientImpl), event.Descriptor}
		cas, err := nh.GetContentAndStat()
		if err != nil {
			return
		}

		//log.Println("got invalidation event for nd:", event.Descriptor)
		cb := s.callbacks[event.Descriptor.Path()]
		if cb != nil {
			cb(event.Descriptor.Path(), cas)
		}
	case server.ContentInvalidationPushEvent:
		cb := s.callbacks[event.Descriptor.Path()]
		//log.Println("got push event for nd:", event.Descriptor)
		if cb != nil {
			cb(event.Descriptor.Path(), event.NodeContentAndStat)
		}
	}
}
