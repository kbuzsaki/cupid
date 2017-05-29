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
		log.Println("Cannot register cupid rpc:", err)
		return
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("RPC server cannot listen:", err)
		ready <- false
		return
	}

	go func() { ready <- true }()

	log.Println("cupid rpc listening on:", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Cannot accept connection:", err)
			break
		}

		log.Println("Accepted new connection from:", conn.RemoteAddr())

		go rpcServer.ServeConn(conn)
	}

	log.Println("RPC Server exited infinite loop")
}
