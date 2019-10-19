package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cryptopunkscc/go-xmpp"
)

var port uint
var host string
var vhost string
var from string
var user string
var pass string

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
		fmt.Println("-", id.Local, id.Space)
	}
	fmt.Println()
}

func init() {
	// Configure command-line flags
	flag.UintVar(&port, "port", 5222, "Set custom port")
	flag.StringVar(&vhost, "vhost", "", "Set virtual host to connect to (default same as hostname)")
	flag.StringVar(&from, "from", "", "Set JID to send in the stream header")
	flag.StringVar(&user, "user", "", "Log in with this user")
	flag.StringVar(&pass, "pass", "", "Log in with this password")
	flag.Usage = usage
}

func main() {
	parse() // parse command line options

	conn, err := xmpp.Connect(host, xmpp.JID(host), nil)

	if err != nil {
		fmt.Println("Error connecting to host:", err)
		os.Exit(1)
	}

	fmt.Println("Connected!")
	printFeatures(conn.Features())

	if err := conn.StartTLS(host); err != nil {
		fmt.Println("Failed upgrading to TLS:", err)
		os.Exit(1)
	}

	fmt.Println("Upgraded to TLS.")
	printFeatures(conn.Features())

	if user != "" {
		if err := conn.Authenticate(user, pass); err != nil {
			fmt.Println("Auth error:", err)
			os.Exit(1)
		}

		fmt.Println("Authenticated!")
		printFeatures(conn.Features())

		jid, err := conn.Bind("probe")
		if err != nil {
			fmt.Println("Bind error", err)
		}

		fmt.Println("Bound to", jid)
	}

	if err := conn.Close(); err != nil {
		fmt.Println("Error closing connection:", err)
		os.Exit(1)
	}
}
