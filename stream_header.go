package xmpp

import (
	"encoding/xml"
)

// Stream namespaces
const (
	NamespaceClient = "jabber:client"
	NamespaceServer = "jabber:server"
	DefaultVersion  = "1.0"
)

// StreamHeader represents stream info exchanged in the initial <stream/> element
type StreamHeader struct {
	Version   string
	Namespace string
	From      string
	To        string
	ID        string
}

// NewHeader returns a stream header in client protocol namespace
func NewHeader(namesapce string, from string, to string) *StreamHeader {
	return &StreamHeader{
		Version:   DefaultVersion,
		Namespace: namesapce,
		From:      from,
		To:        to,
	}
}

// XMLStartElement returns an XML start element containing all stream info
func (header *StreamHeader) XMLStartElement() (start xml.StartElement) {
	start.Name = xml.Name{
		Space: header.Namespace,
		Local: "stream:stream",
	}
	start.Attr = make([]xml.Attr, 0)
	if header.To != "" {
		start.Attr = append(start.Attr, xml.Attr{xml.Name{"", "to"}, header.To})
	}
	if header.From != "" {
		start.Attr = append(start.Attr, xml.Attr{xml.Name{"", "from"}, header.From})
	}
	if header.ID != "" {
		start.Attr = append(start.Attr, xml.Attr{xml.Name{"", "id"}, header.ID})
	}
	start.Attr = append(start.Attr,
		xml.Attr{xml.Name{"", "version"}, header.Version},
		xml.Attr{xml.Name{"", "xml:lang"}, "en"},
		xml.Attr{xml.Name{"", "xmlns:stream"}, nsEtherXStreams},
	)
	return
}

// Reply returns a server stream header
func (header *StreamHeader) Reply(id string) *StreamHeader {
	return &StreamHeader{
		To:        header.From,
		From:      header.To,
		Namespace: header.Namespace,
		Version:   header.Version,
		ID:        id,
	}
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
