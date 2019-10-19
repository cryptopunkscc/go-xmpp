package xmpp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

// Conn represents an XMPP connection
type Conn struct {
	transport    io.ReadWriteCloser
	stream       *Stream
	localHeader  *StreamHeader
	remoteHeader *StreamHeader
	features     *Features
	logger       Logger
	mu           sync.Mutex
}

type Logger interface {
	Received([]byte)
	Sent([]byte)
}

// Connect establishes an XMPP connection
func Connect(addr string, to JID, logger Logger) (*Conn, error) {
	var err error

	c := &Conn{logger: logger}
	c.localHeader = &StreamHeader{
		Namespace: NamespaceClient,
		To:        to,
		Version:   "1.0",
	}

	finalAddr := resolveSRV(addr, "client")

	// If SRV resolution failed, try the defaults
	if finalAddr == "" {
		finalAddr = fmt.Sprintf("%s:%d", addr, defaultClientPort)
	}

	// Establish a TCP connection
	tcp, err := net.Dial("tcp", finalAddr)
	if err != nil {
		return nil, err
	}

	// Establish an XMPP stream over the TCP connection
	if err := c.RestartStream(tcp); err != nil {
		return nil, err
	}

	return c, nil
}

// StartTLS upgrades the connection to TLS if possible and replaces the stream
func (c *Conn) StartTLS(serverName string) error {
	tcp, ok := c.Transport().(net.Conn)
	if !ok {
		return errors.New("tcp transport required to start TLS")
	}
	if c.Features().ChildCount(&StartTLS{}) == 0 {
		return errors.New("stream doesn't support TLS")
	}

	// Start TLS negotiation
	if err := c.Write(&StartTLS{}); err != nil {
		return err
	}
	msg, err := c.Read()
	if err != nil {
		return err
	}

	if _, ok := msg.(*Proceed); !ok {
		id := Identify(msg)
		return fmt.Errorf("unexpected response while upgrading to TLS: %s", id.Local)
	}

	tlsConn := tls.Client(tcp, &tls.Config{ServerName: serverName})
	if tlsConn == nil {
		return errors.New("tls failed")
	}
	return c.RestartStream(tlsConn)
}

// Bind binds the XMPP stream to a resource name
func (c *Conn) Bind(resourceName string) (JID, error) {
	req := &IQ{
		Type: "set",
		ID:   "bind-request",
	}
	req.AddChild(&Bind{
		Resource: resourceName,
	})
	if err := c.Write(req); err != nil {
		return "", err
	}
	msg, err := c.Read()
	if err != nil {
		return "", err
	}
	res, ok := msg.(*IQ)
	if !ok {
		return "", fmt.Errorf("bind: unexpected message: %v", Identify(msg))
	}
	if !res.Result() {
		return "", errors.New("unexpected iq response: invalid type attribute")
	}
	if res.ID != req.ID {
		return "", errors.New("unexpected iq response: invalid id attribute")
	}
	bind, ok := res.Child(&Bind{}).(*Bind)
	if !ok {
		return "", errors.New("unexpected iq response: unexpected element type")
	}
	return bind.JID, nil
}

// Transport returns the transport used for the stream
func (c *Conn) Transport() io.ReadWriteCloser {
	return c.transport
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
func (c *Conn) Close() error {
	if err := c.stream.Close(); err != nil {
		return err
	}
	return c.transport.Close()
}

// Features returns the current stream features
func (c *Conn) Features() *Features {
	return c.features
}

// RestartStream restarts the current stream. If transport is provided, it replaces the currently used transport.
func (c *Conn) RestartStream(transport io.ReadWriteCloser) (err error) {
	if transport != nil {
		c.transport = transport
	}
	c.stream = NewStream(c.transportWithLogger())
	err = c.stream.WriteHeader(c.localHeader)
	if err != nil {
		return err
	}
	c.remoteHeader, err = c.stream.ReadHeader()
	if err != nil {
		return err
	}
	return c.readFeatures()
}

// transportWithLogger returns the current transport wrapped with a logger
func (c *Conn) transportWithLogger() io.ReadWriter {
	if c.logger == nil {
		return c.transport
	}
	return &tee{
		target: c.transport,
		logger: c.logger,
	}
}

// readFeatures reads features from the stream and stores them in the struct
func (c *Conn) readFeatures() error {
	msg, err := c.Read()
	if err != nil {
		return err
	}
	if feats, ok := msg.(*Features); ok {
		c.features = feats
		return nil
	}
	id := Identify(msg)
	return fmt.Errorf("ReadFeatures: unexpected message: %s", id.Local)
}
