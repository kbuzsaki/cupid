package main

import (
	"fmt"
	"log"

	"github.com/kbuzsaki/cupid/client"
)

func main() {
	_, err := client.New("localhost:12345")
	if err != nil {
		log.Printf("error initializing client: %v\n", err)
	}

	fmt.Println("Hello, client")
}
