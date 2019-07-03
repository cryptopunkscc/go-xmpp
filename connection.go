package xmpp

import (
	"encoding/xml"
	"errors"
	"io"
	"sync"
)

// Conn manages a bidirectional XMPP stream using an underlying transport
type Conn struct {
	transport  io.ReadWriter
	decoder    *xml.Decoder
	encoder    *xml.Encoder
	piReceived bool
	streamEnd  xml.EndElement

	features *Features

	outStream *StreamHeader
	inStream  *StreamHeader

	sync.Mutex
}

// Open a new bidirectional XMPP stream on the provided transport
func Open(transport io.ReadWriter, outStream *StreamHeader) (*Conn, error) {
	var err error

	conn := &Conn{
		transport: transport,
		decoder:   xml.NewDecoder(transport),
		encoder:   xml.NewEncoder(transport),
	}

	conn.outStream = outStream
	conn.inStream, err = conn.requestStream(outStream)

	if err != nil {
		return nil, err
	}

	err = conn.readFeatures()

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Transport returns currently used transport object
func (conn *Conn) Transport() io.ReadWriter {
	return conn.transport
}

// Features returns current connection features
func (conn *Conn) Features() *Features {
	return conn.features
}

// Read reads and returns a message from the stream
func (conn *Conn) Read() (interface{}, error) {
	for {
		token, err := conn.decoder.Token()

		if err != nil {
			return nil, err
		}

		switch start := token.(type) {
		case xml.StartElement:
			return StreamContext.DecodeElement(conn.decoder, &start)

		case xml.EndElement: // </stream>
			return nil, ErrEndOfStream

		case xml.CharData:
			// Ignore character data
		default:
			return nil, ErrStreamError
		}
	}
}

// Close closes the XMPP stream. No writes can be performed afterwards.
func (conn *Conn) Close() {
	conn.writeStreamEnd()
	conn.encoder = nil
}

// Write writes an XML message to the XMPP stream
func (conn *Conn) Write(msg interface{}) error {
	conn.Lock()
	defer conn.Unlock()

	if err := conn.encoder.Encode(msg); err != nil {
		return err
	}

	return conn.encoder.Flush()
}

// requestStream tries to open a bidirectional stream
func (conn *Conn) requestStream(out *StreamHeader) (*StreamHeader, error) {
	conn.writeStreamHeader(out)

	return conn.readStreamHeader()
}

// writeStreamHeader writes an XMPP stream opening tag
func (conn *Conn) writeStreamHeader(stream *StreamHeader) error {
	// Write <?xml version="1.0"?>
	err := conn.encoder.EncodeToken(xml.ProcInst{
		Target: "xml",
		Inst:   []byte(`version="1.0"`),
	})

	if err != nil {
		return err
	}

	start := stream.XMLStartElement()
	conn.streamEnd = start.End() // Store matching end element

	conn.Lock()
	defer conn.Unlock()

	err = conn.encoder.EncodeToken(start)
	if err != nil {
		return err
	}

	return conn.encoder.Flush()
}

// readStreamHeader reads and returns a stream header
func (conn *Conn) readStreamHeader() (*StreamHeader, error) {
	for {
		token, err := conn.decoder.Token()

		if err != nil {
			return nil, err
		}

		switch start := token.(type) {
		case xml.ProcInst:
			if conn.piReceived {
				return nil, ErrStreamError
			}
			conn.piReceived = true // PI is only valid once

		case xml.StartElement:
			conn.piReceived = true // PI is only valid before any elements
			return ParseStreamHeader(&start), nil

		default:
			return nil, ErrStreamError
		}
	}
}

// writeStreamEnd writes the XMPP stream end element
func (conn *Conn) writeStreamEnd() error {
	conn.Lock()
	defer conn.Unlock()

	err := conn.encoder.EncodeToken(conn.streamEnd)

	if err != nil {
		return err
	}

	return conn.encoder.Flush()
}

// readFeatures reads stream features from the server. Should only be called
// by a client directly after establishing an XMPP stream.
func (conn *Conn) readFeatures() error {
	for {
		msg, err := conn.Read()

		if err != nil {
			return err
		}

		switch typed := msg.(type) {
		case *Features:
			conn.features = typed
			return nil

		default:
			return errors.New("unexpected message")
		}
	}
}
