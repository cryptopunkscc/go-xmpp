package xmpp

import "encoding/xml"

// Errors are a little tricky to unmarshal, because of their irregular structure:
//
// <stream:error>
//   <defined-condition xmlns='urn:ietf:params:xml:ns:xmpp-streams'/>
//   [<text xmlns='urn:ietf:params:xml:ns:xmpp-streams'
//          xml:lang='langcode'>
//      OPTIONAL descriptive text
//   </text>]
//   [OPTIONAL application-specific condition element]
// </stream:error>
//
// This implementation doesn't handle <text> and custom elements yet

type Error struct {
	Condition string
	Text      string
	Type      string
}

type pError struct {
	XMLName xml.Name `xml:"stream:error"`
	Type    string   `xml:"type,attr,omitempty"`
	Error   struct {
		XMLName xml.Name
	}
}

func (e *Error) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	c := &struct{ XMLName xml.Name }{}

	dec.Decode(&c)

	e.Condition = c.XMLName.Local

	dec.Skip()

	return nil
}

func (e *Error) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	p := &pError{}

	p.Error.XMLName.Local = e.Condition
	p.Error.XMLName.Space = "urn:ietf:params:xml:ns:xmpp-streams"

	enc.Encode(p)

	return nil
}

func addErrorElements() {
	AddElement(&Error{})
}
