package rpcclient

import (
	"log"
	"net"
	"net/rpc"

	"github.com/kbuzsaki/cupid/server"
)

func ServeCupidRPC(s server.Server, addr string, ready chan bool) {
	rpcServer := rpc.NewServer()
	cupidRPC := NewServer(s)
	err := rpcServer.RegisterName("Cupid", cupidRPC)

	if err != nil {
		log.Printf("Cannot register cupid rpc %#v", err)
		return
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("RPC server cannot listen: %#v\n", err)
		ready <- false
		return
	}

	go func() { ready <- true }()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Cannot accept connection: %#v\n", err)
			break
		}

		go rpcServer.ServeConn(conn)
	}

	log.Printf("RPC Server exited infinite loop")
}
