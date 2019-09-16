package xmpp

import (
	"encoding/xml"
	"fmt"
	"io"
	"sync"
)

// Stream manages a bidirectional XMPP stream using an underlying transport
type Stream struct {
	transport io.ReadWriter
	decoder   *xml.Decoder
	encoder   *xml.Encoder
	streamEnd xml.EndElement

	mu sync.Mutex
}

// NewStream instantiates a new stream using the provided transport
func NewStream(transport io.ReadWriter) *Stream {
	return &Stream{
		transport: transport,
		decoder:   xml.NewDecoder(transport),
		encoder:   xml.NewEncoder(transport),
	}
}

// Transport returns the underlying transport object
func (s *Stream) Transport() io.ReadWriter {
	return s.transport
}

// Read reads and returns a message from the stream
func (s *Stream) Read() (interface{}, error) {
	for {
		token, err := s.decoder.Token()
		if err != nil {
			return nil, err
		}

		switch typed := token.(type) {
		case xml.StartElement:
			return StreamContext.DecodeElement(s.decoder, &typed)

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
func (s *Stream) Close() {
	s.writeStreamEnd()
	s.encoder = nil
}

// Write writes an XMPP message to the stream
func (s *Stream) Write(msg interface{}) error {
	// Lock the stream for writing
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.encoder.Encode(msg); err != nil {
		return err
	}
	return s.encoder.Flush()
}

// WriteHeader writes a stream header to the stream.
func (s *Stream) WriteHeader(stream *StreamHeader) error {
	// Lock the stream for writing
	s.mu.Lock()
	defer s.mu.Unlock()

	// Write ProcInst (<?xml version="1.0"?>)
	err := s.encoder.EncodeToken(xml.ProcInst{
		Target: "xml",
		Inst:   []byte(`version="1.0"`),
	})
	if err != nil {
		return err
	}

	// Open the stream element and keep the closing tag for later
	start := stream.XMLStartElement()
	s.streamEnd = start.End() // Store matching end element
	err = s.encoder.EncodeToken(start)
	if err != nil {
		return err
	}
	return s.encoder.Flush()
}

// ReadHeader reads and returns a stream header. It ignores the first ProcInst
// read from the stream.
func (s *Stream) ReadHeader() (*StreamHeader, error) {
	allowProcInst := true
	for {
		token, err := s.decoder.Token()
		if err != nil {
			return nil, err
		}

		switch typed := token.(type) {
		case xml.ProcInst:
			if !allowProcInst {
				return nil, ErrStreamError
			}
			allowProcInst = false

		case xml.StartElement:
			allowProcInst = false // PI is only valid before other elements
			return ParseStreamHeader(&typed), nil

		default:
			return nil, ErrStreamError
		}
	}
}

// ReadFeatures reads stream features from the server. Should only be called
// by the client directly after receiveng a stream header.
func (s *Stream) ReadFeatures() (*Features, error) {
	msg, err := s.Read()
	if err != nil {
		return nil, err
	}
	if feats, ok := msg.(*Features); ok {
		return feats, nil
	}
	n, _ := Identify(msg)
	return nil, fmt.Errorf("ReadFeatures: unexpected message: %s", n)
}

// writeStreamEnd writes the XMPP stream end element
func (s *Stream) writeStreamEnd() error {
	// Lock the stream for writing
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.encoder.EncodeToken(s.streamEnd)
	if err != nil {
		return err
	}
	return s.encoder.Flush()
}
