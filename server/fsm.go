package server

import (
	"log"
)

// TODO: does this need a keepalive? where should keepalive information live? maybe just the front end?
type FSM interface {
	OpenSession() SessionDescriptor
	CloseSession(sd SessionDescriptor)
	GetSession(sd SessionDescriptor) *clientSession
	GetSessionDescriptors() []SessionDescriptor

	OpenNode(sd SessionDescriptor, path string, readOnly bool, config EventsConfig) NodeDescriptor
	CloseNode(nd NodeDescriptor)
	GetNodeDescriptor(nd NodeDescriptor) *nodeDescriptor
	GetUnfinalizedNodes() []*nodeInfo

	SetLocked(nd NodeDescriptor)
	ReleaseLock(nd NodeDescriptor) bool

	GetContentAndStat(nd NodeDescriptor) NodeContentAndStat
	PrepareSetContent(nd NodeDescriptor, cas NodeContentAndStat) bool
	FinalizeSetContent(path string)

	Nop(garbage int)
}

type fsmImpl struct {
	sessions *sessionDescriptorMap
	nodes    *nodeInfoMap
}

func NewFSM() (FSM, error) {
	return &fsmImpl{
		sessions: makeSessionDescriptorMap(),
		nodes:    makeNodeInfoMap(),
	}, nil
}

func (fsm *fsmImpl) OpenSession() SessionDescriptor {
	key := fsm.sessions.OpenSession()
	return SessionDescriptor{Descriptor: key}
}

func (fsm *fsmImpl) CloseSession(sd SessionDescriptor) {
	fsm.sessions.CloseSession(sd)
}

func (fsm *fsmImpl) GetSession(sd SessionDescriptor) *clientSession {
	return fsm.sessions.GetSession(sd.Descriptor)
}

func (fsm *fsmImpl) GetSessionDescriptors() []SessionDescriptor {
	var sds []SessionDescriptor

	keys := fsm.sessions.GetSessionDescriptors()
	for _, key := range keys {
		sds = append(sds, SessionDescriptor{key})
	}

	return sds
}

func (fsm *fsmImpl) OpenNode(sd SessionDescriptor, path string, readOnly bool, config EventsConfig) NodeDescriptor {
	session := fsm.sessions.GetSession(sd.Descriptor)

	ni := fsm.nodes.GetOrCreateNode(path)
	key := session.OpenDescriptor(ni, readOnly, config)
	return NodeDescriptor{
		Session:    sd,
		Descriptor: key,
		Path:       path,
	}
}

func (fsm *fsmImpl) CloseNode(nd NodeDescriptor) {
	session := fsm.sessions.GetSession(nd.Session.Descriptor)
	if session == nil {
		log.Println("trying to delete descriptor from null session!")
		return
	}

	session.CloseDescriptor(nd.Descriptor)
}

func (fsm *fsmImpl) GetNodeDescriptor(nd NodeDescriptor) *nodeDescriptor {
	return fsm.sessions.GetDescriptor(nd)
}

func (fsm *fsmImpl) GetUnfinalizedNodes() []*nodeInfo {
	return fsm.nodes.GetUnfinalizedNodes()
}

func (fsm *fsmImpl) SetLocked(nd NodeDescriptor) {
	nid := fsm.sessions.GetDescriptor(nd)
	if nid == nil {
		log.Println("fsm.TryAcquire got invalid node descriptor:", nd)
		return
	}

	nid.ni.SetLocked(nid)
}

func (fsm *fsmImpl) ReleaseLock(nd NodeDescriptor) bool {
	nid := fsm.sessions.GetDescriptor(nd)
	if nid == nil {
		log.Println("fsm.Release got invalid node descriptor:", nd)
		return false
	}

	return nid.ni.Release(nid) == nil
}

func (fsm *fsmImpl) GetContentAndStat(nd NodeDescriptor) NodeContentAndStat {
	nid := fsm.sessions.GetDescriptor(nd)
	if nid == nil {
		log.Println("fsm.GetContentAndStat got invalid node descriptor:", nd)
		return NodeContentAndStat{}
	}

	return nid.ni.GetContentAndStat()
}

func (fsm *fsmImpl) PrepareSetContent(nd NodeDescriptor, cas NodeContentAndStat) bool {
	nid := fsm.sessions.GetDescriptor(nd)
	if nid == nil {
		log.Println("fsm.GetContentAndStat got invalid node descriptor:", nd)
		return false
	}

	// TODO: fix this to use cas.Stat.LastModified?
	return nid.ni.SetContent(cas.Content, cas.Stat.Generation)
}

func (fsm *fsmImpl) FinalizeSetContent(path string) {
	ni := fsm.nodes.GetNode(path)
	if ni == nil {
		log.Println("fsm.GetContentAndStat got invalid node info:", path)
	}

	// TODO: include generation somehow?
	ni.FinalizeSetContent()
}

func (fsm *fsmImpl) Nop(garbage int) {
}
