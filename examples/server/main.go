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
	fmt.Println("TCP connection from", conn.RemoteAddr())
	defer conn.Close()

	// Create a stream over the TCP connection
	stream := xmpp.NewStream(conn)

	// Read stream header from the client
	header, err := stream.ReadHeader()
	if err != nil {
		panic(err)
	}
	fmt.Println("Client requested connection to", header.To)

	// Send back a stream header
	stream.WriteHeader(header.Reply("fakeid"))
	defer stream.Close()

	// Send stream features
	features := &xmpp.Features{}
	features.AddChild(&xmpp.StartTLS{})
	stream.Write(features)

	// Continue reading/writing XMPP messages
	// ...
}

func init() {
	// Configure command-line flags
	flag.StringVar(&bind, "bind", "127.0.0.1", "Bind to this address")
	flag.UintVar(&port, "port", 5222, "Specify port to listen on")
}

func main() {
	// Parse command-line options
	flag.Parse()

	// Start a TCP server
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

		// Pass the connection to a client handler
		go handleClient(conn)
	}
}
