package xmpp

import "encoding/xml"

// SASL structures (urn:ietf:params:xml:ns:xmpp-sasl)

type Abort struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl abort"`
}

type Auth struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl auth"`
	Mechanism string   `xml:"mechanism,attr"`
	Data      string   `xml:",chardata"`
}

type Challenge struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl challenge"`
	Data    string   `xml:",chardata"`
}

type SASLFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl failure"`
}

type Mechanisms struct {
	XMLName    xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Mechanisms []string `xml:"mechanism"`
}

type Response struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl response"`
	Data    string   `xml:",chardata"`
}

type Success struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
	Data    string   `xml:",chardata"`
}

func addSASLElements() {
	AddElement(&Abort{})
	AddElement(&Auth{})
	AddElement(&Challenge{})
	AddElement(&SASLFailure{})
	AddElement(&Mechanisms{})
	AddElement(&Response{})
	AddElement(&Success{})
}
