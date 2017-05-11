package main

import (
	"fmt"
	"log"

	"github.com/kbuzsaki/cupid/client"
)

func main() {
	var cl client.Client
	cl, err := client.New("localhost:12345")
	if err != nil {
		log.Printf("error initializing client: %v\n", err)
	}

	_, _ = cl.Open("/foo/bar")

	fmt.Println("Hello, client")
}
