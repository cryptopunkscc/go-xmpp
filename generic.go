package xmpp

import "encoding/xml"

type Generic struct {
	Local    string
	Space    string
	Attrs    map[string]string
	Children []Generic `xml:",any"`
	Text     string    `xml:",chardata"`
}

func (u *Generic) Name() string {
	return u.Local
}

func (u *Generic) Namespace() string {
	return u.Space
}

func (u *Generic) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type rawType Generic
	raw := &rawType{}
	dec.DecodeElement(raw, &start)

	u.Children = raw.Children

	u.Local = start.Name.Local
	u.Space = start.Name.Space
	u.Attrs = make(map[string]string, 0)

	for _, attr := range start.Attr {
		u.Attrs[attr.Name.Local] = attr.Value
	}

	return nil
}
