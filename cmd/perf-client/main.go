package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"math"

	"sync"

	"github.com/kbuzsaki/cupid/client"
	"github.com/kbuzsaki/cupid/server"
)

func doPublish(cl client.Client, topic string, count int) {
	nh, err := cl.Open(topic, false, server.EventsConfig{})
	if err != nil {
		log.Fatal("unable to open node handle")
	}

	for i := 0; i < count; i++ {
		ok, err := nh.SetContent("foo", math.MaxUint64)
		if err != nil || !ok {
			log.Fatal("unable to set content")
		}
	}

	ok, err := nh.SetContent("shutdown", math.MaxUint64)
	if err != nil || !ok {
		log.Fatal("unable to set shutdown content")
	}
}

func doSubscribe(cl client.Client, topic string) {
	nh, err := cl.Open(topic, true, server.EventsConfig{})
	if err != nil {
		log.Fatal("unable to open node handle")
	}

	_, err = nh.GetContentAndStat()
	if err != nil {
		log.Fatal("unable to getcas")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	// TODO: get event for if we somehow lose connection?
	nh.Register(func(path string, cas server.NodeContentAndStat) {
		// NO OP
		if cas.Content == "shutdown" {
			wg.Done()
		}
	})

	wg.Wait()
	time.Sleep(1 * time.Second)
}

func main() {
	addrstrp := flag.String("addrs", "", "the server address")
	publisherp := flag.Bool("publish", false, "whether the client will publish or subscribe")
	topicp := flag.String("topic", "", "the topic to publish or subscribe to")
	countp := flag.Int("count", 100, "the number of messages to publish")
	flag.Parse()

	addrs := strings.Split(*addrstrp, ",")
	//fmt.Println("addresses:", addrs)

	cl, err := client.NewRaft(addrs, 5*time.Second)
	if err != nil {
		log.Fatal("error opening client:", err)
	}

	if *publisherp {
		doPublish(cl, *topicp, *countp)
	} else {
		doSubscribe(cl, *topicp)
	}
}
