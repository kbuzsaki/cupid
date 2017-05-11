package main

import (
	"fmt"
	"log"

	"github.com/kbuzsaki/cupid/server"
)

func main() {
	var s server.Server
	s, err := server.New("localhost:12345")
	if err != nil {
		log.Printf("error initializing server: %v\n", err)
	}

	_ = s.KeepAlive()

	fmt.Println("Hello, server")
}
