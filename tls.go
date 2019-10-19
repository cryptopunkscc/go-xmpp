package xmpp

import "encoding/xml"

// TLS structures (urn:ietf:params:xml:ns:xmpp-tls)

type TLSFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls failure"`
}

type Proceed struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

type StartTLS struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required *string  `xml:"required,omitempty"`
}

func addTLSElements() {
	AddElement(&TLSFailure{})
	AddElement(&Proceed{})
	AddElement(&StartTLS{})
}
