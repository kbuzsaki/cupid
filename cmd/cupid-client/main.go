package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kbuzsaki/cupid/client"
	"github.com/kbuzsaki/cupid/server"
)

var (
	addr       = ""
	command    = ""
	path       = ""
	value      = ""
	generation = uint64(0)
)

func parseArgs() {
	index := 0
	addrp := flag.String("addr", "", "the address to connect to")
	generationp := flag.Uint64("gen", 0, "the generation number to restrict a set to")
	flag.Parse()

	if *addrp != "" {
		addr = *addrp
	} else {
		addr = flag.Arg(index)
		index++
	}

	command = flag.Arg(index)
	index++

	path = flag.Arg(index)
	index++

	value = flag.Arg(index)
	index++

	generation = *generationp

	if addr == "" {
		log.Fatal("address required")
	} else if command == "" {
		log.Fatal("command required")
	} else if path == "" {
		log.Fatal("path required")
	}
}

func main() {
	parseArgs()

	cl, err := client.New(addr)
	if err != nil {
		log.Fatalf("error initializing client: %v\n", err)
	}

	switch command {
	case "get":
		nh, err := cl.Open(path, true, server.EventsConfig{})
		if err != nil {
			log.Fatal("get open error:", err)
		}
		cas, err := nh.GetContentAndStat()
		if err != nil {
			log.Fatal("get error:", err)
		}
		fmt.Println(cas.Content)
	case "set":
		nh, err := cl.Open(path, false, server.EventsConfig{})
		if err != nil {
			log.Fatal("set open error:", err)
		}
		ok, err := nh.SetContent(value, generation)
		if err != nil {
			log.Fatal("set error:", err)
		}
		if !ok {
			fmt.Println("set no-oped due to generation")
		}
	default:
		fmt.Printf("Unrecognized command %#v\n", command)
	}
}
