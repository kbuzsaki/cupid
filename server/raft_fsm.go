package server

import (
	"bytes"
	"encoding/gob"
	"log"
	"sync/atomic"
)

const (
	openSessionProposalType = iota
	closeSessionProposalType
	openNodeProposalType
	closeNodeProposalType
	tryAcquireProposalType
	releaseProposalType
	prepareSetContentProposalType
	finalizeSetContentProposalType
)

type Proposal struct {
	Type int
	*OpenSessionProposal
	*CloseSessionProposal
	*OpenNodeProposal
	*CloseNodeProposal
	*TryAcquireProposal
	*ReleaseProposal
	*PrepareSetContentProposal
	*FinalizeSetContentProposal
}

func (p *Proposal) Get() interface{} {
	switch p.Type {
	case openSessionProposalType:
		return *p.OpenSessionProposal
	case closeSessionProposalType:
		return *p.CloseSessionProposal
	case openNodeProposalType:
		return *p.OpenNodeProposal
	case closeNodeProposalType:
		return *p.CloseNodeProposal
	case tryAcquireProposalType:
		return *p.TryAcquireProposal
	case releaseProposalType:
		return *p.ReleaseProposal
	case prepareSetContentProposalType:
		return *p.PrepareSetContentProposal
	case finalizeSetContentProposalType:
		return *p.FinalizeSetContentProposal
	default:
		return nil
	}
}

type OpenSessionProposal struct {
	ID uint64
}

func (osp *OpenSessionProposal) Wrap() Proposal {
	return Proposal{Type: openSessionProposalType, OpenSessionProposal: osp}
}

type CloseSessionProposal struct {
	ID uint64

	SD SessionDescriptor
}

func (csp *CloseSessionProposal) Wrap() Proposal {
	return Proposal{Type: closeSessionProposalType, CloseSessionProposal: csp}
}

type OpenNodeProposal struct {
	ID uint64

	SD       SessionDescriptor
	Path     string
	ReadOnly bool
	Config   EventsConfig
}

func (onp *OpenNodeProposal) Wrap() Proposal {
	return Proposal{Type: openNodeProposalType, OpenNodeProposal: onp}
}

type CloseNodeProposal struct {
	ID uint64

	ND NodeDescriptor
}

func (cnp *CloseNodeProposal) Wrap() Proposal {
	return Proposal{Type: closeNodeProposalType, CloseNodeProposal: cnp}
}

type TryAcquireProposal struct {
	ID uint64
	ND NodeDescriptor
}

func (tap *TryAcquireProposal) Wrap() Proposal {
	return Proposal{Type: tryAcquireProposalType, TryAcquireProposal: tap}
}

type ReleaseProposal struct {
	ID uint64
	ND NodeDescriptor
}

func (rp *ReleaseProposal) Wrap() Proposal {
	return Proposal{Type: releaseProposalType, ReleaseProposal: rp}
}

type PrepareSetContentProposal struct {
	ID  uint64
	ND  NodeDescriptor
	CAS NodeContentAndStat
}

func (scp *PrepareSetContentProposal) Wrap() Proposal {
	return Proposal{Type: prepareSetContentProposalType, PrepareSetContentProposal: scp}
}

type FinalizeSetContentProposal struct {
	ID uint64
	ND NodeDescriptor
}

func (fcp *FinalizeSetContentProposal) Wrap() Proposal {
	return Proposal{Type: finalizeSetContentProposalType, FinalizeSetContentProposal: fcp}
}

func Encode(proposal Proposal) string {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&proposal); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func Decode(s string) interface{} {
	var p Proposal
	if err := gob.NewDecoder(bytes.NewBufferString(s)).Decode(&p); err != nil {
		return nil
	}
	return p.Get()
}

func NewRaftFSM(proposeC chan<- string, committedC <-chan *string, delegate FSM) FSM {
	fsm := &raftFSMImpl{
		proposeC:               proposeC,
		committedC:             committedC,
		delegate:               delegate,
		id:                     0,
		openSessionAcks:        NewAtomicMap(),
		closeSessionAcks:       NewAtomicMap(),
		openNodeAcks:           NewAtomicMap(),
		closeNodeAcks:          NewAtomicMap(),
		tryAcquireAcks:         NewAtomicMap(),
		releaseAcks:            NewAtomicMap(),
		setContentAcks:         NewAtomicMap(),
		finalizeSetContentAcks: NewAtomicMap(),
	}

	go fsm.readFromLog()

	return fsm
}

type raftFSMImpl struct {
	proposeC   chan<- string
	committedC <-chan *string
	delegate   FSM

	id uint64

	openSessionAcks        AtomicMap
	closeSessionAcks       AtomicMap
	openNodeAcks           AtomicMap
	closeNodeAcks          AtomicMap
	tryAcquireAcks         AtomicMap
	releaseAcks            AtomicMap
	setContentAcks         AtomicMap
	finalizeSetContentAcks AtomicMap
}

func (fsm *raftFSMImpl) nextId() uint64 {
	return atomic.AddUint64(&fsm.id, 1)
}

func (fsm *raftFSMImpl) OpenSession() SessionDescriptor {
	id := fsm.nextId()

	ac := make(chan SessionDescriptor)
	fsm.openSessionAcks.Put(id, ac)

	proposal := OpenSessionProposal{ID: id}
	fsm.proposeC <- Encode(proposal.Wrap())

	return <-ac
}

func (fsm *raftFSMImpl) CloseSession(sd SessionDescriptor) {
	id := fsm.nextId()

	ac := make(chan bool)
	fsm.closeSessionAcks.Put(id, ac)

	proposal := CloseSessionProposal{ID: id, SD: sd}
	fsm.proposeC <- Encode(proposal.Wrap())

	<-ac
}

func (fsm *raftFSMImpl) GetSession(sd SessionDescriptor) *clientSession {
	return fsm.delegate.GetSession(sd)
}

func (fsm *raftFSMImpl) GetSessionDescriptors() []SessionDescriptor {
	return fsm.delegate.GetSessionDescriptors()
}

func (fsm *raftFSMImpl) OpenNode(sd SessionDescriptor, path string, readOnly bool, config EventsConfig) NodeDescriptor {
	id := fsm.nextId()

	ac := make(chan NodeDescriptor)
	fsm.openNodeAcks.Put(id, ac)

	proposal := OpenNodeProposal{ID: id, SD: sd, Path: path, ReadOnly: readOnly, Config: config}
	fsm.proposeC <- Encode(proposal.Wrap())

	return <-ac
}

func (fsm *raftFSMImpl) CloseNode(nd NodeDescriptor) {
	id := fsm.nextId()

	ac := make(chan bool)
	fsm.closeNodeAcks.Put(id, ac)

	proposal := CloseNodeProposal{ID: id, ND: nd}
	fsm.proposeC <- Encode(proposal.Wrap())

	<-ac
}

func (fsm *raftFSMImpl) GetNodeDescriptor(nd NodeDescriptor) *nodeDescriptor {
	return fsm.delegate.GetNodeDescriptor(nd)
}

func (fsm *raftFSMImpl) GetUnfinalizedNodes() []*nodeInfo {
	return fsm.delegate.GetUnfinalizedNodes()
}

func (fsm *raftFSMImpl) SetLocked(nd NodeDescriptor) {
	id := fsm.nextId()

	ac := make(chan bool)
	fsm.tryAcquireAcks.Put(id, ac)

	proposal := TryAcquireProposal{ID: id, ND: nd}
	fsm.proposeC <- Encode(proposal.Wrap())

	<-ac
}

func (fsm *raftFSMImpl) ReleaseLock(nd NodeDescriptor) bool {
	id := fsm.nextId()

	ac := make(chan bool)
	fsm.releaseAcks.Put(id, ac)

	proposal := ReleaseProposal{ID: id, ND: nd}
	fsm.proposeC <- Encode(proposal.Wrap())

	return <-ac
}

func (fsm *raftFSMImpl) GetContentAndStat(nd NodeDescriptor) NodeContentAndStat {
	return fsm.delegate.GetContentAndStat(nd)
}

func (fsm *raftFSMImpl) PrepareSetContent(nd NodeDescriptor, cas NodeContentAndStat) bool {
	id := fsm.nextId()

	ac := make(chan bool)
	fsm.setContentAcks.Put(id, ac)

	proposal := PrepareSetContentProposal{ID: id, ND: nd, CAS: cas}
	fsm.proposeC <- Encode(proposal.Wrap())

	return <-ac
}

func (fsm *raftFSMImpl) FinalizeSetContent(nd NodeDescriptor) {
	id := fsm.nextId()

	ac := make(chan bool)
	fsm.finalizeSetContentAcks.Put(id, ac)

	proposal := FinalizeSetContentProposal{ID: id, ND: nd}
	fsm.proposeC <- Encode(proposal.Wrap())

	<-ac
}

func (fsm *raftFSMImpl) readFromLog() {
	for operation := range fsm.committedC {
		if operation == nil {
			log.Println("nil operation")
			continue
		}

		proposal := Decode(*operation)
		switch p := proposal.(type) {
		case OpenSessionProposal:
			sd := fsm.delegate.OpenSession()
			if ch := fsm.openSessionAcks.Get(p.ID); ch != nil {
				ch.(chan SessionDescriptor) <- sd
			}
		case CloseSessionProposal:
			fsm.delegate.CloseSession(p.SD)
			if ch := fsm.closeSessionAcks.Get(p.ID); ch != nil {
				ch.(chan bool) <- true
			}
		case OpenNodeProposal:
			nd := fsm.delegate.OpenNode(p.SD, p.Path, p.ReadOnly, p.Config)
			if ch := fsm.openNodeAcks.Get(p.ID); ch != nil {
				ch.(chan NodeDescriptor) <- nd
			}
		case CloseNodeProposal:
			fsm.delegate.CloseNode(p.ND)
			if ch := fsm.closeNodeAcks.Get(p.ID); ch != nil {
				ch.(chan bool) <- true
			}
		case TryAcquireProposal:
			fsm.delegate.SetLocked(p.ND)
			if ch := fsm.tryAcquireAcks.Get(p.ID); ch != nil {
				ch.(chan bool) <- true
			}
		case ReleaseProposal:
			succ := fsm.delegate.ReleaseLock(p.ND)
			if ch := fsm.releaseAcks.Get(p.ID); ch != nil {
				ch.(chan bool) <- succ
			}
		case PrepareSetContentProposal:
			succ := fsm.delegate.PrepareSetContent(p.ND, p.CAS)
			if ch := fsm.setContentAcks.Get(p.ID); ch != nil {
				ch.(chan bool) <- succ
			}
		case FinalizeSetContentProposal:
			fsm.delegate.FinalizeSetContent(p.ND)
			if ch := fsm.finalizeSetContentAcks.Get(p.ID); ch != nil {
				ch.(chan bool) <- true
			}
		default:
			log.Println("unrecognized operation:", proposal)
		}
	}
}
