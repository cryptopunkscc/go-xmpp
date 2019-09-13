package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/cryptopunkscc/go-xmpp"
)

var port uint
var host string
var vhost string
var from string

func usage() {
	fmt.Printf("This tool connects to an XMPP server and prints its stream header and features\n\n")
	fmt.Printf("Usage:\n\n")
	fmt.Printf("  %s [flags] <hostname>\n\n", os.Args[0])
	fmt.Printf("Flags:\n\n")
	flag.PrintDefaults()
}

func init() {
	flag.UintVar(&port, "port", 5222, "Set custom port")
	flag.StringVar(&vhost, "vhost", "", "Set virtual host to connect to (default same as hostname)")
	flag.StringVar(&from, "from", "", "Set JID to send in the stream header")
	flag.Usage = usage
}

func parse() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}
	host = flag.Arg(0)
	if vhost == "" {
		vhost = host
	}
}

func printStreamHeader(h *xmpp.StreamHeader) {
	fmt.Println("- namespace: ", h.Namespace)
	fmt.Println("- version:   ", h.Version)
	fmt.Println("- from:      ", h.From)
	fmt.Println("- to:        ", h.To)
	fmt.Println("- id:        ", h.ID)
}

func printFeatures(features *xmpp.Features) {
	for _, f := range features.Children {
		fmt.Println("-", xmpp.ResolveName(f))
	}
}

func main() {
	parse()

	address := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Connecting to %s...\n\n", address)
	tcp, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	stream, err := xmpp.Open(tcp, xmpp.ClientHeader(from, vhost))
	if err != nil {
		panic(err)
	}

	fmt.Println("Remote header:")
	printStreamHeader(stream.RemoteHeader())

	fmt.Println("\nStream features:")
	printFeatures(stream.Features())
	fmt.Println()
}
