package xep0115

import (
	"encoding/xml"

	"github.com/cryptopunkscc/go-xmpp"
)

// Capability represents the <c/> element described in XEP-0115
type Capability struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/caps c"`
	Hash    string   `xml:"hash,attr"`
	Node    string   `xml:"node,attr"`
	Ver     string   `xml:"ver,attr"`
}

func init() {
	// https://xmpp.org/extensions/xep-0115.html#howitworks
	// https://xmpp.org/extensions/xep-0115.html#stream
	xmpp.AddElement(&Capability{})
}
