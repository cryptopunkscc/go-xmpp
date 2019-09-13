package xmpp

import (
	"encoding/xml"
	"errors"
	"io"
	"sync"
)

// Stream manages a bidirectional XMPP stream using an underlying transport
type Stream struct {
	transport  io.ReadWriter
	decoder    *xml.Decoder
	encoder    *xml.Encoder
	piReceived bool
	streamEnd  xml.EndElement

	features *Features

	localHeader  *StreamHeader
	remoteHeader *StreamHeader

	sync.Mutex
}

type StreamFunc func(*StreamHeader) *StreamHeader

// Open a new bidirectional XMPP stream on the provided transport
func Open(transport io.ReadWriter, header *StreamHeader) (*Stream, error) {
	var err error

	s := &Stream{
		transport: transport,
		decoder:   xml.NewDecoder(transport),
		encoder:   xml.NewEncoder(transport),
	}

	s.localHeader = header
	s.remoteHeader, err = s.requestStream(header)

	if err != nil {
		return nil, err
	}

	err = s.readFeatures()

	if err != nil {
		return nil, err
	}

	return s, nil
}

func Accept(transport io.ReadWriter, header *StreamHeader) (*Stream, error) {
	s := &Stream{
		transport: transport,
		decoder:   xml.NewDecoder(transport),
		encoder:   xml.NewEncoder(transport),
	}

	remoteHeader, err := s.readStreamHeader()

	if err != nil {
		return nil, err
	}

	s.remoteHeader = remoteHeader
	s.writeStreamHeader(header)

	return s, nil
}

func AcceptFunc(transport io.ReadWriter, streamFunc StreamFunc) (*Stream, error) {
	var err error
	s := &Stream{
		transport: transport,
		decoder:   xml.NewDecoder(transport),
		encoder:   xml.NewEncoder(transport),
	}
	s.remoteHeader, err = s.readStreamHeader()
	if err != nil {
		return nil, err
	}
	s.localHeader = streamFunc(s.remoteHeader)
	if s.localHeader == nil {
		return nil, errors.New("Stream function returned no header")
	}
	return s, s.writeStreamHeader(s.localHeader)
}

// RemoteHeader returns a stream header received from the remote host
func (s *Stream) RemoteHeader() *StreamHeader {
	return s.remoteHeader
}

// LocalHeader returns a stream header that was sent to the remote host
func (s *Stream) LocalHeader() *StreamHeader {
	return s.localHeader
}

// Transport returns currently used transport object
func (s *Stream) Transport() io.ReadWriter {
	return s.transport
}

// Features returns current connection features
func (s *Stream) Features() *Features {
	return s.features
}

// Read reads and returns a message from the stream
func (s *Stream) Read() (interface{}, error) {
	for {
		token, err := s.decoder.Token()

		if err != nil {
			return nil, err
		}

		switch start := token.(type) {
		case xml.StartElement:
			return StreamContext.DecodeElement(s.decoder, &start)

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

// Write writes an XML message to the XMPP stream
func (s *Stream) Write(msg interface{}) error {
	s.Lock()
	defer s.Unlock()

	if err := s.encoder.Encode(msg); err != nil {
		return err
	}

	return s.encoder.Flush()
}

// requestStream tries to open a bidirectional stream
func (s *Stream) requestStream(out *StreamHeader) (*StreamHeader, error) {
	s.writeStreamHeader(out)

	return s.readStreamHeader()
}

// writeStreamHeader writes an XMPP stream opening tag
func (s *Stream) writeStreamHeader(stream *StreamHeader) error {
	// Write <?xml version="1.0"?>
	err := s.encoder.EncodeToken(xml.ProcInst{
		Target: "xml",
		Inst:   []byte(`version="1.0"`),
	})

	if err != nil {
		return err
	}

	start := stream.XMLStartElement()
	s.streamEnd = start.End() // Store matching end element

	s.Lock()
	defer s.Unlock()

	err = s.encoder.EncodeToken(start)
	if err != nil {
		return err
	}

	return s.encoder.Flush()
}

// readStreamHeader reads and returns a stream header
func (s *Stream) readStreamHeader() (*StreamHeader, error) {
	for {
		token, err := s.decoder.Token()

		if err != nil {
			return nil, err
		}

		switch start := token.(type) {
		case xml.ProcInst:
			if s.piReceived {
				return nil, ErrStreamError
			}
			s.piReceived = true // PI is only valid once

		case xml.StartElement:
			s.piReceived = true // PI is only valid before any elements
			return ParseStreamHeader(&start), nil

		default:
			return nil, ErrStreamError
		}
	}
}

// writeStreamEnd writes the XMPP stream end element
func (s *Stream) writeStreamEnd() error {
	s.Lock()
	defer s.Unlock()

	err := s.encoder.EncodeToken(s.streamEnd)

	if err != nil {
		return err
	}

	return s.encoder.Flush()
}

// readFeatures reads stream features from the server. Should only be called
// by a client directly after establishing an XMPP stream.
func (s *Stream) readFeatures() error {
	for {
		msg, err := s.Read()

		if err != nil {
			return err
		}

		switch typed := msg.(type) {
		case *Features:
			s.features = typed
			return nil

		default:
			return errors.New("unexpected message")
		}
	}
}
