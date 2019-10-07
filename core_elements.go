package xmpp

import "encoding/xml"

type Bind struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	JID      JID      `xml:"jid,omitempty"`
	Resource string   `xml:"resource,omitempty"`
}

type Compression struct {
	XMLName xml.Name `xml:"http://jabber.org/features/compress compression"`
	Methods []string `xml:"method"`
}

type Register struct {
	XMLName xml.Name `xml:"http://jabber.org/features/iq-register register"`
}

func addCoreElements() {
	AddElement(&Bind{})
	AddElement(&Compression{})
	AddElement(&Register{})
}
