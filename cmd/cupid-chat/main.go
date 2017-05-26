package main

import (
	"flag"
	"log"

	"fmt"
	"strings"
	"sync"

	"github.com/kbuzsaki/cupid/client"
	"github.com/kbuzsaki/cupid/server"
	"bufio"
	"os"
	"math"
)

var (
	addr    = ""
	channel = ""
	nick    = ""
)

func parseArgs() {
	flag.StringVar(&addr, "addr", "", "the server address")
	flag.StringVar(&channel, "channel", "", "the channel name to join")
	flag.StringVar(&nick, "nick", "", "the nickname to use")
	flag.Parse()

	index := 0
	if addr == "" {
		addr = flag.Arg(index)
		index++
	}
	if channel == "" {
		channel = flag.Arg(index)
		index++
	}
	if nick == "" {
		channel = flag.Arg(index)
	}

	if addr == "" || channel == "" || nick == "" {
		log.Fatal("addr, channel, and nick required")
	}
}

func advertiseNick(chanHandle client.NodeHandle, nick string) ([]string, error) {
	err := chanHandle.Acquire()
	if err != nil {
		return nil, err
	}
	defer chanHandle.Release()

	oldNicks := ""
	done := false
	for !done {
		cas, err := chanHandle.GetContentAndStat()
		if err != nil {
			return nil, err
		}

		oldNicks = cas.Content
		newNicks := oldNicks + "\n" + nick
		done, err = chanHandle.SetContent(newNicks, cas.Stat.Generation+1)
		if err != nil {
			return nil, err
		}
	}

	return strings.Split(oldNicks, "\n"), nil
}

type message struct {
	sender string
	body   string
}

type Channel struct {
	cl       client.Client
	messages chan message

	lock     sync.Mutex
	chatters map[string]struct{}
}

func (c *Channel) printMessages() {
	for m := range c.messages {
		fmt.Println(m.sender + ": " + m.body + "\n> ")
	}
}

func (c *Channel) registerChatter(chatter string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.chatters[chatter]; ok {
		return
	}
	c.chatters[chatter] = struct{}{}


	chatterHandle, err := c.cl.Open(channel + "/" + chatter, true, server.EventsConfig{})
	if err != nil {
		log.Fatal("unable to open chatter handle:", err)
	}

	chatterHandle.Register(func(path string, cas server.NodeContentAndStat) {
		pathParts := strings.SplitN(path, "/", 2)
		if len(pathParts) != 2 {
			log.Println("got bad message path: ", path)
		}

		sender := pathParts[1]
		m := message{sender, cas.Content}
		c.messages <- m
	})
}

func main() {
	parseArgs()

	cl, err := client.New(addr)
	if err != nil {
		log.Fatal("error opening client:", err)
	}

	messages := make(chan message, 100)
	ch := Channel{cl: cl, messages: messages, chatters: make(map[string]struct{})}
	go ch.printMessages()

	chanHandle, err := cl.Open(channel, false, server.EventsConfig{})
	if err != nil {
		log.Fatal("error opening channel:", err)
	}
	chatters, err := advertiseNick(chanHandle, nick)
	if err != nil {
		log.Fatal("unable to advertise nick")
	}

	for _, chatter := range chatters {
		ch.registerChatter(chatter)
	}

	chanHandle.Register(func(path string, cas server.NodeContentAndStat) {
		chatters := strings.Split(cas.Content, "\n")
		for _, chatter := range chatters {
			ch.registerChatter(chatter)
		}
	})

	nickPath := channel + "/" + nick
	nickHandle, err := cl.Open(nickPath, false, server.EventsConfig{})
	if err != nil {
		log.Fatal("error opening nick:", err)
	}


	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("error reading from stdin:", err)
		}

		if line == "" {
			continue
		}

		nickHandle.SetContent(line, math.MaxUint32)
	}
}
