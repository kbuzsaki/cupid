package main

import (
	"log"

	"github.com/kbuzsaki/cupid/rpcclient"
	"github.com/kbuzsaki/cupid/server"
)

func main() {
	var s server.Server
	s, err := server.New()

	if err != nil {
		log.Printf("error initializing server: %v\n", err)
	}

	ready := make(chan bool)
	rpcclient.ServeCupidRPC(s, "localhost:3000", ready)
}
