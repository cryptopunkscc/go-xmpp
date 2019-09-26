package xmpp

import "encoding/xml"

type Stanza interface {
	GetID() string
	GetFrom() JID
	GetTo() JID
	GetType() string
	GetLang() string

	SetID(string)
	SetFrom(JID)
	SetTo(JID)
	SetType(string)
	SetLang(string)
}

func stanzaToStartElement(s Stanza) xml.StartElement {
	se := xml.StartElement{
		Name: Identify(s),
	}
	se.Attr = make([]xml.Attr, 0)

	set := func(f, v string) {
		if v != "" {
			se.Attr = append(se.Attr, xml.Attr{Name: xml.Name{Local: f}, Value: v})
		}
	}

	set("id", s.GetID())
	set("from", s.GetFrom().String())
	set("to", s.GetTo().String())
	set("type", s.GetType())
	set("xml:lang", s.GetLang())

	return se
}

func startElementToStanza(se xml.StartElement, s Stanza) {
	get := func(f string) string {
		for _, attr := range se.Attr {
			if attr.Name.Local == f {
				return attr.Value
			}
		}
		return ""
	}

	s.SetID(get("id"))
	s.SetFrom(JID(get("from")))
	s.SetTo(JID(get("to")))
	s.SetType(get("type"))
	s.SetLang(get("xml:lang"))
}
