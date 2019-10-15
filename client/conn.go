package xmppc

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/cryptopunkscc/go-xmpp"
)

// Conn represents an XMPP connection
type Conn struct {
	jid xmpp.JID

	// Internal state
	stream       *xmpp.Stream
	localHeader  *xmpp.StreamHeader
	remoteHeader *xmpp.StreamHeader
	features     *xmpp.Features
	mu           sync.Mutex
}

// Credentials holds authentication information
type Credentials struct {
	Username string
	Password string
}

// Connect establishes a connection to the XMPP server
func Connect(host string, jid xmpp.JID, password string) (*Conn, error) {
	var err error

	c := &Conn{}

	c.localHeader = &xmpp.StreamHeader{
		Namespace: xmpp.NamespaceClient,
		To:        jid.Domain(),
		Version:   "1.0",
	}
	// If no host was provided, extract it from JID
	if host == "" {
		host = jid.Domain().String()
	}

	// Add default port if none is specified
	if !strings.Contains(host, ":") {
		host = host + ":5222"
	}

	// Establish a TCP connection
	tcp, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	// Establish an XMPP stream over the TCP connection
	c.stream = xmpp.NewStream(tcp)
	if err := c.openStream(); err != nil {
		return nil, err
	}

	// Establish a TLS socket over the XMPP stream
	err = c.upgradeToTLS(jid.Domain().String())
	if err != nil {
		return nil, err
	}

	// Establish an XMPP stream over the TLS socket
	if err := c.openStream(); err != nil {
		return nil, err
	}

	// Authenticate via SASL over the XMPP stream
	err = c.authenticate(jid.Local(), password)
	if err != nil {
		return nil, err
	}

	// Reestablish an XMPP stream over the TLS socket
	c.stream = xmpp.NewStream(c.stream.Transport())
	if err := c.openStream(); err != nil {
		return nil, err
	}

	// Bind
	err = c.bind(jid.Resource())
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Read reads the next XMPP message from the stream
func (c *Conn) Read() (interface{}, error) {
	return c.stream.Read()
}

// Write writes an XMPP message to the stream
func (c *Conn) Write(msg interface{}) error {
	return c.stream.Write(msg)
}

// Close closes the XMPP connection
func (c *Conn) Close() {
	c.stream.Close()
}

// JID returns JID the connection is bound to
func (c *Conn) JID() xmpp.JID {
	return c.jid
}

// Features returns the current stream features
func (c *Conn) Features() *xmpp.Features {
	return c.features
}

// openStream opens a bidirectional stream
func (c *Conn) openStream() (err error) {
	err = c.stream.WriteHeader(c.localHeader)
	if err != nil {
		return
	}
	c.remoteHeader, err = c.stream.ReadHeader()
	if err != nil {
		return
	}
	c.features, err = c.stream.ReadFeatures()
	return
}

// upgradeToTLS establishes a TLS session over an XMPP stream
func (c *Conn) upgradeToTLS(serverName string) error {
	tcp, ok := c.stream.Transport().(net.Conn)
	if !ok {
		return errors.New("tcp transport required to start TLS")
	}
	if c.Features().ChildCount(&xmpp.StartTLS{}) == 0 {
		return errors.New("stream doesn't support TLS")
	}

	// Start TLS negotiation
	if err := c.Write(&xmpp.StartTLS{}); err != nil {
		return err
	}
	msg, err := c.Read()
	if err != nil {
		return err
	}

	if _, ok := msg.(*xmpp.Proceed); !ok {
		id := xmpp.Identify(msg)
		return fmt.Errorf("unexpected response while upgrading to TLS: %s", id.Local)
	}

	tlsConn := tls.Client(tcp, &tls.Config{ServerName: serverName})
	if tlsConn == nil {
		return errors.New("tls failed")
	}
	c.stream = xmpp.NewStream(tlsConn)
	return nil
}

// bind binds the XMPP stream to a resource name
func (c *Conn) bind(resourceName string) error {
	req := &xmpp.IQ{
		Type: "set",
		From: c.jid,
		To:   "",
		ID:   "bind-request",
		Lang: "en",
	}
	req.AddChild(&xmpp.Bind{
		Resource: resourceName,
	})

	if err := c.stream.Write(req); err != nil {
		panic(err)
	}

	msg, err := c.stream.Read()
	if err != nil {
		return err
	}

	res, ok := msg.(*xmpp.IQ)
	if !ok {
		return fmt.Errorf("bind: unexpected message: %v", xmpp.Identify(msg))
	}
	if !res.Result() {
		return errors.New("unexpected iq response: invalid type attribute")
	}
	if res.ID != req.ID {
		return errors.New("unexpected iq response: invalid id attribute")
	}
	bind, ok := res.Child(&xmpp.Bind{}).(*xmpp.Bind)
	if !ok {
		return errors.New("unexpected iq response: unexpected element type")
	}
	c.jid = bind.JID

	return nil
}
