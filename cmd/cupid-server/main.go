package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/pkg/capnslog"
	"github.com/kbuzsaki/cupid/rpcclient"
	"github.com/kbuzsaki/cupid/server"
)

func main() {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	cupidport := flag.Int("port", 9121, "cupid rpc server port")
	join := flag.Bool("join", false, "join an existing cluster")
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	flag.Parse()

	if !*verbose {
		capnslog.SetGlobalLogLevel(capnslog.ERROR)
	}

	go http.ListenAndServe(fmt.Sprintf("localhost:%d", *cupidport+1), nil)

	proposeC := make(chan string)
	defer close(proposeC)
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)

	// raft provides a commit stream for the proposals from the http api
	fsm, err := server.NewFSM()
	if err != nil {
		log.Fatal("unable to open fsm:", err)
	}

	getSnapshot := func() ([]byte, error) {
		log.Println("no snapshotting available")
		return nil, nil
	}
	commitC, errorC, stateC, snapshotterReady := newRaftNode(*id, strings.Split(*cluster, ","), *join, getSnapshot, proposeC, confChangeC)

	// TODO: what to do with these things?
	_ = errorC
	_ = snapshotterReady

	raftFSM := server.NewRaftFSM(proposeC, commitC, fsm)
	s, err := server.NewFrontendWithFSM(raftFSM, stateC)
	if err != nil {
		log.Fatalf("error initializing server: %v\n", err)
	}

	cupidaddr := fmt.Sprintf("localhost:%d", *cupidport)
	log.Println("starting cupid-server on", cupidaddr)
	ready := make(chan bool)
	rpcclient.ServeCupidRPC(s, cupidaddr, ready)
}
