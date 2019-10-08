package xmpp

import (
	"encoding/xml"
	"errors"
)

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

// DefinedConditions is a list of conditions defined by RFC6120
var DefinedConditions = []string{
	"bad-format",
	"bad-namespace-prefix",
	"conflict",
	"connection-timeout",
	"host-gone",
	"host-unknown",
	"improper-addressing",
	"internal-server-error",
	"invalid-from",
	"invalid-namespace",
	"invalid-xml",
	"not-authorized",
	"not-well-formed",
	"policy-violation",
	"remote-connection-failed",
	"reset",
	"resource-constraint",
	"restricted-xml",
	"see-other-host",
	"system-shutdown",
	"undefined-condition",
	"unsupported-encoding",
	"unsupported-feature",
	"unsupported-stanza-type",
	"unsupported-version",
}

type Error struct {
	Condition string
	Text      string
	Extra     interface{}
}

func (e *Error) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	// First element should be the condition
	token, err := dec.Token()
	if err != nil {
		return err
	}
	start, ok := token.(xml.StartElement)
	if !ok {
		return errors.New("failed to parse error: unexpected element")
	}
	e.Condition = start.Name.Local
	if err := dec.Skip(); err != nil {
		return err
	}

	// Get the optional text element
	token, err = dec.Token()
	if err != nil {
		return err
	}
	start, ok = token.(xml.StartElement)
	if !ok {
		return errors.New("failed to parse error: unexpected element")
	}
	if start.Name.Local == "text" {
		textElem := &struct {
			Text string `xml:",chardata"`
		}{}
		if err := dec.DecodeElement(textElem, &start); err != nil {
			return err
		}
		e.Text = textElem.Text
	} else {
		p := &proxy{}
		if err := dec.DecodeElement(p, &start); err != nil {
			return err
		}
		e.Extra = p.Object
		return dec.Skip()
	}

	// Get the optional application specific element
	token, err = dec.Token()
	if err != nil {
		return err
	}
	if _, ok := token.(xml.EndElement); ok {
		return nil
	}
	start, ok = token.(xml.StartElement)
	if !ok {
		return errors.New("failed to parse error: unexpected element")
	}
	p := &proxy{}
	if err := dec.DecodeElement(p, &start); err != nil {
		return err
	}
	e.Extra = p.Object
	return dec.Skip()
}

func (e *Error) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "error"
	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	// Write error
	condition := &struct {
		XMLName xml.Name
	}{
		XMLName: xml.Name{
			Space: nsStreams,
			Local: e.Condition,
		},
	}
	if err := enc.Encode(condition); err != nil {
		return err
	}

	// Write text
	if e.Text != "" {
		text := &struct {
			XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-streams text"`
			Text    string   `xml:",chardata"`
		}{
			Text: e.Text,
		}
		if err := enc.Encode(text); err != nil {
			return err
		}
	}

	// Write extra element
	if e.Extra != nil {
		return enc.Encode(e.Extra)
	}

	return enc.EncodeToken(start.End())
}

func addErrorElements() {
	AddElement(&Error{})
}
