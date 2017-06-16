package main

import (
	"flag"
	"log"
	"math"
	"sync"
	"time"

	"fmt"

	"strings"

	"github.com/kbuzsaki/cupid/client"
	"github.com/kbuzsaki/cupid/rpcclient"
	"github.com/kbuzsaki/cupid/server"
)

var (
	addrs []string
)

func doPublish(topic string, count int, created *sync.WaitGroup, finished *sync.WaitGroup) {
	cl, err := client.NewRaft(addrs, 5*time.Second)
	if err != nil {
		log.Fatal("Could not create client in doPublish", err)
	}

	nh, err := cl.Open(topic, false, server.EventsConfig{})
	if err != nil {
		log.Fatal("unable to open node handle")
	}
	defer nh.Close()

	created.Done()
	created.Wait()

	start := time.Now()

	for i := 0; i < count; i++ {
		ok, err := nh.SetContent("foo", math.MaxUint64)
		if err != nil || !ok {
			log.Fatal("unable to set content")
		}
	}

	delta := time.Since(start).Nanoseconds()
	fmt.Println(delta)

	ok, err := nh.SetContent("shutdown", math.MaxUint64)
	if err != nil || !ok {
		log.Fatal("unable to set shutdown content")
	}

	finished.Done()
}

func launchPublishers(topic string, numMessages int, numPubs int) {

	created := sync.WaitGroup{}
	finished := sync.WaitGroup{}
	created.Add(numPubs)
	finished.Add(numPubs)

	for i := 0; i < numPubs; i++ {
		go doPublish(topic, numMessages, &created, &finished)
	}

	finished.Wait()
}

func doSubscribe(topic string, created *sync.WaitGroup, finished *sync.WaitGroup) {
	cl, err := client.NewRaft(addrs, 5*time.Second)
	if err != nil {
		log.Fatal("Could not launch new raft in doSubscribe", err)
	}

	nh, err := cl.Open(topic, true, server.EventsConfig{})
	if err != nil {
		log.Fatal("unable to open node handle")
	}
	defer nh.Close()

	created.Done()
	created.Wait()

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
	finished.Done()
}

func launchSubscribers(topic string, count int) {

	created := sync.WaitGroup{}
	finished := sync.WaitGroup{}
	created.Add(count)
	finished.Add(count)

	for i := 0; i < count; i++ {
		go doSubscribe(topic, &created, &finished)
	}
	finished.Wait()
}

func doLocker(topic string, iterations int, created *sync.WaitGroup, finished *sync.WaitGroup) {

	cl, err := client.NewRaft(addrs, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	nh, err := cl.Open(topic, false, server.EventsConfig{})
	if err != nil {
		log.Fatal(err)
	}

	created.Done()
	created.Wait()

	start := time.Now()
	for i := 0; i < iterations; i++ {

		err = nh.Acquire()
		if err != nil {
			log.Fatal("unable to acquire lock")
		}
		err = nh.Release()
		if err != nil {
			log.Fatal("unable to release lock")
		}

	}

	fmt.Println(time.Since(start).Nanoseconds())
	finished.Done()
}

func launchLockers(topic string, numGoRoutines, iterations int) {
	created := sync.WaitGroup{}
	finished := sync.WaitGroup{}
	created.Add(numGoRoutines)
	finished.Add(numGoRoutines)

	for i := 0; i < numGoRoutines; i++ {
		go doLocker(topic, iterations, &created, &finished)
	}

	finished.Wait()
}

func doNoper(topic string, numOps uint64, iterations int, created *sync.WaitGroup, finished *sync.WaitGroup) {
	cl, err := client.NewRaft(addrs, 5*time.Second)
	nh, err := cl.Open(topic, false, server.EventsConfig{})
	if err != nil {
		log.Fatal("Failed to open file in noper")
	}
	created.Done()
	created.Wait()
	start := time.Now()
	for i := 0; i < iterations; i++ {
		nh.Nop(numOps)
	}
	if err != nil {
		log.Fatal("Could not do noper")
	}
	fmt.Println(time.Since(start).Nanoseconds())
	finished.Done()
}

func launchNopers(topic string, numOps uint64, count, iterations int) {
	created := sync.WaitGroup{}
	finished := sync.WaitGroup{}
	created.Add(count)
	finished.Add(count)

	for i := 0; i < count; i++ {
		go doNoper(topic, numOps, iterations, &created, &finished)
	}

	finished.Wait()
}

func launchFailover(topic string, id int) {
	keepAliveDelay := 10 * time.Millisecond
	log.Println("opening rpc")
	s := rpcclient.New(addrs[0], keepAliveDelay)
	log.Println("opening session")
	sd, err := s.OpenSession()
	if err != nil {
		log.Fatal("error opening session:", err)
	}

	log.Println("opening node")
	nd, err := s.Open(sd, topic, true, server.EventsConfig{})
	if err != nil {
		log.Fatal("error opening node:", err)
	}

	s = rpcclient.New(addrs[id], keepAliveDelay)
	for {
		//log.Println("pinging:", addrs[0])
		_, err := s.GetContentAndStat(nd)
		if err != nil {
			fmt.Println(time.Now(), ":", err)
		} else {
			fmt.Println(time.Now(), ":", "success!")
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	addrstrp := flag.String("addrs", "", "the server address")

	publisherp := flag.Bool("publish", false, "whether the client will publish")
	lockp := flag.Bool("locker", false, "whether the client will run locker benchmark")
	nopp := flag.Bool("nop", false, "whether the client will run nops")
	failoverp := flag.Bool("failover", false, "whether the client will do a failover benchmark")

	idp := flag.Int("id", 0, "index to ping")
	pubsp := flag.Int("pubs", 0, "Number of concurrent publishers")
	subsp := flag.Int("subs", 0, "Number of concurrent subscribers")
	topicp := flag.String("topic", "", "the topic to publish or subscribe to")
	numGoRoutinesP := flag.Int("numGoRoutines", 100, "the number of messages to publish")
	iterationsp := flag.Int("iterations", 100, "the number of messages to publish")
	numOpsp := flag.Uint64("numOps", 0, "the number of operations to perform on the server in a single call")

	flag.Parse()

	addrs = strings.Split(*addrstrp, ",")

	/*
		cl, err := client.NewRaft(addrs, 5*time.Second)
		if err != nil {
			log.Fatal("error opening client:", err)
		}
		defer cl.Close()
	*/

	if *publisherp {
		launchSubscribers(*topicp, *subsp)
		launchPublishers(*topicp, *numGoRoutinesP, *pubsp)
	} else if *lockp {
		launchLockers(*topicp, *numGoRoutinesP, *iterationsp)
	} else if *nopp {
		launchNopers(*topicp, *numOpsp, *numGoRoutinesP, *iterationsp)
	} else if *failoverp {
		launchFailover(*topicp, *idp)
	} else { //keep alive test
		launchSubscribers(*topicp, *subsp)
	}
}
