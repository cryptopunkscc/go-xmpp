package xmpp

import "encoding/xml"

var StreamContext = NewContext(nil)

type Abort struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl abort"`
}

type Auth struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl auth"`
	Mechanism string   `xml:"mechanism,attr"`
	Data      string   `xml:",chardata"`
}

type Bind struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	JID      string   `xml:"jid,omitempty"`
	Resource string   `xml:"resource,omitempty"`
}

type Challenge struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl challenge"`
	Data    string   `xml:",chardata"`
}

type TLSFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls failure"`
}

type SASLFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl failure"`
}

type Proceed struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

type Response struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl response"`
	Data    string   `xml:",chardata"`
}

type StartTLS struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
}

type Success struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
}

func initStreamContext() {
	StreamContext.Add(&Abort{})
	StreamContext.Add(&Auth{})
	StreamContext.Add(&Bind{})
	StreamContext.Add(&Challenge{})
	StreamContext.Add(&Error{})
	StreamContext.Add(&TLSFailure{})
	StreamContext.Add(&SASLFailure{})
	StreamContext.Add(&Features{})
	StreamContext.Add(&Proceed{})
	StreamContext.Add(&StartTLS{})
	StreamContext.Add(&Success{})
}
