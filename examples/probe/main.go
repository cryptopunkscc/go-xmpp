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

// usage displays help
func usage() {
	fmt.Printf("This tool connects to an XMPP server and prints its stream header and features\n\n")
	fmt.Printf("Usage:\n\n")
	fmt.Printf("  %s [flags] <hostname>\n\n", os.Args[0])
	fmt.Printf("Flags:\n\n")
	flag.PrintDefaults()
}

// parse checks and parses command-line options
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
	fmt.Println("Stream header:")
	fmt.Println("- namespace: ", h.Namespace)
	fmt.Println("- version:   ", h.Version)
	fmt.Println("- from:      ", h.From)
	fmt.Println("- to:        ", h.To)
	fmt.Println("- id:        ", h.ID)
	fmt.Println()
}

func printFeatures(features *xmpp.Features) {
	fmt.Println("Stream features:")
	for _, f := range features.Children {
		id := xmpp.Identify(f)
		fmt.Println("-", id.Local)
	}
	fmt.Println()
}

func init() {
	// Configure command-line flags
	flag.UintVar(&port, "port", 5222, "Set custom port")
	flag.StringVar(&vhost, "vhost", "", "Set virtual host to connect to (default same as hostname)")
	flag.StringVar(&from, "from", "", "Set JID to send in the stream header")
	flag.Usage = usage
}

func main() {
	parse() // parse command line options

	// Establish a TCP connection
	address := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Connecting to %s...\n\n", address)
	tcp, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	// Start an XMPP stream
	stream := xmpp.NewStream(tcp)

	err = stream.WriteHeader(xmpp.NewHeader(xmpp.NamespaceClient, xmpp.JID(from), xmpp.JID(vhost)))
	if err != nil {
		panic(err)
	}

	// Read stream header from the server
	header, err := stream.ReadHeader()
	if err != nil {
		panic(err)
	}
	printStreamHeader(header)

	// Read stream features
	features, err := stream.ReadFeatures()
	if err != nil {
		panic(err)
	}
	printFeatures(features)

	// Ready to exchange XMPP messages!

	stream.Close()
}
