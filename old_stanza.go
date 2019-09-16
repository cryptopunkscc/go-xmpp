package xmpp

import "encoding/xml"

type OldStanza struct {
	Container
	Context *Context
	Stanza  string
	ID      string
	From    string
	To      string
	Type    string
	Lang    string
}

func (stanza *OldStanza) Name() string {
	return stanza.Stanza
}

func (stanza *OldStanza) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var err error

	stanza.copyStartElement(&start)
	stanza.Children, err = stanza.Context.DecodeAll(dec)

	return err
}

func (stanza *OldStanza) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	s := stanza.startElement()
	enc.EncodeToken(s)

	if err := EncodeAll(enc, stanza.Children); err != nil {
		return err
	}

	return enc.EncodeToken(s.End())
}

func (stanza *OldStanza) startElement() xml.StartElement {
	s := xml.StartElement{
		Name: xml.Name{
			Local: stanza.Stanza,
		},
	}

	s.Attr = make([]xml.Attr, 0)

	set := func(f, v string) {
		if v != "" {
			s.Attr = append(s.Attr, xml.Attr{Name: xml.Name{Local: f}, Value: v})
		}
	}

	set("id", stanza.ID)
	set("from", stanza.From)
	set("to", stanza.To)
	set("type", stanza.Type)
	set("xml:lang", stanza.Lang)

	return s
}

func (stanza *OldStanza) copyStartElement(s *xml.StartElement) {
	get := func(f string) string {
		for _, attr := range s.Attr {
			if attr.Name.Local == f {
				return attr.Value
			}
		}
		return ""
	}

	stanza.Stanza = s.Name.Local
	stanza.ID = get("id")
	stanza.From = get("from")
	stanza.To = get("to")
	stanza.Type = get("type")
	stanza.Lang = get("xml:lang")
}
