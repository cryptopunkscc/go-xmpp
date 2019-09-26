package xep0030

import (
	"encoding/xml"

	"github.com/cryptopunkscc/go-xmpp"
)

// QueryInfo represents a query for disco info
type QueryInfo struct {
	XMLName    xml.Name   `xml:"http://jabber.org/protocol/disco#info query"`
	Node       string     `xml:"node,attr,omitempty"`
	Identities []Identity `xml:"identity"`
	Features   []Feature  `xml:"feature"`
}

// QueryItems represents a query for disco items
type QueryItems struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#items query"`
	Node    string   `xml:"node,attr,omitempty"`
	Items   []Item   `xml:"item"`
}

// Identity represents a XEp0030 identitiy
type Identity struct {
	Category string `xml:"category,attr"`
	Type     string `xml:"type,attr"`
	Name     string `xml:"name,attr"`
	Lang     string `xml:"lang,attr,omitempty"`
}

// Feature represents a XEP0030 feature
type Feature struct {
	Var string `xml:"var,attr"`
}

// Item represents a XEP0030 item
type Item struct {
	JID  xmpp.JID `xml:"jid,attr"`
	Node string   `xml:"node,attr,omitempty"`
	Name string   `xml:"name,attr,omitempty"`
}

func init() {
	xmpp.AddElement(&QueryInfo{})
	xmpp.AddElement(&QueryItems{})
}
