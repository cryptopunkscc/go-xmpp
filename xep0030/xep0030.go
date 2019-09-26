package xep0030

import (
	"encoding/xml"

	"github.com/cryptopunkscc/go-xmpp"
)

type QueryInfo struct {
	XMLName    xml.Name   `xml:"http://jabber.org/protocol/disco#info query"`
	Node       string     `xml:"node,attr,omitempty"`
	Identities []Identity `xml:"identity"`
	Features   []Feature  `xml:"feature"`
}

type Identity struct {
	Category string `xml:"category,attr"`
	Type     string `xml:"type,attr"`
	Name     string `xml:"name,attr"`
	Lang     string `xml:"lang,attr,omitempty"`
}

type Feature struct {
	Var string `xml:"var,attr"`
}

type QueryItems struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#items query"`
	Node    string   `xml:"node,attr,omitempty"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	JID  string `xml:"jid,attr"`
	Name string `xml:"name,attr"`
}

func init() {
	xmpp.AddElement(&QueryInfo{})
	xmpp.AddElement(&QueryItems{})
}
