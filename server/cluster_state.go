package server

import (
	"encoding/json"
	"errors"
	"log"
)

var (
	ErrNoLeader = errors.New("no leader")
)

type ClusterState struct {
	IsLeader   bool
	LeaderID   int
	LeaderAddr string
}

func (cs *ClusterState) MakeRedirectError() error {
	if cs.LeaderAddr == "" {
		log.Println("returned no leader error")
		return ErrNoLeader
	}
	log.Println("returned redirect to:", cs.LeaderAddr)
	return LeaderRedirectError{cs.LeaderID, cs.LeaderAddr}
}

type LeaderRedirectError struct {
	LeaderID   int
	LeaderAddr string
}

func (lre LeaderRedirectError) Error() string {
	b, err := json.Marshal(&lre)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
