package server

import (
	"bytes"
	"encoding/gob"
	"log"
	"sync/atomic"
)

const (
	openSessionProposalType = iota
	openNodeProposalType    = iota
	tryAcquireProposalType  = iota
	releaseProposalType     = iota
	setContentProposalType  = iota
)

type Proposal struct {
	Type int
	*OpenSessionProposal
	*OpenNodeProposal
	*TryAcquireProposal
	*ReleaseProposal
	*SetContentProposal
}

func (p *Proposal) Get() interface{} {
	switch p.Type {
	case openSessionProposalType:
		return *p.OpenSessionProposal
	case openNodeProposalType:
		return *p.OpenNodeProposal
	case tryAcquireProposalType:
		return *p.TryAcquireProposal
	case releaseProposalType:
		return *p.ReleaseProposal
	case setContentProposalType:
		return *p.SetContentProposal
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

type OpenNodeProposal struct {
	ID uint64

	SD       SessionDescriptor
	Path     string
	ReadOnly bool
}

func (onp *OpenNodeProposal) Wrap() Proposal {
	return Proposal{Type: openNodeProposalType, OpenNodeProposal: onp}
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

type SetContentProposal struct {
	ID  uint64
	ND  NodeDescriptor
	CAS NodeContentAndStat
}

func (scp *SetContentProposal) Wrap() Proposal {
	return Proposal{Type: setContentProposalType, SetContentProposal: scp}
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
		proposeC:        proposeC,
		committedC:      committedC,
		delegate:        delegate,
		id:              0,
		openSessionAcks: NewAtomicMap(),
		openNodeAcks:    NewAtomicMap(),
		tryAcquireAcks:  NewAtomicMap(),
		releaseAcks:     NewAtomicMap(),
		setContentAcks:  NewAtomicMap(),
	}

	go fsm.readFromLog()

	return fsm
}

type raftFSMImpl struct {
	proposeC   chan<- string
	committedC <-chan *string
	delegate   FSM

	id uint64

	openSessionAcks AtomicMap
	openNodeAcks    AtomicMap
	tryAcquireAcks  AtomicMap
	releaseAcks     AtomicMap
	setContentAcks  AtomicMap
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

func (fsm *raftFSMImpl) GetSession(sd SessionDescriptor) *clientSession {
	return fsm.delegate.GetSession(sd)
}

func (fsm *raftFSMImpl) GetSessionDescriptors() []SessionDescriptor {
	return fsm.delegate.GetSessionDescriptors()
}

func (fsm *raftFSMImpl) OpenNode(sd SessionDescriptor, path string, readOnly bool) NodeDescriptor {
	id := fsm.nextId()

	ac := make(chan NodeDescriptor)
	fsm.openNodeAcks.Put(id, ac)

	proposal := OpenNodeProposal{ID: id, SD: sd, Path: path, ReadOnly: readOnly}
	fsm.proposeC <- Encode(proposal.Wrap())

	return <-ac
}

func (fsm *raftFSMImpl) GetNodeDescriptor(nd NodeDescriptor) *nodeDescriptor {
	return fsm.delegate.GetNodeDescriptor(nd)
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

func (fsm *raftFSMImpl) SetContent(nd NodeDescriptor, cas NodeContentAndStat) bool {
	id := fsm.nextId()

	ac := make(chan bool)
	fsm.setContentAcks.Put(id, ac)

	proposal := SetContentProposal{ID: id, ND: nd, CAS: cas}
	fsm.proposeC <- Encode(proposal.Wrap())

	return <-ac
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
		case OpenNodeProposal:
			nd := fsm.delegate.OpenNode(p.SD, p.Path, p.ReadOnly)
			if ch := fsm.openNodeAcks.Get(p.ID); ch != nil {
				ch.(chan NodeDescriptor) <- nd
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
		case SetContentProposal:
			succ := fsm.delegate.SetContent(p.ND, p.CAS)
			if ch := fsm.setContentAcks.Get(p.ID); ch != nil {
				ch.(chan bool) <- succ
			}
		default:
			log.Println("unrecognized operation:", proposal)
		}
	}
}
