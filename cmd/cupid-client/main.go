package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"bufio"
	"os"
	"strings"

	"strconv"

	"github.com/kbuzsaki/cupid/client"
	"github.com/kbuzsaki/cupid/server"
)

var (
	addr       = ""
	debug      = false
	command    = ""
	path       = ""
	value      = ""
	generation = uint64(0)
	cl         client.Client
	subscriber client.Subscriber
	handles    = make(map[string]client.NodeHandle)
)

const (
	cmdHelp = "Command List:\n" +
		"\tget <path>\n" +
		"\tset <path> <value> <generation>" +
		"\tlock <name>" +
		"\ttrylock <name>" +
		"\tunlock <name>"
	prompt = "> "
)

func parseArgs() []string {
	addrp := flag.String("addr", "", "the address to connect to")
	debugp := flag.Bool("debug", false, "whether to print event information")
	flag.Parse()

	if *addrp == "" {
		log.Fatal("Address is required")
	}

	addr = *addrp
	debug = *debugp

	args := flag.Args()

	if len(args) >= 1 {
		command = args[0]
		return args[1:]
	}

	return args

}

func maybePrintHelp(ok bool) bool {
	if !ok {
		fmt.Println(cmdHelp)
	}

	return !ok
}

func mustGetNodeHandle(path string) client.NodeHandle {
	if _, ok := handles[path]; !ok {
		nh, err := cl.Open(path, false, server.EventsConfig{})
		if err != nil {
			log.Fatalf("open error:", err)
		}
		handles[path] = nh
	}

	return handles[path]
}

func parseGet(args []string) bool {
	if len(args) == 0 {
		path = ""
	} else {
		path = args[0]
		return true
	}
	return false
}

func handleGet(args []string) bool {
	if maybePrintHelp(parseGet(args)) {
		return false
	}

	nh := mustGetNodeHandle(path)
	cas, err := nh.GetContentAndStat()
	if err != nil {
		log.Fatalf("get error:", err)
	}

	//Maybe we want this switched on a flag?
	fmt.Printf("Generation: %v\n", cas.Stat.Generation)
	fmt.Println(cas.Content)

	return true
}

func parseSet(args []string) bool {
	if len(args) == 0 {
		path = ""
	}

	//check if generation is well formed before messing with path and value
	if len(args) >= 3 {
		gen, err := strconv.ParseUint(args[2], 10, 64)
		if err != nil {
			fmt.Println("Generation expected to be a number")
			return false
		}
		generation = gen
	}

	//generation was good so now parse path and value
	if len(args) >= 2 {
		path = args[0]
		value = args[1]

		return true
	}

	return false
}

func handleSet(args []string) bool {
	if maybePrintHelp(parseSet(args)) {
		return true
	}

	nh := mustGetNodeHandle(path)
	ok, err := nh.SetContent(value, generation)
	if err != nil {
		log.Fatalf("set error:", err)
	}
	if !ok {
		fmt.Println("set no-oped due to generation")
	}

	return true
}

func handleLock(args []string) bool {
	if maybePrintHelp(parseGet(args)) {
		return true
	}

	nh := mustGetNodeHandle(path)

	err := nh.Acquire()
	handles[path] = nh
	if err != nil {
		log.Fatalf("lock error: %v\n", err)
	}

	return true
}

func handleUnlock(args []string) bool {
	if maybePrintHelp(parseGet(args)) {
		return true
	}

	nh := mustGetNodeHandle(path)

	err := nh.Release()
	if err != nil {
		log.Fatalf("lock release error: %v\n", err)
	}

	return true
}

func handleTryLock(args []string) bool {
	if maybePrintHelp(parseGet(args)) {
		return true
	}

	nh := mustGetNodeHandle(path)

	succ, e := nh.TryAcquire()
	if e != nil {
		log.Fatalf("try lock error:%v\n", e)
	}

	if succ {
		fmt.Println("Acquired Lock")
	} else {
		fmt.Println("Failed to Acquire Lock")
	}

	return true
}

func handleSubscribe(args []string) bool {
	if maybePrintHelp(parseGet(args)) {
		return true
	}

	nh := mustGetNodeHandle(path)
	nh.GetContentAndStat()

	subscriber.Register(path, func(path string, cas server.NodeContentAndStat) {
		fmt.Println("Event published for path:", path, "content:", cas)
	})

	fmt.Println("Subscribed!")

	return true
}

func runCmd(args []string) bool {

	switch command {
	case "get":
		return handleGet(args)
	case "lock":
		return handleLock(args)
	case "unlock":
		return handleUnlock(args)
	case "set":
		return handleSet(args)
	case "trylock":
		return handleTryLock(args)
	case "subscribe":
		return handleSubscribe(args)
	case "wait":
		time.Sleep(10 * time.Second)
		return true
	case "exit":
		return false
	default:
		fmt.Printf("Unrecognized command %#v\n", command)
	}

	return false
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)

	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Fields(line)

		if len(args) > 0 {
			command = args[0]
			if len(args) > 0 {
				args = args[1:]
			}

			if !runCmd(args) {
				break
			}

		}
		fmt.Print(prompt)
	}

	err := scanner.Err()
	if err != nil {
		panic(err)
	}

}

func printEvents(in <-chan server.Event) {
	for rawEvent := range in {
		switch event := rawEvent.(type) {
		case server.LockInvalidationEvent:
			fmt.Println("\nLock Invalidation Event:", event)
		case server.ContentInvalidationEvent:
			fmt.Println("\nContent Invalidation Event:", event)
		case server.ContentInvalidationPushEvent:
			fmt.Println("\nContent Invalidation Push Event:", event)
		default:
			fmt.Println("\nUnrecognized event:", rawEvent)
		}
		fmt.Print(prompt)
	}
}

func main() {
	args := parseArgs()

	tmp_cl, err := client.New(addr)
	if err != nil {
		log.Fatalf("error initializing client: %v\n", err)
	}
	cl = tmp_cl

	subscriber, err = client.NewSubscriber(cl)
	if err != nil {
		log.Fatal("failed to open subscriber")
	}

	if command == "" {
		runPrompt()
	} else {
		runCmd(args)
	}

}
