package xmpp

import (
	"encoding/xml"
)

// Stream namespaces
const (
	NamespaceClient = "jabber:client"
	NamespaceServer = "jabber:server"
)

// StreamHeader represents stream info exchanged in the initial <stream/> element
type StreamHeader struct {
	Version   string
	Namespace string
	From      string
	To        string
	ID        string
}

// ClientHeader returns a stream header in client protocol namespace
func ClientHeader(from string, to string) *StreamHeader {
	return &StreamHeader{
		Version:   "1.0",
		Namespace: NamespaceClient,
		From:      from,
		To:        to,
	}
}

// XMLStartElement returns an XML start element containing all stream info
func (stream *StreamHeader) XMLStartElement() (start xml.StartElement) {
	start.Name = xml.Name{
		Space: stream.Namespace,
		Local: "stream:stream",
	}
	start.Attr = make([]xml.Attr, 0)
	if stream.To != "" {
		start.Attr = append(start.Attr, xml.Attr{xml.Name{"", "to"}, stream.To})
	}
	if stream.From != "" {
		start.Attr = append(start.Attr, xml.Attr{xml.Name{"", "from"}, stream.From})
	}
	if stream.ID != "" {
		start.Attr = append(start.Attr, xml.Attr{xml.Name{"", "id"}, stream.ID})
	}
	start.Attr = append(start.Attr,
		xml.Attr{xml.Name{"", "version"}, "1.0"},
		xml.Attr{xml.Name{"", "xml:lang"}, "en"},
		xml.Attr{xml.Name{"", "xmlns:stream"}, "http://etherx.jabber.org/streams"},
	)
	return
}

// ParseStreamHeader decodes the opening of the XML stream without touching
// the content of the stream
func ParseStreamHeader(start *xml.StartElement) *StreamHeader {
	s := &StreamHeader{}
	for _, a := range start.Attr {
		switch a.Name.Local {
		case "to":
			s.To = a.Value
		case "from":
			s.From = a.Value
		case "id":
			s.ID = a.Value
		case "xmlns":
			s.Namespace = a.Value
		case "version":
			s.Version = a.Value
		}
	}
	return s
}
