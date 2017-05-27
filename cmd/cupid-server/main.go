package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/kbuzsaki/cupid/rpcclient"
	"github.com/kbuzsaki/cupid/server"
)

var (
	addr = ""
)

func parseArgs() {
	index := 0
	addrp := flag.String("addr", "", "the address to listen on")
	flag.Parse()

	if *addrp != "" {
		addr = *addrp
	} else {
		addr = flag.Arg(index)
		index++
	}

	if addr == "" {
		log.Fatal("addr required")
	}
}

func main() {
	parseArgs()

	go http.ListenAndServe("localhost:6060", nil)

	s, err := server.New()
	if err != nil {
		log.Fatalf("error initializing server: %v\n", err)
	}

	log.Println("starting cupid-server on", addr)
	ready := make(chan bool)
	rpcclient.ServeCupidRPC(s, addr, ready)
}
