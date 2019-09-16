package xmpp

import "encoding/xml"

type Stanza struct {
	ID   string `xml:"id,attr,omitempty"`
	To   string `xml:"to,attr,omitempty"`
	From string `xml:"from,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
	Lang string `xml:"lang,attr,omitempty"`
}

func (s *Stanza) startElement(name string) xml.StartElement {
	se := xml.StartElement{
		Name: xml.Name{
			Local: name,
		},
	}
	se.Attr = make([]xml.Attr, 0)

	set := func(f, v string) {
		if v != "" {
			se.Attr = append(se.Attr, xml.Attr{Name: xml.Name{Local: f}, Value: v})
		}
	}

	set("id", s.ID)
	set("from", s.From)
	set("to", s.To)
	set("type", s.Type)
	set("xml:lang", s.Lang)

	return se
}

func (s *Stanza) copyStartElement(se *xml.StartElement) {
	get := func(f string) string {
		for _, attr := range se.Attr {
			if attr.Name.Local == f {
				return attr.Value
			}
		}
		return ""
	}

	s.ID = get("id")
	s.From = get("from")
	s.To = get("to")
	s.Type = get("type")
	s.Lang = get("xml:lang")
}
