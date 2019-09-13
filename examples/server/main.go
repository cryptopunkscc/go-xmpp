package main

import (
	"fmt"
	"net"

	"github.com/cryptopunkscc/go-xmpp"
)

func acceptHeader(header *xmpp.StreamHeader) *xmpp.StreamHeader {
	if header.To != "example.com" {
		return nil
	}

	return &xmpp.StreamHeader{
		To:        header.From,
		From:      "example.com",
		Version:   "1.0",
		Namespace: xmpp.NamespaceClient,
		ID:        "fakeid",
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:5222")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		s, err := xmpp.AcceptFunc(conn, acceptHeader)
		if err != nil {
			panic(err)
		}
		s.Write(&xmpp.Features{})
		fmt.Println("Accepted new stream", s.RemoteHeader().From)
		s.Close()
		conn.Close()
	}
}
