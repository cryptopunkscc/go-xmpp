package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/cryptopunkscc/go-xmpp"
)

var bind string
var port uint

func handleClient(conn net.Conn) {
	fmt.Println("Connection from", conn.RemoteAddr())
	defer conn.Close()

	// Wait for the client to open an XMPP stream
	stream, err := xmpp.Accept(conn)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	// Send back our stream header and features
	header := stream.RemoteHeader().Reply("fakeid")
	stream.WriteStreamHeader(header)
	stream.Write(&xmpp.Features{})
}

func init() {
	flag.StringVar(&bind, "bind", "127.0.0.1", "Bind to this address")
	flag.UintVar(&port, "port", 5222, "Specify port to listen on")
}

func main() {
	flag.Parse()
	address := fmt.Sprintf("%s:%d", bind, port)
	fmt.Println("Starting server at", address)
	srv, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	for {
		// Accept a TCP connection
		conn, err := srv.Accept()
		if err != nil {
			panic(err)
		}

		// Pass the connection to an XMPP handler
		go handleClient(conn)
	}
}
